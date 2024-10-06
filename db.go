package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "dacidon"
	dbname = "tgBot"
)

func connectString() string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"dbname=%s sslmode=disable", host, port, user, dbname)

	return psqlInfo
}

func AddUser(user_id int64, username string, group_id string) {

	db, err := sql.Open("postgres", connectString())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	sqlStatement := `
		INSERT INTO users (user_id, username, group_id)
		VALUES ($1, $2, $3)`

	_, err = db.Exec(sqlStatement, user_id, username, group_id)
	if err != nil {
		panic(err)
	}
}

func UpdateUser(user_id int64, group_id string) {

	db, err := sql.Open("postgres", connectString())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	sqlStatement := `
		UPDATE users
		SET group_id = $2
		WHERE user_id = $1;`
	_, err = db.Exec(sqlStatement, user_id, group_id)
	if err != nil {
		panic(err)
	}
}

func RetrieveGroup(user_id int64) string {

	db, err := sql.Open("postgres", connectString())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var group_id string
	sqlStatement := `
		SELECT group_id FROM users
		WHERE user_id = $1;`

	row := db.QueryRow(sqlStatement, user_id)

	switch err := row.Scan(&group_id); err {
	case sql.ErrNoRows:
		return ""
	case nil:
		return group_id
	default:
		panic(err)
	}
}

func CheckUser(user_id int64) bool {

	db, err := sql.Open("postgres", connectString())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	sqlStatement := `
		SELECT * FROM users
		WHERE user_id = $1`

	row := db.QueryRow(sqlStatement, user_id)

	var (
		i, u, g string
	)

	switch err := row.Scan(&i, &u, &g); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		panic(err)
	}
}
