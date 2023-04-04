package main

import (
	"BastetSoftware/backend/api"
	"BastetSoftware/backend/database"
	"database/sql"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "<h1>Homepage</h1>")
}

var db *sql.DB // TODO: make local

func handleRequest(conn net.Conn, r *api.Request) (*api.Response, error) {
	switch r.Func {
	case api.FPing:
		response := api.Response{
			Code: 0,
			Data: nil,
		}

		return &response, nil

	case api.FUserCreate:
		var response api.Response

		var args api.ArgsFUserCreate
		err := msgpack.Unmarshal(r.Args, &args)
		if err != nil {
			response = api.Response{
				Code: 1,
				Data: nil,
			}
			return &response, nil
		}

		passHash, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
		if err != nil {
			response = api.Response{
				Code: 1,
				Data: nil,
			}
			return &response, nil
		}

		var patr *string = nil
		if args.Patronymic != "" {
			patr = &args.Patronymic // TODO: this is unsafe
		}
		userInfo := database.UserInfo{
			Id:         0,
			Login:      args.Login,
			PassHash:   passHash,
			FirstName:  args.FirstName,
			LastName:   args.LastName,
			Patronymic: patr,
			Role:       database.RoleUser,
		}
		err = userInfo.Register(db)
		if err != nil {
			response = api.Response{
				Code: 1,
				Data: nil,
			}
			return &response, nil
		}

		response = api.Response{
			Code: 0,
			Data: nil,
		}
		return &response, nil
	}

	response := api.Response{
		Code: 255,
		Data: nil,
	}
	return &response, nil
}

func main() {
	var err error

	db, err = database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	err = api.Listen("/tmp/estate.sock", handleRequest)
	if err != nil {
		log.Fatal()
	}

	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
