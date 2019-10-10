package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"web/models"
)

var t *template.Template
var db *sql.DB

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	user := models.User{Username: "Tony"}
	if err := user.Create(db); err != nil {
		log.Printf("createUserHandler:\n\t%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		t.ExecuteTemplate(w, "error", nil)
		return
	}

	// user.Username = "TEST"
	// if err := user.Update(db); err != nil {
	// 	log.Printf("createUserHandler:\n\t%s", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	t.ExecuteTemplate(w, "error", nil)
	// 	return
	// }

	// if err := user.Delete(db); err != nil {
	// 	log.Printf("createUserHandler:\n\t%s", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	t.ExecuteTemplate(w, "error", nil)
	// 	return
	// }

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers(db)
	if err != nil {
		log.Printf("usersHandler:\n\t%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		t.ExecuteTemplate(w, "error", nil)
		return
	}

	t.ExecuteTemplate(w, "users", users)
}

func main() {
	var err error
	db, err = models.InitDB("db/test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	t = template.Must(template.ParseGlob("./templates/*"))

	// With gorilla mux
	// r := mux.NewRouter()
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/create_user", createUserHandler)
	http.HandleFunc("/users", usersHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
