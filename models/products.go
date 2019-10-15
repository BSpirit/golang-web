package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Product struct {
	ID     int64
	Name   string
	UserID int64
}

func (p *Product) Create(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO products(name, user_id) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("product.Create: could not prepare query\n\t%s", err)
	}
	res, err := stmt.Exec(p.Name, p.UserID)
	if err != nil {
		return fmt.Errorf("product.Create: could not execute query\n\t%s", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("product.Create: could not retrieve last inserted id\n\t%s", err)
	}
	p.ID = id

	return nil
}

func (p *Product) Update(db *sql.DB) error {
	if p.ID == 0 {
		return errors.New("product.Update: no ID")
	}
	stmt, err := db.Prepare("UPDATE products SET name=?, user_id=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("product.Update: could not prepare query\n\t%s", err)
	}
	_, err = stmt.Exec(p.Name, p.UserID, p.ID)
	if err != nil {
		return fmt.Errorf("product.Update: could not execute query\n\t%s", err)
	}

	return nil
}

func (p *Product) Delete(db *sql.DB) error {
	if p.ID == 0 {
		return errors.New("product.Delete: no ID")
	}
	stmt, err := db.Prepare("DELETE FROM products WHERE id=?")
	if err != nil {
		return fmt.Errorf("product.Delete: could not prepare query\n\t%s", err)
	}
	_, err = stmt.Exec(p.ID)
	if err != nil {
		return fmt.Errorf("product.Delete: could not execute query\n\t%s", err)
	}
	p.ID = 0

	return nil
}

func GetProduct(pk int64, db *sql.DB) (*Product, error) {
	rows, err := db.Query("SELECT * FROM products WHERE id = ?", pk)
	if err != nil {
		return nil, fmt.Errorf("GetProduct: could not execute query\n\t%s", err)
	}
	defer rows.Close()
	product := &Product{}
	for rows.Next() {
		err := rows.Scan(&product.ID, &product.Name, &product.UserID)
		if err != nil {
			return nil, fmt.Errorf("GetProduct: could not scan row\n\t%s", err)
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetProduct: got error when fetching row\n\t%s", err)
	}

	return product, nil
}

func GetAllProducts(db *sql.DB) ([]*Product, error) {
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		return nil, fmt.Errorf("GetAllProducts: could not execute query\n\t%s", err)
	}
	defer rows.Close()
	products := make([]*Product, 0)
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.UserID)
		if err != nil {
			return nil, fmt.Errorf("GetAllProducts: could not scan row\n\t%s", err)
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetAllProducts: got error when fetching rows\n\t%s", err)
	}

	return products, nil
}
