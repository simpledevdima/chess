package main

import (
	"github.com/simpledevdima/chess/server"
	"html/template"
	"log"
	"net/http"
)

func client(w http.ResponseWriter, _ *http.Request) {
	tpl, err := template.ParseFiles("client.gohtml")
	if err != nil {
		log.Println(err)
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func handlers() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.Handle("/javascript/", http.StripPrefix("/javascript/", http.FileServer(http.Dir("./javascript/"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images/"))))
	http.HandleFunc("/chess", client)
}

// chess client application start
func main() {
	// chess server start
	go server.Start("server/config.yaml")

	handlers()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
	}
}
