package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type CiscoIPPhoneDirectory struct {
	Title          string
	Prompt         string
	DirectoryEntry []phonenumber
}

type CiscoIPPhoneMenu struct {
	Title    string
	Prompt   string
	MenuItem []menuitem
}

type menuitem struct {
	Name string
	URL  string
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	var mi menuitem
	mi.Name = "Telefonbuch"
	mi.URL = "http://" + r.Host + "/dir" + r.URL.String()

	var menu CiscoIPPhoneMenu
	menu.MenuItem = append(menu.MenuItem, mi)

	x, err := xml.Marshal(menu)
	if err != nil {
		log.Fatal(err)
	}

	xs := []byte(xml.Header + string(x))

	w.Header().Set("Content-type", "text/xml")

	w.Write(xs)
}

func phonebookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Access")
	vars := mux.Vars(r)
	extension := vars["extension"]
	fmt.Println(extension)

	pns, err := GetNumbersforExtension(extension)
	if err != nil {
		log.Println(err)
	}

	var dir CiscoIPPhoneDirectory
	dir.Title = "Telefonbuch"
	dir.Prompt = ""
	dir.DirectoryEntry = pns

	x, err := xml.Marshal(dir)
	if err != nil {
		log.Fatal(err)
	}

	xs := []byte(xml.Header + string(x))

	w.Header().Set("Content-type", "text/xml")

	w.Write(xs)

}

func main() {
	var dbc DBConnection
	dbc.Driver = "sqlite3"
	dbc.Connection = "./phonebook.db"

	initialisation(&dbc)

	r := mux.NewRouter()
	r.HandleFunc("/dir/{extension}", phonebookHandler).Methods("GET")
	r.HandleFunc("/{extension}", menuHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8958", r))
}
