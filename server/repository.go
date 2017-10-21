package server

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"fmt"
	"errors"
)

var connection *sql.DB

type EmailRow struct {
	email string
	hasPassword bool
}

func connect(host string, port int, user string, password string, database string) *sql.DB {
	var connectString string = user + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + database

	fmt.Println("Try connect to the database: " + connectString)

	db, err := sql.Open("mysql", connectString)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to the database: " + connectString)

	return db
}

func insertEmail(emailAddress string) bool {
	stmtIns, err := connection.Prepare("INSERT INTO emails (email) VALUES(?)")
	defer stmtIns.Close()

	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = stmtIns.Exec(emailAddress)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func getEmails(offset uint, limit uint) ([]EmailRow, error) {
	rows, err := connection.Query("SELECT email, password != '' FROM emails ORDER BY email LIMIT ?, ?", offset, limit)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("database error")
	}

	var emailRows []EmailRow = []EmailRow{}

	for rows.Next() {
		var emailRow EmailRow
		rows.Scan(&emailRow.email, &emailRow.hasPassword)

		emailRows = append(emailRows, emailRow)
	}

	return emailRows, nil
}

func countAllEmails() (uint, error) {
	stmtOut, err := connection.Prepare("SELECT count(email) FROM emails")
	var count uint
	err = stmtOut.QueryRow().Scan(&count)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("database error")
	}

	return count, nil
}
