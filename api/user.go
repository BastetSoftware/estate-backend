package api

import (
	"BastetSoftware/backend/database"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/bcrypt"
)

func HandleFUserCreate(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFUserCreate
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return Response{Code: EArgsInval}, err
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		return Response{Code: EUnknown}, err
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
	err = userInfo.Register(Db)
	switch err {
	case nil:
		break
	case database.ErrUserExists:
		return Response{Code: EExists}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	// success
	return Response{Code: 0}, nil
}

func HandleFLogIn(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFLogIn
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return Response{Code: EArgsInval}, err
	}

	session, err := database.OpenSession(Db, args.Login, args.Password)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return Response{Code: ENoEntry}, nil
	case database.ErrPassWrong:
		return Response{Code: EPassWrong}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	return RespFLogIn{Code: 0, Token: string(session.Token)}, nil
}

func HandleFLogOut(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFLogOut
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return Response{Code: EArgsInval}, err
	}

	err = database.CloseSession(Db, []byte(args.Token))
	if err != nil {
		return Response{Code: EUnknown}, err
	}

	return Response{Code: 0}, nil
}

func HandleFUserInfo(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFUserInfo
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return Response{Code: EArgsInval}, err
	}

	_, err = database.VerifySession(Db, []byte(args.Token))
	switch err {
	case nil:
		break
	case database.ErrNotLoggedIn:
		return Response{Code: ENotLoggedIn}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	userinfo, err := database.FindUserInfo(Db, args.Login)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	return RespFUserInfo{
		Code:       0,
		Login:      userinfo.Login,
		FirstName:  userinfo.FirstName,
		LastName:   userinfo.LastName,
		Patronymic: userinfo.Patronymic,
	}, nil
}

func HandleFUserEdit(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFUserEdit
	err := msgpack.Unmarshal(r, &args)
	if err != nil || args.Token == "" {
		return Response{Code: EArgsInval}, err
	}

	session, err := database.VerifySession(Db, []byte(args.Token))
	switch err {
	case nil:
		break
	case database.ErrNotLoggedIn:
		return Response{Code: ENotLoggedIn}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	uid := session.User

	if args.Login != nil {
		err = database.UserChangeLogin(Db, uid, *args.Login)
		switch err {
		case nil:
			break
		case database.ErrNoUser:
			return Response{Code: ENoEntry}, nil
		case database.ErrUserExists:
			return Response{Code: EExists}, nil
		default:
			return Response{Code: EUnknown}, err
		}
	}

	if args.Password != nil {
		passHash, err := bcrypt.GenerateFromPassword([]byte(*args.Password), bcrypt.DefaultCost)
		if err != nil {
			return Response{Code: EUnknown}, err
		}
		err = database.UserChangePasswordHash(Db, uid, passHash)
		switch err {
		case nil:
			break
		case database.ErrNoUser:
			return Response{Code: ENoEntry}, nil
		default:
			return Response{Code: EUnknown}, err
		}
	}

	var names = [3]*string{args.FirstName, args.LastName, args.Patronymic}
	for i, n := range names {
		if n == nil {
			continue
		}

		err = database.UserChangeName(Db, uid, i, *n)
		switch err {
		case nil:
			break
		case database.ErrNoUser:
			return Response{Code: ENoEntry}, nil
		default:
			return Response{Code: EUnknown}, err
		}
	}

	return Response{Code: 0}, nil
}

func HandleFUserSetManagesGroups(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFUserSetManagesGroups
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return Response{Code: EArgsInval}, err
	}

	// check that user can manage groups
	resp, err := verifyManagesGroups(args.Token)
	if resp != nil {
		return resp, err
	}

	// find target user
	userinfo, err := database.FindUserInfo(Db, args.Login)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	err = database.UserSetManagesGroups(Db, userinfo.Id, args.Value)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	return Response{Code: 0}, nil
}

func HandleFUserListGroups(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFUserListGroups
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return Response{Code: EArgsInval}, err
	}

	// find target user
	userinfo, err := database.FindUserInfo(Db, args.Login)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	// get user groups
	gids, err := database.ListGroupsOrUsers(Db, database.UserListGroups, userinfo.Id)
	if err != nil {
		return Response{Code: EUnknown}, err
	}

	return RespFUserListGroups{
		Code:  0,
		Gids:  gids,
		Count: len(gids),
	}, nil
}
