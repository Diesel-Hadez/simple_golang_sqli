package main
import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

import (
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	host = "db"
	port = 5432
	user = "user"
	password = "password"
	dbname = "users"
)

var db *sql.DB

func getLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /login request\n")
	fmt.Println("GET params were: ", r.URL.Query())
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	var counter int
	queryString := fmt.Sprintf("SELECT count(*) FROM users WHERE username='%s' AND pass='%s'", username, password)
	db.QueryRow(queryString).Scan(&counter)
	fmt.Printf("Query: \" %s \" ", queryString)
	fmt.Printf("found only %d\n", counter)
	if counter == 1 {
		io.WriteString(w, "Admin portal is on TODO list to be implemented. Please see: Jira ticket TPE-3174.\n")
	} else {
		io.WriteString(w, "User not found!\n")
	}
}

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to db!")
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	db.Exec(`DROP TABLE users`)
	db.Exec(`CREATE TABLE users (username TEXT, pass TEXT)`)
	insertDynStmt := `insert into "users"("username", "pass") values($1, $2)`
	_, err = db.Exec(insertDynStmt, "admin", "CTF{SAMPLE_FLAG}")
	if err != nil {
		panic(err)
	}


	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/login", getLogin)

	err = http.ListenAndServe(":8080", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
		defer db.Close()
	}	else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		defer db.Close()
		os.Exit(1)
	}
}
