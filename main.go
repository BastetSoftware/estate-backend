package main

import (
	"BastetSoftware/backend/api"
	"BastetSoftware/backend/database"
	"bytes"
	"database/sql"
	"fmt"
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

func handleFReserved(r *api.Request) (*api.Response, error) {
	return &api.Response{Code: api.ENoFun, Data: nil}, nil
}

func handleFPing(r *api.Request) (*api.Response, error) {
	return &api.Response{Code: 0, Data: nil}, nil
}

func handleFUserCreate(r *api.Request) (*api.Response, error) {
	// parse args
	var args api.ArgsFUserCreate
	err := CustomUnmarshal(r.Args, &args)
	if err != nil {
		return &api.Response{Code: api.EArgsInval, Data: nil}, err
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	if args.Patronymic == "" {
		args.Patronymic = "-"
	}
	userInfo := database.UserInfo{
		Id:            0,
		Login:         args.Login,
		PassHash:      passHash,
		FirstName:     args.FirstName,
		LastName:      args.LastName,
		Patronymic:    args.Patronymic,
		ManagesGroups: false,
	}
	err = userInfo.Register(db)
	switch err {
	case nil:
		break
	case database.ErrUserExists:
		return &api.Response{Code: api.EExists, Data: nil}, nil
	default:
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	// success
	return &api.Response{Code: 0, Data: nil}, nil
}

func handleFLogIn(r *api.Request) (*api.Response, error) {
	// parse args
	var args api.ArgsFLogIn
	err := CustomUnmarshal(r.Args, &args)
	if err != nil {
		return &api.Response{Code: api.EArgsInval, Data: nil}, err
	}

	session, err := database.OpenSession(db, args.Login, args.Password)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return &api.Response{Code: api.ENoEntry, Data: nil}, nil
	case database.ErrPassWrong:
		return &api.Response{Code: api.EPassWrong, Data: nil}, nil
	default:
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	var resp = api.RespFLogIn{
		Token: string(session.Token),
	}
	data, err := msgpack.Marshal(resp)
	if err != nil {
		// TODO: close session
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	return &api.Response{Code: 0, Data: data}, nil
}

func handleFLogOut(r *api.Request) (*api.Response, error) {
	// parse args
	var args api.ArgsFLogOut
	err := CustomUnmarshal(r.Args, &args)
	if err != nil {
		return &api.Response{Code: api.EArgsInval, Data: nil}, err
	}

	err = database.CloseSession(db, []byte(args.Token))
	if err != nil {
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	return &api.Response{Code: 0, Data: nil}, nil
}

func handleFUserInfo(r *api.Request) (*api.Response, error) {
	// parse args
	var args api.ArgsFUserInfo
	err := CustomUnmarshal(r.Args, &args)
	if err != nil {
		return &api.Response{Code: api.EArgsInval, Data: nil}, err
	}

	_, err = database.VerifySession(db, []byte(args.Token))
	switch err {
	case nil:
		break
	case database.ErrNotLoggedIn:
		return &api.Response{Code: api.ENotLoggedIn, Data: nil}, nil
	default:
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	userinfo, err := database.FindUserInfo(db, args.Login)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return &api.Response{Code: api.ENoEntry, Data: nil}, nil
	default:
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	resp := api.RespFUserInfo{
		Login:      userinfo.Login,
		FirstName:  userinfo.FirstName,
		LastName:   userinfo.LastName,
		Patronymic: userinfo.Patronymic,
	}
	data, err := msgpack.Marshal(resp)
	if err != nil {
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	return &api.Response{Code: 0, Data: data}, nil
}

func handleFUserEdit(r *api.Request) (*api.Response, error) {
	// parse args
	var args api.ArgsFUserEdit
	err := msgpack.Unmarshal(r.Args, &args)
	if err != nil || args.Token == "" {
		return &api.Response{Code: api.EArgsInval, Data: nil}, err
	}

	session, err := database.VerifySession(db, []byte(args.Token))
	switch err {
	case nil:
		break
	case database.ErrNotLoggedIn:
		return &api.Response{Code: api.ENotLoggedIn, Data: nil}, nil
	default:
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	uid := session.User

	if args.Login != nil {
		err = database.UserChangeLogin(db, uid, *args.Login)
		switch err {
		case nil:
			break
		case database.ErrNoUser:
			return &api.Response{Code: api.ENoEntry, Data: nil}, nil
		case database.ErrUserExists:
			return &api.Response{Code: api.EExists, Data: nil}, nil
		default:
			return &api.Response{Code: api.EUnknown, Data: nil}, err
		}
	}

	if args.Password != nil {
		passHash, err := bcrypt.GenerateFromPassword([]byte(*args.Password), bcrypt.DefaultCost)
		if err != nil {
			return &api.Response{Code: api.EUnknown, Data: nil}, err
		}
		err = database.UserChangePasswordHash(db, uid, passHash)
		switch err {
		case nil:
			break
		case database.ErrNoUser:
			return &api.Response{Code: api.ENoEntry, Data: nil}, nil
		default:
			return &api.Response{Code: api.EUnknown, Data: nil}, err
		}
	}

	var names = [3]*string{args.FirstName, args.LastName, args.Patronymic}
	for i, n := range names {
		if n == nil {
			continue
		}

		err = database.UserChangeName(db, uid, i, *n)
		switch err {
		case nil:
			break
		case database.ErrNoUser:
			return &api.Response{Code: api.ENoEntry, Data: nil}, nil
		default:
			fmt.Println(err)
			return &api.Response{Code: api.EUnknown, Data: nil}, err
		}
	}

	return &api.Response{Code: 0, Data: nil}, nil
}

// verifyManagesGroups: check that user can manage groups
func verifyManagesGroups(token string) (*api.Response, error) {
	session, err := database.VerifySession(db, []byte(token))
	switch err {
	case nil:
		break
	case database.ErrNotLoggedIn:
		return &api.Response{Code: api.ENotLoggedIn, Data: nil}, nil
	default:
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	userinfo, err := database.GetUserInfo(db, session.User)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return &api.Response{Code: api.ENoEntry, Data: nil}, nil
	default:
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	if !userinfo.ManagesGroups {
		return &api.Response{Code: api.EAccessDenied, Data: nil}, nil
	}

	return nil, nil
}

func handleFGroupCreate(r *api.Request) (*api.Response, error) {
	// parse args
	var args api.ArgsFGroupCreateRemove
	err := CustomUnmarshal(r.Args, &args)
	if err != nil {
		return &api.Response{Code: api.EArgsInval, Data: nil}, err
	}

	// check that user can manage groups
	resp, err := verifyManagesGroups(args.Token)
	if resp != nil {
		return resp, err
	}

	_, err = database.CreateGroup(db, args.Name)
	switch err {
	case nil:
		break
	case database.ErrGroupExists:
		return &api.Response{Code: api.EExists, Data: nil}, nil
	default:
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	return &api.Response{Code: 0, Data: nil}, nil
}

func handleFGroupRemove(r *api.Request) (*api.Response, error) {
	// parse args
	var args api.ArgsFGroupCreateRemove
	err := CustomUnmarshal(r.Args, &args)
	if err != nil {
		return &api.Response{Code: api.EArgsInval, Data: nil}, err
	}

	// check that user can manage groups
	resp, err := verifyManagesGroups(args.Token)
	if resp != nil {
		return resp, err
	}

	group, err := database.FindGroup(db, args.Name)
	switch err {
	case nil:
		break
	case database.ErrNoGroup:
		return &api.Response{Code: api.ENoEntry, Data: nil}, nil
	default:
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	err = database.RemoveGroup(db, group.Id)
	switch err {
	case nil:
		break
	case database.ErrNoGroup:
		return &api.Response{Code: api.ENoEntry, Data: nil}, nil
	default:
		return &api.Response{Code: api.EUnknown, Data: nil}, err
	}

	return &api.Response{Code: 0, Data: nil}, nil
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
	apiFHandlers[api.FUserInfo] = handleFUserInfo
	apiFHandlers[api.FUserEdit] = handleFUserEdit

	apiFHandlers[api.FGroupCreate] = handleFGroupCreate
	apiFHandlers[api.FGroupRemove] = handleFGroupRemove

	/* =(setup handlers)= */

	db, err = database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	// err = api.Listen("/tmp/estate.sock", handleRequest)
	err = api.Listen("localhost:8080", handleRequest)
	if err != nil {
		log.Fatal(err)
	}

	//http.HandleFunc("/", rootHandler)
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
