package main

import (
	"fmt"
	"net/http"
  "html/template"
)

type Page struct {
  User string
  Valuetype string
  Order string
  Message string
	Color string
	Station string
}

type Message struct {
	Success bool
	User string
	Message string
}

func main() {
	http.HandleFunc("/", pick)
  http.HandleFunc("/scan", scan)
	http.ListenAndServe(":8080", nil)
}

func scan(w http.ResponseWriter, r *http.Request) {
	var typevalue string
	var override bool
	station := r.FormValue("station")
	color := "9ccdff"
  fmt.Println("type:"+r.FormValue("type"))
  fmt.Println("station:"+station)
  if r.FormValue("type") == "user" {
		message := userauth(r.FormValue("value"))
		if message.Success == false {
			typevalue = "user"
			color="ff2200"
		} else {
			typevalue = "order"
			color="00ff7b"
		}
    http.Redirect(w, r, "/?user="+message.User+"&type="+typevalue+"&message="+message.Message+"&color="+color+"&station="+station, http.StatusSeeOther)
  }
  if r.FormValue("type") == "order" {
		if r.FormValue("override") == "override" {override=true} else {override=false}
		message := insert(r.FormValue("user"),r.FormValue("value"),r.FormValue("station"),override)
		if message.Success == false {
			typevalue = "order"
			color="ff2200"
		} else {
			typevalue = "user"
		}
    fmt.Println(r.FormValue("station"))
	  http.Redirect(w, r, "/?user="+r.FormValue("user")+"&type="+typevalue+"&message="+message.Message+"&color="+color+"&order="+r.FormValue("value")+"&station="+r.FormValue("station"), http.StatusSeeOther)
  }
}

func pick(w http.ResponseWriter, r *http.Request) {
  var Data Page
  Data.User=r.URL.Query().Get("user")
  Data.Valuetype=r.URL.Query().Get("type")
	Data.Message=r.URL.Query().Get("message")
	Data.Color=r.URL.Query().Get("color")
	Data.Order=r.URL.Query().Get("order")
	Data.Station=r.URL.Query().Get("station")
	fmt.Println(Data.Station)
	if Data.Color == "" {
		Data.Color = "9ccdff"
	}
  if Data.User == "notfound" || Data.User == "" {
    Data.Valuetype = "user"
  }
	tmpl, err := template.ParseFiles("layout.html")
  fmt.Println(err)
  tmpl.Execute(w, Data)
}
