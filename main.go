package main

import (
	"BastetSoftware/backend/api"
	"BastetSoftware/backend/database"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "<h1>Homepage</h1>")
}

var db *sql.DB // TODO: make local

func CustomUnmarshal(data []byte, v interface{}) error {
	dec := msgpack.GetDecoder()

	dec.Reset(bytes.NewReader(data))
	dec.DisallowUnknownFields(true) // <- this is customized
	err := dec.Decode(v)

	msgpack.PutDecoder(dec)

	return err
}

func handleFPing(r *api.Request) (*api.Response, error) {
	return &api.Response{Code: 0, Data: nil}, nil
}

func handleFUserCreate(r *api.Request) (*api.Response, error) {
	// parse args
	var args api.ArgsFUserCreate
	err := CustomUnmarshal(r.Args, &args)
	if err != nil {
		return &api.Response{Code: api.EArgsInval, Data: nil}, nil
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		return &api.Response{Code: api.EUnknown, Data: nil}, nil
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
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == 1062 {
				return &api.Response{Code: api.EExists, Data: nil}, nil
			}
		}
		return &api.Response{Code: api.EUnknown, Data: nil}, nil
	}

	// success
	return &api.Response{Code: 0, Data: nil}, nil
}

func handleFLogIn(r *api.Request) (*api.Response, error) {
	// parse args
	var args api.ArgsFLogIn
	err := CustomUnmarshal(r.Args, &args)
	if err != nil {
		return &api.Response{Code: api.EArgsInval, Data: nil}, nil
	}

	session, err := database.OpenSession(db, args.Login, args.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return &api.Response{Code: api.ENoEntry, Data: nil}, nil
		} else if err.Error() == "incorrect password" {
			return &api.Response{Code: api.EPassWrong, Data: nil}, nil
		}
		return &api.Response{Code: api.EUnknown, Data: nil}, nil
	}

	var resp = api.RespFLogIn{
		Token: string(session.Token),
	}
	data, err := msgpack.Marshal(resp)
	if err != nil {
		// TODO: close session
		return &api.Response{Code: api.EUnknown, Data: nil}, nil
	}

	return &api.Response{Code: 0, Data: data}, nil
}

func handleFLogOut(r *api.Request) (*api.Response, error) {
	// parse args
	var args api.ArgsFLogOut
	err := CustomUnmarshal(r.Args, &args)
	if err != nil {
		return &api.Response{Code: api.EArgsInval, Data: nil}, nil
	}

	err = database.CloseSession(db, []byte(args.Token))
	if err != nil {
		return &api.Response{Code: api.EUnknown, Data: nil}, nil // TODO: proper error handling
	}

	return &api.Response{Code: 0, Data: nil}, err
}

var apiFHandlers [api.FNull]api.RequestHandler

func handleRequest(r *api.Request) (*api.Response, error) {
	if int(r.Func) >= len(apiFHandlers) {
		return &api.Response{Code: api.ENoFun, Data: nil}, nil
	}

	return apiFHandlers[r.Func](r)
}

func main() {
	var err error

	/* setup handlers */

	apiFHandlers[api.FPing] = handleFPing

	apiFHandlers[api.FUserCreate] = handleFUserCreate
	apiFHandlers[api.FLogIn] = handleFLogIn
	apiFHandlers[api.FLogOut] = handleFLogOut

	/* =(setup handlers)= */

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
