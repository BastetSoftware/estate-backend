package main

import (
	"BastetSoftware/backend/api"
	"BastetSoftware/backend/database"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "<h1>Homepage</h1>")
}

func main() {
	var err error
	var db *sql.DB

	db, err = database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	pass, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	admin := database.UserInfo{
		Id:         0,
		Login:      "admin",
		PassHash:   pass,
		FirstName:  "Admin",
		LastName:   "McServer",
		Patronymic: nil,
		Role:       database.RoleAdmin,
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

	session, err := database.OpenSession(db, "admin", "12345678")
	if err != nil {
		log.Fatal(err)
	}

	info, err := database.GetUserInfo(db, session.User)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v: %s\n", session.Id, info.Format())

	err = api.Listen("/tmp/estate.sock", func(conn *net.Conn, r *api.Request) {
		fmt.Printf("%d: %v\n", r.Func, r.Args)
	})
	if err != nil {
		log.Fatal()
	}

	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
