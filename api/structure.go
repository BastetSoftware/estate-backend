package api

import (
	"BastetSoftware/backend/database"

	"github.com/vmihailenco/msgpack/v5"
)

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
		Actual_user: args.Actual_user,
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
		Code:       0,
		Structures: structsInfo,
	}, nil
}

func HandleFDeleteStruct(r []byte) (interface{}, error) {
	var args ArgsFDeleteStruct
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
	err = database.DeleteStruct(Db, args.Id)
	switch err {
	case nil:
		break
	case database.ErrNoStruct:
		return Response{Code: ENoEntry}, nil
	default:
		return Response{Code: EUnknown}, err
	}
	return Response{Code: 0}, nil
}

func HandleFStructEdit(r []byte) (interface{}, error) {
	var args ArgsFStructEdit
	err := msgpack.Unmarshal(r, &args)
	if err != nil || args.Token == "" {
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

	uid := args.Id

	if args.Name != nil {
		err = database.StructChangeName(Db, uid, *args.Name)
		switch err {
		case nil:
			break
		default:
			return Response{Code: EUnknown}, err
		}
	}

	if args.Description != nil {
		err = database.StructChangeDescription(Db, uid, *args.Description)
		switch err {
		case nil:
			break
		default:
			return Response{Code: EUnknown}, err
		}
	}

	if args.District != nil {
		err = database.StructChangeDistrict(Db, uid, *args.District)
		switch err {
		case nil:
			break
		default:
			return Response{Code: EUnknown}, err
		}
	}

	if args.Region != nil {
		err = database.StructChangeRegion(Db, uid, *args.Region)
		switch err {
		case nil:
			break
		default:
			return Response{Code: EUnknown}, err
		}
	}

	if args.Address != nil {
		err = database.StructChangeAddress(Db, uid, *args.Address)
		switch err {
		case nil:
			break
		default:
			return Response{Code: EUnknown}, err
		}
	}

	if args.Type != nil {
		err = database.StructChangeType(Db, uid, *args.Type)
		switch err {
		case nil:
			break
		default:
			return Response{Code: EUnknown}, err
		}
	}

	if args.State != nil {
		err = database.StructChangeState(Db, uid, *args.State)
		switch err {
		case nil:
			break
		default:
			return Response{Code: EUnknown}, err
		}
	}

	if args.Area != nil {
		err = database.StructChangeArea(Db, uid, *args.Area)
		switch err {
		case nil:
			break
		default:
			return Response{Code: EUnknown}, err
		}
	}

	if args.Owner != nil {
		err = database.StructChangeOwner(Db, uid, *args.Owner)
		switch err {
		case nil:
			break
		default:
			return Response{Code: EUnknown}, err
		}
	}

	if args.Actual_user != nil {
		err = database.StructChangeActualUser(Db, uid, *args.Actual_user)
		switch err {
		case nil:
			break
		default:
			return Response{Code: EUnknown}, err
		}
	}

	if args.Permissions != nil {
		err = database.StructChangePermissions(Db, uid, *args.Permissions)
		switch err {
		case nil:
			break
		default:
			return Response{Code: EUnknown}, err
		}
	}

	return Response{Code: 0}, nil
}
