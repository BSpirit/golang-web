package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewNullInt64(s string) sql.NullInt64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: n,
		Valid: true,
	}
}

func WhereClause(entries map[string][]string) (string, []interface{}) {
	var whereClause string
	var values []interface{}

	var clauseBuilder []string
	for key, value := range entries {
		if value[0] != "" {
			values = append(values, value[0])
			clauseBuilder = append(clauseBuilder, fmt.Sprintf("%s LIKE '%%' || ? || '%%'", key))
		}
	}

	if len(clauseBuilder) > 0 {
		whereClause = " WHERE " + strings.Join(clauseBuilder, " AND ")
	}

	return whereClause, values
}
