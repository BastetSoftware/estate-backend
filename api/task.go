package api

import "BastetSoftware/backend/database"

func HandleFTaskCreate(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFTaskCreate
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

	task := database.Task{
		Id:          0,
		Name:        args.Name,
		Description: args.Description,
		Deadline:    args.Deadline,
		Status:      args.Status,
		Object:      args.Object,
		Maintainer:  args.Maintainer,
		Gid:         args.Gid,
		Permissions: args.Permissions,
	}
	task.Id, err = database.CreateTask(Db, &task)
	switch err {
	case nil:
		break
	case database.ErrTaskExists:
		return Response{Code: EExists}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	// success
	return RespFTaskCreate{Code: 0, Id: task.Id}, nil
}

func HandleFTaskRemove(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFTaskRemove
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

	err = database.RemoveTask(Db, args.Id)
	switch err {
	case nil:
		break
	case database.ErrNoTask:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	// success
	return Response{Code: 0}, nil
}
