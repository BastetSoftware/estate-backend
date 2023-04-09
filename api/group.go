package api

import "BastetSoftware/backend/database"

// verifyManagesGroups: check that user can manage groups
func verifyManagesGroups(token string) (interface{}, error) {
	session, err := database.VerifySession(Db, []byte(token))
	switch err {
	case nil:
		break
	case database.ErrNotLoggedIn:
		return Response{Code: ENotLoggedIn}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	userinfo, err := database.GetUserInfo(Db, session.User)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	if !userinfo.ManagesGroups {
		return Response{Code: EAccessDenied}, nil
	}

	return nil, nil
}

func HandleFGroupCreate(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFGroupCreateRemove
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return Response{Code: EArgsInval}, err
	}

	// check that user can manage groups
	resp, err := verifyManagesGroups(args.Token)
	if resp != nil {
		return resp, err
	}

	_, err = database.CreateGroup(Db, args.Name)
	switch err {
	case nil:
		break
	case database.ErrGroupExists:
		return Response{Code: EExists}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	return Response{Code: 0}, nil
}

func HandleFGroupRemove(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFGroupCreateRemove
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return Response{Code: EArgsInval}, err
	}

	// check that user can manage groups
	resp, err := verifyManagesGroups(args.Token)
	if resp != nil {
		return resp, err
	}

	group, err := database.FindGroup(Db, args.Name)
	switch err {
	case nil:
		break
	case database.ErrNoGroup:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	err = database.RemoveGroup(Db, group.Id)
	switch err {
	case nil:
		break
	case database.ErrNoGroup:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	return Response{Code: 0}, nil
}

func HandleFGroupAddRemoveUser(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFGroupAddRemoveUser
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return Response{Code: EArgsInval}, err
	}

	// check that user can manage groups
	resp, err := verifyManagesGroups(args.Token)
	if resp != nil {
		return resp, err
	}

	group, err := database.FindGroup(Db, args.Group)
	switch err {
	case nil:
		break
	case database.ErrNoGroup:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	user, err := database.FindUserInfo(Db, args.Login)
	switch err {
	case nil:
		break
	case database.ErrNoUser:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	if args.Action {
		err = database.GroupAddUser(Db, user.Id, group.Id)
	} else {
		err = database.GroupRemoveUser(Db, user.Id, group.Id)
	}

	switch err {
	case nil:
		break
	case database.ErrAlreadyInGroup:
		return Response{Code: EExists}, nil
	case database.ErrNoGroup:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	return Response{Code: 0}, nil
}

func HandleFGroupGetInfo(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFGroupGetInfo
	err := CustomUnmarshal(r, &args)
	if err != nil {
		return Response{Code: EArgsInval}, err
	}

	// get group info
	group, err := database.GetGroup(Db, args.Gid)
	switch err {
	case nil:
		break
	case database.ErrNoGroup:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	// get group's users
	uids, err := database.ListGroupsOrUsers(Db, database.GroupListUsers, args.Gid)
	if err != nil {
		return Response{Code: EUnknown}, err
	}

	return RespFGroupGetInfo{
		Code:  0,
		Name:  group.Name,
		Uids:  uids,
		Count: len(uids),
	}, nil
}
