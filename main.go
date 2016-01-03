package main

import (
	"net/http"
	"fmt"
	//"log"
	"html/template"
)

func echoInput(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("public/index.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		username := template.HTMLEscapeString(r.Form.Get("username"))
		password := template.HTMLEscapeString(r.Form.Get("password"))
		fmt.Println("username:", username)
		fmt.Println("password:", password)
	}
}


func main() {
	wrangleJSON(Token, 20, true)
	/*http.HandleFunc("/", echoInput)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} */
}