package main

import (
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

type User struct {
	gorm.Model
	Username string
	Address string
	Balance	float64
}

type Params struct {
	Limit int
}

type Tip struct {
	gorm.Model
	FromId int
	ToId int
	MessageId int
	Amount float64
	From User
	To User
}

var DB *gorm.DB

func initDB() {
	log.Println("Connecting to DB...")
	var err error
	DB, err = gorm.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@/"+os.Getenv("DB_NAME")+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to DB")
}

func (tip *Tip) Find(p Params) ([]Tip, error) {
	var tips []Tip
	if p.Limit == 0 {
		DB.Preload("From").Preload("To").Find(&tips ,tip)
		return tips, DB.Error
	}
	DB.Preload("From").Preload("To").Order("created_at desc").Limit(p.Limit).Find(&tips ,tip)
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
