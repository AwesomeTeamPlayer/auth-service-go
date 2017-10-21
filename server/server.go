package server

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"net/http"
	"os"
	"strconv"
	"fmt"
	"errors"
)

type EmailRequest struct {
	Email string
}

type PaginationRequest struct {
	Page uint
	Limit uint
}

type EmailRowsResponse struct {
	EmailRows []EmailRow
	Count uint
}

type App int

func (t *App) Register (r *http.Request, request *EmailRequest, result *bool) error {
	fmt.Println("Register")
	isInserted := insertEmail(request.Email)

	if isInserted {
		userCreated(request.Email)
	}

	fmt.Println("Email registered")

	*result = isInserted
	return nil
}

func (t *App) GetEmails (r *http.Request, request *PaginationRequest, result *EmailRowsResponse) error {
	fmt.Println("GetEmails")
	limit := uint(request.Limit)
	emails, err := getEmails(uint(request.Page * limit), limit)
	if err != nil {
		return errors.New("database error")
	}

	count, err := countAllEmails()
	if err != nil {
		return errors.New("database error")
	}

	*result = EmailRowsResponse{emails, count}
	return nil
}

func StartServer() {
	port, _ := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	appPort := os.Getenv("APP_PORT")

	connection = connect(
		os.Getenv("MYSQL_HOST"),
		port,
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_DATABASE"),
	)

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	app := new(App)
	s.RegisterService(app, "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)

	fmt.Println("Server has started on port " + appPort)
	http.ListenAndServe(":" + appPort, r)
}