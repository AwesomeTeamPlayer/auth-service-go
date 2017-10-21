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

type RegisterRequest struct {
	Email string
	Password string
}

type LoginRequest struct {
	Email string
	Password string
	Label string
}

type LogoutRequest struct {
	Email string
	SessionKey string
}

type PaginationRequest struct {
	Page uint
	Limit uint
}

type EmailRowsResponse struct {
	EmailRows []EmailRow
	Count uint
}

type LoggedEmailsResponse struct {
	LoggedEmails []string
	Count uint
}

type App int

func (t *App) CreateUser (r *http.Request, request *EmailRequest, result *bool) error {
	fmt.Println("CreateUser")
	isInserted := insertEmail(request.Email)

	if isInserted {
		userCreated(request.Email)
	}

	fmt.Println("User created")

	*result = isInserted
	return nil
}

func (t *App) Register (r *http.Request, request *RegisterRequest, result *bool) error {
	fmt.Println("Register")
	isInserted := setPassword(request.Email, hashPassword(request.Password))

	if isInserted {
		userCreated(request.Email)
	}

	fmt.Println("User registered")

	*result = isInserted
	return nil
}

func (t *App) Login (r *http.Request, request *LoginRequest, result *string) error {
	fmt.Println("Login")

	*result = ""

	if request.Password == "" {
		return nil
	}

	hashedPassword, err := getHashedPassword(request.Email)
	if err != nil {
		return err
	}

	if hashedPassword != hashPassword(request.Password) {
		return nil
	}

	sessionKey := generateRandomHash(100)
	isInserted := insertSession(request.Email, sessionKey, request.Label);

	if isInserted {
		userLoggedIn(request.Email)
	}

	fmt.Println("Logged in")

	*result = sessionKey
	return nil
}

func generateRandomHash(length uint) string {
	return "asd"
}

func hashPassword(password string) string {
	return password
}

func (t *App) Logout (r *http.Request, request *LogoutRequest, result *bool) error {
	fmt.Println("Logout")

	isRemoved := removeSession(request.Email, request.SessionKey)

	if isRemoved {
		userLoggedOut(request.Email)
	}

	fmt.Println("Logged out")

	*result = isRemoved
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

func (t *App) GetLoggedUsers (r *http.Request, request *PaginationRequest, result *LoggedEmailsResponse) error {
	fmt.Println("GetLoggedUsers")
	limit := uint(request.Limit)

	emails := getLoggedEmails(uint(request.Page * limit), limit)
	count := countAllLoggedEmails()

	*result = LoggedEmailsResponse{emails, count}
	return nil
}

func (t *App) GetSessions (r *http.Request, request *EmailRequest, result []SessionRow) error {
	fmt.Println("GetSessions")

	sessionsRows, err := getSessionsRows(request.Email)
	if err != nil {
		return errors.New("database error")
	}

	result = sessionsRows
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