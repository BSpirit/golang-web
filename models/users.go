package models

import (
	"database/sql"
	"errors"
	"fmt"
	"web/utils"
)

type User struct {
	ID       int64
	Username string
	Age      sql.NullInt64
}

func (u *User) Create(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO users(username, age) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf(utils.Trace(err))
	}
	res, err := stmt.Exec(u.Username, u.Age)
	if err != nil {
		return fmt.Errorf(utils.Trace(err))
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf(utils.Trace(err))
	}
	u.ID = id

	return nil
}

func (u *User) Update(db *sql.DB) error {
	if u.ID == 0 {
		return errors.New("no ID")
	}
	stmt, err := db.Prepare("UPDATE users SET username=?, age=? WHERE id=?")
	if err != nil {
		return fmt.Errorf(utils.Trace(err))
	}
	_, err = stmt.Exec(u.Username, u.Age, u.ID)
	if err != nil {
		return fmt.Errorf(utils.Trace(err))
	}

	return nil
}

func (u *User) Delete(db *sql.DB) error {
	if u.ID == 0 {
		return errors.New("no ID")
	}
	stmt, err := db.Prepare("DELETE FROM users WHERE id=?")
	if err != nil {
		return fmt.Errorf(utils.Trace(err))
	}
	_, err = stmt.Exec(u.ID)
	if err != nil {
		return fmt.Errorf(utils.Trace(err))
	}
	u.ID = 0

	return nil
}

func GetUser(pk int64, db *sql.DB) (*User, error) {
	user := &User{}
	err := db.QueryRow("SELECT * FROM users WHERE id = ?", pk).Scan(&user.ID, &user.Username, &user.Age)
	if err != nil {
		return nil, fmt.Errorf(utils.Trace(err))
	}

	return user, nil
}

func (u *User) GetRelatedProducts(db *sql.DB) ([]*Product, error) {
	rows, err := db.Query("SELECT * FROM products WHERE user_id=?", u.ID)
	if err != nil {
		return nil, fmt.Errorf(utils.Trace(err))
	}
	defer rows.Close()
	products := make([]*Product, 0)
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.UserID)
		if err != nil {
			return nil, fmt.Errorf(utils.Trace(err))
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf(utils.Trace(err))
	}

	return products, nil
}

func GetUsersByFilter(entries map[string][]string, db *sql.DB) ([]*User, error) {
	whereClause, values := WhereClause(entries)
	rows, err := db.Query("SELECT * FROM users"+whereClause, values...)
	if err != nil {
		return nil, fmt.Errorf(utils.Trace(err))
	}
	defer rows.Close()
	users := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Age)
		if err != nil {
			return nil, fmt.Errorf(utils.Trace(err))
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf(utils.Trace(err))
	}

	return users, nil
}
