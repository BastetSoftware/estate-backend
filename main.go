package main

import (
	"BastetSoftware/backend/api"
	"BastetSoftware/backend/database"
	"bytes"
	"database/sql"
	"io"
	"log"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/bcrypt"
)

func writeResponse(w http.ResponseWriter, v interface{}) error {
	data, err := msgpack.Marshal(v)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	var buf []byte
	var n int
	handler := apiFHandlers[r.URL.Path[len("/api/"):]]
	if handler == nil {
		// unknown API function, no arguments
		handler = unknownFPlug
		buf = nil
		n = 0
	} else {
		// valid API function, read request body
		buf = make([]byte, 4096)
		var err error
		n, err = r.Body.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}

	response, err := handler(buf[:n])
	if err != nil {
		log.Println(err)
	}

	err = writeResponse(w, response)
	if err != nil {
		log.Println(err)
	}
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

func unknownFPlug(_ []byte) (interface{}, error) {
	return api.Response{Code: api.ENoFun}, nil
}

func handleFPing(_ []byte) (interface{}, error) {
	return api.Response{Code: 0}, nil
}

func handleFUserCreate(r []byte) (interface{}, error) {
	// parse args
	var args api.ArgsFUserCreate
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return api.Response{Code: api.EArgsInval}, err
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		return api.Response{Code: api.EUnknown}, err
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
		return api.Response{Code: api.EExists}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	// success
	return api.Response{Code: 0}, nil
}

func handleFLogIn(r []byte) (interface{}, error) {
	// parse args
	var args api.ArgsFLogIn
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return api.Response{Code: api.EArgsInval}, err
	}

	session, err := database.OpenSession(db, args.Login, args.Password)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return api.Response{Code: api.ENoEntry}, nil
	case database.ErrPassWrong:
		return api.Response{Code: api.EPassWrong}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	return api.RespFLogIn{Code: 0, Token: string(session.Token)}, nil
}

func handleFLogOut(r []byte) (interface{}, error) {
	// parse args
	var args api.ArgsFLogOut
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return api.Response{Code: api.EArgsInval}, err
	}

	err = database.CloseSession(db, []byte(args.Token))
	if err != nil {
		return api.Response{Code: api.EUnknown}, err
	}

	return api.Response{Code: 0}, nil
}

func handleFUserInfo(r []byte) (interface{}, error) {
	// parse args
	var args api.ArgsFUserInfo
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return api.Response{Code: api.EArgsInval}, err
	}

	_, err = database.VerifySession(db, []byte(args.Token))
	switch err {
	case nil:
		break
	case database.ErrNotLoggedIn:
		return api.Response{Code: api.ENotLoggedIn}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	userinfo, err := database.FindUserInfo(db, args.Login)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return api.Response{Code: api.ENoEntry}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	return api.RespFUserInfo{
		Code:       0,
		Login:      userinfo.Login,
		FirstName:  userinfo.FirstName,
		LastName:   userinfo.LastName,
		Patronymic: userinfo.Patronymic,
	}, nil
}

func handleFUserEdit(r []byte) (interface{}, error) {
	// parse args
	var args api.ArgsFUserEdit
	err := msgpack.Unmarshal(r, &args)
	if err != nil || args.Token == "" {
		return api.Response{Code: api.EArgsInval}, err
	}

	session, err := database.VerifySession(db, []byte(args.Token))
	switch err {
	case nil:
		break
	case database.ErrNotLoggedIn:
		return api.Response{Code: api.ENotLoggedIn}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	uid := session.User

	if args.Login != nil {
		err = database.UserChangeLogin(db, uid, *args.Login)
		switch err {
		case nil:
			break
		case database.ErrNoUser:
			return api.Response{Code: api.ENoEntry}, nil
		case database.ErrUserExists:
			return api.Response{Code: api.EExists}, nil
		default:
			return api.Response{Code: api.EUnknown}, err
		}
	}

	if args.Password != nil {
		passHash, err := bcrypt.GenerateFromPassword([]byte(*args.Password), bcrypt.DefaultCost)
		if err != nil {
			return api.Response{Code: api.EUnknown}, err
		}
		err = database.UserChangePasswordHash(db, uid, passHash)
		switch err {
		case nil:
			break
		case database.ErrNoUser:
			return api.Response{Code: api.ENoEntry}, nil
		default:
			return api.Response{Code: api.EUnknown}, err
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
			return api.Response{Code: api.ENoEntry}, nil
		default:
			return api.Response{Code: api.EUnknown}, err
		}
	}

	return api.Response{Code: 0}, nil
}

func handleFUserSetManagesGroups(r []byte) (interface{}, error) {
	// parse args
	var args api.ArgsFUserSetManagesGroups
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return api.Response{Code: api.EArgsInval}, err
	}

	// check that user can manage groups
	resp, err := verifyManagesGroups(args.Token)
	if resp != nil {
		return resp, err
	}

	// find target user
	userinfo, err := database.FindUserInfo(db, args.Login)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return api.Response{Code: api.ENoEntry}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	err = database.UserSetManagesGroups(db, userinfo.Id, args.Value)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return api.Response{Code: api.ENoEntry}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	return api.Response{Code: 0}, nil
}

