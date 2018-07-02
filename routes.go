package main

import (
	"net/http"
	"log"
)

func initializeRoutes() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", HomepageFunc)
	http.HandleFunc("/stats", StatpageFunc)
}

func StatpageFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/stats" {
		http.NotFound(w,r)
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

	err = templates.ExecuteTemplate(w, "statpage", data)
	if err != nil {
		log.Println(err)
	}
}


func HomepageFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	err := templates.ExecuteTemplate(w, "homepage", nil)
	if err != nil {
		log.Println(err)
	}
}