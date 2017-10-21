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
	Email string
	HasPassword bool
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

func insertSession(emailAddress string, sessionKey string, label string) bool {
	id := findEmailId(emailAddress)
	if id == 0 {
		return false
	}

	stmtIns, err := connection.Prepare("INSERT INTO sessions (email_id, session_key, label) VALUES(?, ?, ?)")
	defer stmtIns.Close()

	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = stmtIns.Exec(id, sessionKey, label)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func removeSession(emailAddress string, sessionKey string) bool {
	emailId := findEmailId(emailAddress)
	if emailId == 0 {
		return false
	}

	var id uint = 0
	err := connection.QueryRow("SELECT id FROM sessions WHERE email_id = ? AND session_key=?", emailId, sessionKey).Scan(&id)
	if err != nil {
		return false
	}

	if id == 0 {
		return false
	}

	_, err = connection.Exec("DELETE FROM sessions WHERE id = ?", id)
	if err != nil {
		return false
	}

	return true
}

func setPassword(emailAddress string, hashedPassword string) bool {
	id := findEmailId(emailAddress)
	if id == 0 {
		return false
	}

	stmtIns, err := connection.Prepare("UPDATE emails SET password = ? WHERE id = ?")
	defer stmtIns.Close()

	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = stmtIns.Exec(hashedPassword, id)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func findEmailId(emailAddress string) uint {
	query, err := connection.Prepare("SELECT id FROM emails WHERE email=?")
	if err != nil {
		fmt.Println(err)
		return 0
	}

	var id uint
	err = query.QueryRow(emailAddress).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return id
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
		rows.Scan(&emailRow.Email, &emailRow.HasPassword)

		emailRows = append(emailRows, emailRow)
	}

	return emailRows, nil
}

func getHashedPassword(email string) (string, error) {
	stmtOut, err := connection.Prepare("SELECT password FROM emails WHERE email = ?")
	if err != nil {
		fmt.Println(err)
		return "", errors.New("database error")
	}

	var password string;
	err = stmtOut.QueryRow(email).Scan(&password)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("database error")
	}

	return password, nil
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
