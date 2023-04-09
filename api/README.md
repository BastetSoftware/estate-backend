# API

## Request format

A POST request to api/function_name with binary msgpack-encoded arguments in the body.

## Response format

| field | type  | description                                    |
|-------|-------|------------------------------------------------|
| Code  | uint8 | 0 if no errors occurred, error code otherwise  |
| ...   | ...   | function-specific return values (if no errors) |

## Error codes

|     name     | code |
|:------------:|:----:|
|   EExists    |  1   |
|   ENoEntry   |  2   |
|  EPassWrong  |  3   |
| ENotLoggedIn |  4   |
|  EArgsInval  | 253  |
|    ENoFun    | 254  |
|   EUnknown   | 255  |

## Functions

### Service

#### ping

##### Request args

None

##### Response data

No specific response data

##### Possible errors

Always successful

### User data manipulation

#### user_create

##### Request args

| argument   | type   | description           |
|------------|--------|-----------------------|
| Login      | string | unique login          |
| Password   | string | password              |
| FirstName  | string | first name            |
| LastName   | string | last name             |
| Patronymic | string | patronymic (optional) |

##### Response data

No specific response data

##### Possible errors

| error      | description                         |
|------------|-------------------------------------|
| EArgsInval | invalid request arguments           |
| EExists    | user with this login already exists |
| EUnknown   | unknown error                       |

#### user_log_in

##### Request args

| argument   | type   | description |
|------------|--------|-------------|
| Login      | string | login       |
| Password   | string | password    |

##### Response data

| field | type   | description   |
|-------|--------|---------------|
| Token | string | session token |

##### Possible errors

| error      | description               |
|------------|---------------------------|
| EArgsInval | invalid request arguments |
| ENoEntry   | user does not exist       |
| EPassWrong | wrong password            |
| EUnknown   | unknown error             |

#### user_log_out

##### Request args

| argument | type   | description   |
|----------|--------|---------------|
| Token    | string | session token |

##### Response data

No specific response data

##### Possible errors

| error      | description               |
|------------|---------------------------|
| EArgsInval | invalid request arguments |
| EUnknown   | unknown error             |

#### user_get_info

##### Request args

| argument | type   | description              |
|----------|--------|--------------------------|
| Token    | string | session token            |
| Login    | string | login of the target user |

##### Response data

| argument   | type   | description |
|------------|--------|-------------|
| Login      | string | login       |
| FirstName  | string | first name  |
| LastName   | string | last name   |
| Patronymic | string | patronymic  |

##### Possible errors

| error        | description                                                 |
|--------------|-------------------------------------------------------------|
| EArgsInval   | invalid request arguments                                   |
| ENotLoggedIn | request sender is not logged in or session token is invalid |
| ENoEntry     | target user does not exist                                  |
| EUnknown     | unknown error                                               |

#### user_edit

##### Request args

| argument   | type    | description                        |
|------------|---------|------------------------------------|
| Token      | string  | session token                      |
| Login      | *string | new login (null if unchanged)      |
| Password   | *string | new password (null if unchanged)   |
| FirstName  | *string | new first name (null if unchanged) |
| LastName   | *string | new last name (null if unchanged)  |
| Patronymic | *string | new patronymic (null if unchanged) |

##### Response data

No specific response data

##### Possible errors

| error        | description                                                 |
|--------------|-------------------------------------------------------------|
| EArgsInval   | invalid request arguments                                   |
| ENotLoggedIn | request sender is not logged in or session token is invalid |
| ENoEntry     | user does not exist                                         |
| EExists      | user with this new login already exists                     |
| EUnknown     | unknown error                                               |

#### user_set_manages_groups

##### Request args

| argument | type   | description                                 |
|----------|--------|---------------------------------------------|
| Token    | string | session token                               |
| Login    | string | login of the target user                    |
| Value    | bool   | permission (true to grant, false to revoke) |

##### Response data

No specific response data

##### Possible errors

| error         | description                                                 |
|---------------|-------------------------------------------------------------|
| EArgsInval    | invalid request arguments                                   |
| ENotLoggedIn  | request sender is not logged in or session token is invalid |
| ENoEntry      | user does not exist                                         |
| EAccessDenied | user has no rights to grant this permission                 |
| EUnknown      | unknown error                                               |

#### user_list_groups

##### Request args

| argument | type   | description       |
|----------|--------|-------------------|
| Token    | string | session token     |
| Login    | string | target user login |

##### Response data

| argument | type    | description            |
|----------|---------|------------------------|
| Name     | string  | group name             |
| Gids     | int64[] | groups user belongs to |
| Count    | int     | number of groups       |

##### Possible errors

| error        | description               |
|--------------|---------------------------|
| EArgsInval   | invalid request arguments |
| ENoEntry     | user does not exist       |
| EUnknown     | unknown error             |

### Groups manipulation

#### group_create

##### Request args

| argument | type   | description   |
|----------|--------|---------------|
| Token    | string | session token |
| Name     | string | group name    |

##### Response data

No specific response data

##### Possible errors

| error         | description                                                 |
|---------------|-------------------------------------------------------------|
| EArgsInval    | invalid request arguments                                   |
| ENotLoggedIn  | request sender is not logged in or session token is invalid |
| ENoEntry      | request sending user does not exist                         |
| EAccessDenied | user has no rights to manage groups                         |
| EExists       | group already exists                                        |
| EUnknown      | unknown error                                               |

#### group_remove

##### Request args

| argument | type   | description   |
|----------|--------|---------------|
| Token    | string | session token |
| Name     | string | group name    |

##### Response data

No specific response data

##### Possible errors

| error         | description                                                 |
|---------------|-------------------------------------------------------------|
| EArgsInval    | invalid request arguments                                   |
| ENotLoggedIn  | request sender is not logged in or session token is invalid |
| ENoEntry      | group does not exist                                        |
| ENoEntry (2)  | request sending user does not exist                         |
| EAccessDenied | user has no rights to manage groups                         |
| EUnknown      | unknown error                                               |

#### group_add_remove_user

Add or remove a user from a group

##### Request args

| argument | type   | description                |
|----------|--------|----------------------------|
| Token    | string | session token              |
| Group    | string | group name                 |
| Login    | string | user login                 |
| Action   | bool   | true - add, false - remove |

##### Response data

No specific response data

##### Possible errors

| error         | description                                                 |
|---------------|-------------------------------------------------------------|
| EArgsInval    | invalid request arguments                                   |
| ENotLoggedIn  | request sender is not logged in or session token is invalid |
| ENoEntry      | group or user does not exist                                |
| EAccessDenied | user has no rights to manage groups                         |
| EExists       | user is already in group                                    |
| EUnknown      | unknown error                                               |

#### group_get_info

##### Request args

| argument | type   | description   |
|----------|--------|---------------|
| Token    | string | session token |
| Gid      | int64  | GID           |

##### Response data

| argument | type    | description                  |
|----------|---------|------------------------------|
| Name     | string  | group name                   |
| Uids     | int64[] | users in the group           |
| Count    | int     | number of users in the group |

##### Possible errors

| error        | description               |
|--------------|---------------------------|
| EArgsInval   | invalid request arguments |
| ENoEntry     | group does not exist      |
| EUnknown     | unknown error             |
