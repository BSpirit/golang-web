package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	ID       int64
	Username string
	Age      sql.NullInt64
}

func (u *User) Create(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO users(username, age) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("user.Create: could not prepare query\n\t%s", err)
	}
	res, err := stmt.Exec(u.Username, u.Age)
	if err != nil {
		return fmt.Errorf("user.Create: could not execute query\n\t%s", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("user.Create: could not retrieve last inserted id\n\t%s", err)
	}
	u.ID = id

	return nil
}

func (u *User) Update(db *sql.DB) error {
	if u.ID == 0 {
		return errors.New("user.Update: no ID")
	}
	stmt, err := db.Prepare("UPDATE users SET username=?, age=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("user.Update: could not prepare query\n\t%s", err)
	}
	_, err = stmt.Exec(u.Username, u.Age, u.ID)
	if err != nil {
		return fmt.Errorf("user.Update: could not execute query\n\t%s", err)
	}

	return nil
}

func (u *User) Delete(db *sql.DB) error {
	if u.ID == 0 {
		return errors.New("user.Delete: no ID")
	}
	stmt, err := db.Prepare("DELETE FROM users WHERE id=?")
	if err != nil {
		return fmt.Errorf("user.Delete: could not prepare query\n\t%s", err)
	}
	_, err = stmt.Exec(u.ID)
	if err != nil {
		return fmt.Errorf("user.Delete: could not execute query\n\t%s", err)
	}
	u.ID = 0

	return nil
}

func GetUser(pk int64, db *sql.DB) (*User, error) {
	rows, err := db.Query("SELECT * FROM users WHERE id = ?", pk)
	if err != nil {
		return nil, fmt.Errorf("GetUser: could not execute query\n\t%s", err)
	}
	defer rows.Close()
	user := &User{}
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Age, &user.Username)
		if err != nil {
			return nil, fmt.Errorf("GetUser: could not scan row\n\t%s", err)
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetUser: got error when fetching row\n\t%s", err)
	}

	return user, nil
}

func (u *User) GetRelatedProducts(db *sql.DB) ([]*Product, error) {
	rows, err := db.Query("SELECT * FROM products WHERE user_id=?", u.ID)
	if err != nil {
		return nil, fmt.Errorf("GetRelatedProducts: could not execute query\n\t%s", err)
	}
	defer rows.Close()
	products := make([]*Product, 0)
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.UserID)
		if err != nil {
			return nil, fmt.Errorf("GetRelatedProducts: could not scan row\n\t%s", err)
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetRelatedProducts: got error when fetching rows\n\t%s", err)
	}

	return products, nil
}

func GetAllUsers(db *sql.DB) ([]*User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: could not execute query\n\t%s", err)
	}
	defer rows.Close()
	users := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Age)
		if err != nil {
			return nil, fmt.Errorf("GetAllUsers: could not scan row\n\t%s", err)
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: got error when fetching rows\n\t%s", err)
	}

	return users, nil
}
