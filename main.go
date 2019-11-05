package main

import (
	"html/template"
	"log"
	"net/http"
	"web/handlers"
	"web/models"
)

func main() {
	env := &handlers.Env{}
	var err error

	env.Templates = template.Must(template.ParseGlob("./templates/*"))
	env.DB, err = models.InitDB("db/test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer env.DB.Close()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/fill_db", handlers.Handler{
		Env:         env,
		HandlerFunc: handlers.FillDBHandler,
	})

	http.Handle("/create_user", handlers.Handler{
		Env:         env,
		HandlerFunc: handlers.CreateUserHandler,
	})

	http.Handle("/users", handlers.Handler{
		Env:         env,
		HandlerFunc: handlers.UsersHandler,
	})

	http.Handle("/user", handlers.Handler{
		Env:         env,
		HandlerFunc: handlers.UserHandler,
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
