package api

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

type Request struct {
	Func uint8
	Args []byte
}

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
	Code   uint8
	Groups []string
	Count  int
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

type RequestHandler func(r []byte) (interface{}, error)
