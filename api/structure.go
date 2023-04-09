package api

import "BastetSoftware/backend/database"

func HandleFStructCreate(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFStructCreate
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
	err = structInfo.AddStruct(Db)
	switch err {
	case nil:
		break
	case database.ErrStructExists:
		return Response{Code: EExists}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	// success
	return RespFStructCreate{Code: 0, Id: structInfo.Id}, nil
}

func HandleFStructInfo(r []byte) (interface{}, error) {
	// parse args
	var args ArgsFStructInfo
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

	structInfo, err := database.GetStructInfo(Db, args.Id)
	switch err {
	case nil:
		break
	case database.ErrNoStruct:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	return RespFStructInfo{
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
func HandleFStructFind(r []byte) (interface{}, error) {
	var args database.ArgsFStructFind
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

	structsInfo, err := database.FindStructures(Db, args)
	switch err {
	case nil:
		break
	case database.ErrNoStruct:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}

	return RespFStructFind{
		Structures: structsInfo,
	}, nil
}