// verifyManagesGroups: check that user can manage groups
func verifyManagesGroups(token string) (interface{}, error) {
	session, err := database.VerifySession(db, []byte(token))
	switch err {
	case nil:
		break
	case database.ErrNotLoggedIn:
		return api.Response{Code: api.ENotLoggedIn}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	userinfo, err := database.GetUserInfo(db, session.User)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return api.Response{Code: api.ENoEntry}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	if !userinfo.ManagesGroups {
		return api.Response{Code: api.EAccessDenied}, nil
	}

	return nil, nil
}

func handleFGroupCreate(r []byte) (interface{}, error) {
	// parse args
	var args api.ArgsFGroupCreateRemove
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return api.Response{Code: api.EArgsInval}, err
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
		return api.Response{Code: api.EExists}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	return api.Response{Code: 0}, nil
}

func handleFGroupRemove(r []byte) (interface{}, error) {
	// parse args
	var args api.ArgsFGroupCreateRemove
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return api.Response{Code: api.EArgsInval}, err
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
		return api.Response{Code: api.ENoEntry}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	err = database.RemoveGroup(db, group.Id)
	switch err {
	case nil:
		break
	case database.ErrNoGroup:
		return api.Response{Code: api.ENoEntry}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	return api.Response{Code: 0}, nil
}

func handleFGroupAddRemoveUser(r []byte) (interface{}, error) {
	// parse args
	var args api.ArgsFGroupAddRemoveUser
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return api.Response{Code: api.EArgsInval}, err
	}

	// check that user can manage groups
	resp, err := verifyManagesGroups(args.Token)
	if resp != nil {
		return resp, err
	}

	group, err := database.FindGroup(db, args.Group)
	switch err {
	case nil:
		break
	case database.ErrNoGroup:
		return api.Response{Code: api.ENoEntry}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	user, err := database.FindUserInfo(db, args.Login)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return api.Response{Code: api.ENoEntry}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	if args.Action {
		err = database.GroupAddUser(db, user.Id, group.Id)
	} else {
		err = database.GroupRemoveUser(db, user.Id, group.Id)
	}

	switch err {
	case nil:
		break
	case database.ErrAlreadyInGroup:
		return api.Response{Code: api.EExists}, nil
	case database.ErrNoGroup:
		return api.Response{Code: api.ENoEntry}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	return api.Response{Code: 0}, nil
}

func handleFStructCreate(r []byte) (interface{}, error) {
	// parse args
	var args api.ArgsFStructCreate
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return api.Response{Code: api.EArgsInval}, err
	}

	_, err = database.VerifySession(db, []byte(args.Token))
	switch err {
	case nil:
		break
	case database.ErrNotLoggedIn:
		return api.Response{Code: api.ENotLoggedIn}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	structInfo := database.StructInfo{
		Id:          0,
		Name:        args.Name,
		Description: args.Description,
		District:    args.District,
		Region:      args.Region,
		Address:     args.Address,
		Type:        args.Type,
		State:       args.State,
		Area:        args.Area,
		Owner:       args.Owner,
		Gid:         args.Gid,
		Permissions: args.Permissions,
	}
	err = structInfo.AddStruct(db)
	switch err {
	case nil:
		break
	case database.ErrStructExists:
		return api.Response{Code: api.EExists}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	// success
	return api.RespFStructCreate{Code: 0, Id: structInfo.Id}, nil
}

func handleFStructInfo(r []byte) (interface{}, error) {
	// parse args
	var args api.ArgsFStructInfo
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return api.Response{Code: api.EArgsInval}, err
	}

	_, err = database.VerifySession(db, []byte(args.Token))
	switch err {
	case nil:
		break
	case database.ErrNotLoggedIn:
		return api.Response{Code: api.ENotLoggedIn}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	structInfo, err := database.GetStructInfo(db, args.Id)
	switch err {
	case nil:
		break
	case database.ErrNoStruct:
		return api.Response{Code: api.ENoEntry}, nil
	default:
		return api.Response{Code: api.EUnknown}, err
	}

	return api.RespFStructInfo{
		Name:        structInfo.Name,
		Description: structInfo.Description,
		District:    structInfo.District,
		Region:      structInfo.Region,
		Address:     structInfo.Address,
		Type:        structInfo.Type,
		State:       structInfo.State,
		Area:        structInfo.Area,
		Owner:       structInfo.Owner,
		Actual_user: structInfo.Actual_user,
		Gid:         structInfo.Gid,
		Permissions: structInfo.Permissions,
	}, nil
}

var apiFHandlers map[string]api.RequestHandler

func main() {
	var err error

	/* setup handlers */

	apiFHandlers = make(map[string]api.RequestHandler)

	apiFHandlers["ping"] = handleFPing

	apiFHandlers["user_create"] = handleFUserCreate
	apiFHandlers["user_log_in"] = handleFLogIn
	apiFHandlers["user_log_out"] = handleFLogOut
	apiFHandlers["user_get_info"] = handleFUserInfo
	apiFHandlers["user_edit"] = handleFUserEdit
	apiFHandlers["user_set_manages_groups"] = handleFUserSetManagesGroups

	apiFHandlers["group_create"] = handleFGroupCreate
	apiFHandlers["group_remove"] = handleFGroupRemove
	apiFHandlers["group_add_remove_user"] = handleFGroupAddRemoveUser

	apiFHandlers["object_create"] = handleFStructCreate
	apiFHandlers["object_get_info"] = handleFStructInfo
	/* =(setup handlers)= */

	db, err = database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/", apiHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
