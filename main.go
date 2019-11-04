package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"web/models"
)

type env struct {
	db        *sql.DB
	templates *template.Template
}

func createUserHandler(env *env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := models.User{Username: "NullTest"}
		user.Create(env.db)
		user = models.User{Username: "Tony", Age: models.NewNullInt64("29")}
		if err := user.Create(env.db); err != nil {
			log.Printf("createUserHandler:\n\t%s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		var product models.Product
		product = models.Product{Name: "PS4", UserID: user.ID}
		product.Create(env.db)
		product = models.Product{Name: "SWITCH", UserID: user.ID}
		product.Create(env.db)

		http.Redirect(w, r, "/users", http.StatusSeeOther)
	})
}

func usersHandler(env *env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := models.GetAllUsers(env.db)
		if err != nil {
			log.Printf("usersHandler:\n\t%s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		env.templates.ExecuteTemplate(w, "users", users)
	})
}

func userHandler(env *env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := r.URL.Query().Get("id")
		id, _ := strconv.ParseInt(s, 10, 64)
		user, err := models.GetUser(id, env.db)
		if err != nil {
			log.Printf("usersHandler:\n\t%s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		products, _ := user.GetRelatedProducts(env.db)
		env.templates.ExecuteTemplate(w, "user", struct {
			User     *models.User
			Products []*models.Product
		}{user, products})
	})
}

func main() {
	env := &env{}

	var err error
	env.db, err = models.InitDB("db/test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer env.db.Close()

	env.templates = template.Must(template.ParseGlob("./templates/*"))

	// With gorilla mux
	// r := mux.NewRouter()
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/create_user", createUserHandler(env))
	http.Handle("/users", usersHandler(env))
	http.Handle("/detail", userHandler(env))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
