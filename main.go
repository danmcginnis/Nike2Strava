//Handles moving run data from Nike+ into Strava
package nike2strava

import (
	"net/http"
	"html/template"
    "log"
)

var userToken string

func echoInput(w http.ResponseWriter, r *http.Request) {
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
