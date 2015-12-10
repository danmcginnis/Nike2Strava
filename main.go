package main

import (
	"net/http"
	"fmt"
	"log"
	"html/template"
)

func echoInput(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("public/index.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}


func main() {
	http.HandleFunc("/", echoInput)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}