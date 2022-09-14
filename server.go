package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
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

	http.HandleFunc("/", Home)
	http.HandleFunc("/registration", Registration)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)

	http.ListenAndServe(":8080", nil)

}

type Hedgehog struct {
	Id          int
	Name        string
	Description string
}

func Home(w http.ResponseWriter, r *http.Request) {

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
	templates.ExecuteTemplate(w, "home", hedgeArray)
}

func Registration(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "registration", nil)
}

func Create(w http.ResponseWriter, r *http.Request) {
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

func Delete(w http.ResponseWriter, r *http.Request) {
	hedgeId := r.URL.Query().Get("ID")
	hedgeIdInt, err1 := strconv.Atoi(hedgeId)
	if err1 != nil {
		panic(err1.Error())
	}
	establishedConn := connBD()
	_, err2 := establishedConn.Exec("DELETE FROM hedgehogs WHERE id= $1", hedgeIdInt)
	if err2 != nil {
		panic(err2.Error())
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	hedgeId := r.URL.Query().Get("id")
	hedgeIdInt, err1 := strconv.Atoi(hedgeId)
	if err1 != nil {
		panic(err1.Error())
	}
	establishedConn := connBD()
	hedgeItem, err2 := establishedConn.Query("SELECT * FROM hedgehogs WHERE id=$1", hedgeIdInt)
	if err2 != nil {
		panic(err2.Error())
	}
	hedgehog := Hedgehog{}
	for hedgeItem.Next() {
		var id int
		var name, description string
		err2 = hedgeItem.Scan(&id, &name, &description)
		if err2 != nil {
			panic(err2.Error())
		}
		hedgehog.Id = id
		hedgehog.Name = name
		hedgehog.Description = description
	}
	templates.ExecuteTemplate(w, "edit", hedgehog)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		idInt, err1 := strconv.Atoi(id)
		if err1 != nil {
			panic(err1.Error())
		}
		name := r.FormValue("name")
		description := r.FormValue("description")
		establishedConn := connBD()
		_, err2 := establishedConn.Exec("UPDATE hedgehogs SET name=$1, description=$2 WHERE id=$3", name, description, idInt)
		if err2 != nil {
			fmt.Print("Error de insercion")
			panic(err2.Error())
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}
