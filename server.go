package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

type Tip struct {
	gorm.Model
	FromId int
	ToId int
	MessageId int
	Amount float64
	From User
	To User
}

type User struct {
	gorm.Model
	Username string
	Address string
	Balance	float64
}

type Params struct {
	Limit int
}

func (tip *Tip) Find(p Params) ([]Tip, error) {
	var tips []Tip
	if p.Limit == 0 {
		DB.Preload("From").Preload("To").Find(&tips ,tip)
		return tips, DB.Error
	}
	DB.Preload("From").Preload("To").Limit(p.Limit).Find(&tips ,tip)
	return tips, DB.Error
}

func Count() (int64, error) {
	var count int64
	DB.Table("tips").Count(&count)
	return count, DB.Error
}

func TippedAmount() (float64, error) {
	type Result struct {
		Total float64
	}

	var result Result
	DB.Model(&Tip{}).Select("sum(amount) as total").Scan(&result)

	return result.Total, DB.Error
}



func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	initDB()
	defer DB.Close()

	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/templates/", http.StripPrefix("/templates/", fs))

	http.HandleFunc("/", ServeTemplate)

	fmt.Println("Listening...")
	err = http.ListenAndServe(GetPort(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return
	}
}

func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func ServeTemplate(w http.ResponseWriter, r *http.Request) {
	lp := path.Join("templates", "layout.html")
	fp := path.Join("templates", r.URL.Path)

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	templates, err := template.ParseFiles(lp, fp)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	tip := Tip{}
	tips, err := tip.Find(Params{Limit:100})
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	tip = Tip{}
	count, err := Count()
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	total_tips, err := TippedAmount()
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	data := struct {
		Tips []Tip
		Total float64
		Count int64
	}{
		tips,
		total_tips,
		count,
	}

	err = templates.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Println(err)
	}
}

func initDB() {
	log.Println("Connecting to DB...")
	var err error
	DB, err = gorm.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@/"+os.Getenv("DB_NAME")+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to DB")
}
