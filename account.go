package main

import (
	"html/template"
	"net/http"
)

func accountBindForm(w http.ResponseWriter, r *http.Request) {
	t := template.New("main")
	_, err := t.ParseFiles("templates/accountbindform.html")
	if err != nil {
		return
	}
	t.Execute(w, nil)
}
