package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"

	"database/sql"
	"github.com/go-sql-driver/mysql"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "<h1>Homepage</h1>")
}

func main() {
	var err error
	var db *sql.DB

	db, err = dbOpen()
	if err != nil {
		log.Fatal(err)
	}

	pass, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	admin := UserInfo{
		Id:         0,
		Login:      "admin",
		PassHash:   pass,
		FirstName:  "Admin",
		LastName:   "McServer",
		Patronymic: nil,
		Role:       RoleAdmin,
	}
	err = admin.Register(db)
	if err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == 1062 {
				log.Println(err)
			} else {
				log.Fatal(err)
			}
		default:
			log.Fatal(err)
		}
	}

	session, err := openSession(db, "admin", "12345678")
	if err != nil {
		log.Fatal(err)
	}

	info, err := getUserInfo(db, session.User)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v: %s\n", session.Id, info.Format())

	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
