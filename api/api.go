package api

import (
	"bytes"
	"database/sql"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	EExists uint8 = iota + 1 // record exists
	ENoEntry
	EPassWrong
	ENotLoggedIn
	EAccessDenied

	EArgsInval uint8 = 253 // invalid arguments
	ENoFun     uint8 = 254 // function does not exist
	EUnknown   uint8 = 255 // unknown error
)

/*
 * API functions' arguments and responses
 */

type Request struct {
	Func uint8
	Args []byte
}

// Response: basic response
type Response struct {
	Code uint8
}

/* FUserCreate */

type ArgsFUserCreate struct {
	Login      string
	Password   string
	FirstName  string
	LastName   string
	Patronymic string
}

/* FLogIn */

type ArgsFLogIn struct {
	Login    string
	Password string
}

type RespFLogIn struct {
	Code  uint8
	Token string
}

/* FLogOut */

type ArgsFLogOut struct {
	Token string
}

/* FUserInfo */

type ArgsFUserInfo struct {
	Token string
	Login string // login
}

type RespFUserInfo struct {
	Code       uint8
	Login      string
	FirstName  string
	LastName   string
	Patronymic string
}

/* FUserEdit */

type ArgsFUserEdit struct {
	Token      string
	Login      *string
	Password   *string
	FirstName  *string
	LastName   *string
	Patronymic *string
}

/* FUserSetManagesGroups */

type ArgsFUserSetManagesGroups struct {
	Token string
	Login string
	Value bool
}

/* FUserListGroups */

type ArgsFUserListGroups struct {
	Token string
	Login string
}

type RespFUserListGroups struct {
	Code  uint8
	Gids  []int64
	Count int
}

/* FGroupCreate */
/* FGroupRemove */

type ArgsFGroupCreateRemove struct {
	Token string
	Name  string
}

/* FGroupAddRemoveUser */

type ArgsFGroupAddRemoveUser struct {
	Token  string
	Group  string
	Login  string
	Action bool // true - add, false - remove
}

/* FGroupGetInfo */

type ArgsFGroupGetInfo struct {
	Token string
	Gid   int64
}

type RespFGroupGetInfo struct {
	Code  uint8
	Name  string
	Uids  []int64
	Count int
}

/* FStructCreate */

type ArgsFStructCreate struct {
	Token       string
	Name        string
	Description string
	District    string
	Region      string
	Address     string
	Type        string
	State       string
	Area        int32
	Owner       string
	Actual_user string
	Gid         int64
	Permissions int8
}

type RespFStructCreate struct {
	Code uint8
	Id   int64
}

/* FStructInfo */

type ArgsFStructInfo struct {
	Token string
	Id    int64
}

type RespFStructInfo struct {
	Name        string
	Description string
	District    string
	Region      string
	Address     string
	Type        string
	State       string
	Area        int32
	Owner       string
	Actual_user string
	Gid         int64
	Permissions int8
}

/* FTaskCreate */

type ArgsFTaskCreate struct {
	Token       string
	Name        string
	Description string
	Deadline    int64
	Status      string
	Object      int64
	Maintainer  int64
	Gid         int64
	Permissions uint8
}

type RespFTaskCreate struct {
	Code uint8
	Id   int64
}

/* FTaskRemove */

type ArgsFTaskRemove struct {
	Token string
	Id    int64
}

/* FTaskGetInfo */

type ArgsFTaskGetInfo struct {
	Token string
	Id    int64
}

type RespFTaskGetInfo struct {
	Code        uint8
	Name        string
	Description string
	Deadline    int64
	Status      string
	Object      int64
	Maintainer  int64
	Gid         int64
	Permissions uint8
}

/*
 * Common
 */

type RequestHandler func(r []byte) (interface{}, error)

func CustomUnmarshal(data []byte, v interface{}) error {
	dec := msgpack.GetDecoder()

	dec.Reset(bytes.NewReader(data))
	dec.DisallowUnknownFields(true) // <- this is customized
	err := dec.Decode(v)

	msgpack.PutDecoder(dec)

	return err
}

var Db *sql.DB // Db reference
