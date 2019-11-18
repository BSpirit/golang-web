package handlers

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"web/models"
	"web/utils"
)

func FillDBHandler(env *Env, w http.ResponseWriter, r *http.Request) *StatusError {
	csvFile, err := os.Open("data.csv")
	if err != nil {
		return &StatusError{Code: 500, Err: utils.Trace(err)}
	}

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return &StatusError{Code: 500, Err: utils.Trace(err)}
	}

	tx, err := env.DB.Begin()
	if err != nil {
		return &StatusError{Code: 500, Err: utils.Trace(err)}
	}
	defer tx.Rollback()

	for _, record := range records {
		user := &models.User{
			Username: record[0],
			Age:      models.NewNullInt64(record[1]),
		}

		err := user.Create(tx)
		if err != nil {
			return &StatusError{Code: 500, Err: utils.Trace(err)}
		}

		products := make([]*models.Product, 0)
		err = json.Unmarshal([]byte(record[2]), &products)
		if err != nil {
			return &StatusError{Code: 500, Err: utils.Trace(err)}
		}

		for _, product := range products {
			product.UserID = user.ID
			err := product.Create(tx)
			if err != nil {
				return &StatusError{Code: 500, Err: utils.Trace(err)}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return &StatusError{Code: 500, Err: utils.Trace(err)}
	}

	return nil
}

func CreateUserHandler(env *Env, w http.ResponseWriter, r *http.Request) *StatusError {
	if r.Method == http.MethodGet {
		err := env.Templates.ExecuteTemplate(w, "create_user", nil)
		if err != nil {
			return &StatusError{Code: 500, Err: utils.Trace(err)}
		}

	} else if r.Method == http.MethodPost {
		user := models.User{
			Username: r.FormValue("username"),
			Age:      models.NewNullInt64(r.FormValue("age")),
		}

		tx, err := env.DB.Begin()
		if err != nil {
			return &StatusError{Code: 500, Err: utils.Trace(err)}
		}
		defer tx.Rollback()

		err = user.Create(tx)
		if err != nil {
			return &StatusError{Code: 500, Err: utils.Trace(err)}
		}

		err = tx.Commit()
		if err != nil {
			return &StatusError{Code: 500, Err: utils.Trace(err)}
		}

		http.Redirect(w, r, "/users", http.StatusSeeOther)
	}

	return nil
}

func UsersHandler(env *Env, w http.ResponseWriter, r *http.Request) *StatusError {
	err := r.ParseForm()
	if err != nil {
		return &StatusError{Code: 500, Err: utils.Trace(err)}
	}

	users, err := models.GetUsersByFilter(r.Form, env.DB)
	if err != nil {
		return &StatusError{Code: 500, Err: utils.Trace(err)}
	}

	err = env.Templates.ExecuteTemplate(w, "filter_users", users)
	if err != nil {
		return &StatusError{Code: 500, Err: utils.Trace(err)}
	}

	return nil
}

func UserHandler(env *Env, w http.ResponseWriter, r *http.Request) *StatusError {
	s := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return &StatusError{Code: 500, Err: utils.Trace(err)}
	}

	user, err := models.GetUser(id, env.DB)
	if err != nil {
		return &StatusError{Code: 500, Err: utils.Trace(err)}
	}

	products, err := user.GetRelatedProducts(env.DB)
	if err != nil {
		return &StatusError{Code: 500, Err: utils.Trace(err)}
	}

	err = env.Templates.ExecuteTemplate(w, "user", struct {
		User     *models.User
		Products []*models.Product
	}{user, products})
	if err != nil {
		return &StatusError{Code: 500, Err: utils.Trace(err)}
	}

	return nil
}
