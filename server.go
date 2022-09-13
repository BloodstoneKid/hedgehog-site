package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "andro747"
	dbname   = "hedgehogdb"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func connBD() (conn *sql.DB) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err.Error())
	}
	return conn
}

func main() {

	http.HandleFunc("/", home)
	http.HandleFunc("/registration", registration)
	http.HandleFunc("/create", create)
	http.HandleFunc("/delete", delete)

	http.ListenAndServe(":8080", nil)

}

func home(w http.ResponseWriter, r *http.Request) {

	type Hedgehog struct {
		Id          int
		Name        string
		Description string
	}

	establishedConn := connBD()
	hedgeList, err := establishedConn.Query("SELECT * FROM hedgehogs")
	if err != nil {
		panic(err.Error())
	}
	hedgehog := Hedgehog{}
	hedgeArray := []Hedgehog{}
	for hedgeList.Next() {
		var id int
		var name, description string
		err = hedgeList.Scan(&id, &name, &description)
		if err != nil {
			panic(err.Error())
		}
		hedgehog.Id = id
		hedgehog.Name = name
		hedgehog.Description = description
		hedgeArray = append(hedgeArray, hedgehog)
	}
	templates.ExecuteTemplate(w, "home", nil)
}

func registration(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "registration", nil)
}

func create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		description := r.FormValue("description")
		establishedConn := connBD()
		_, err := establishedConn.Exec("INSERT INTO hedgehogs(name,description) VALUES($1,$2)", name, description)
		if err != nil {
			fmt.Print("Error de insercion")
			panic(err.Error())
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	hedgeId := r.URL.Query().Get("id")

	establishedConn := connBD()
	_, err := establishedConn.Exec("DELETE FROM hedgehogs WHERE id= $1", hedgeId)
	if err != nil {
		panic(err.Error())
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
