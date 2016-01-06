package main

import (
	"fmt"
	"net/http"
	//"log"
	"html/template"
)

var userToken string

func echoInput(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("public/index.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		userToken = template.HTMLEscapeString(r.Form.Get("token"))
	}
}

func main() {

	webDebug := false

	if webDebug {
		http.HandleFunc("/", echoInput)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else {
        userToken = Token
    }
    
	wrangleJSON(userToken, 20, true)
}
