# API

## Request format

| field | type    | description                 |
|-------|---------|-----------------------------|
| Func  | uint8   | function number             |
| Args  | uint8[] | function-specific arguments |

## Response format

| field | type    | description                                   |
|-------|---------|-----------------------------------------------|
| Code  | uint8   | 0 if no errors occurred, error code otherwise |
| Data  | uint8[] | function-specific return values               |

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

#### Ping

Func = 0

##### Request args

None

##### Response data

No specific response data

##### Possible errors

Always successful

### User data manipulation

#### Create user

Func = 1

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

#### Log in

Func = 2

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

#### Log out

Func = 3

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

#### Get user info

Func = 4

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

#### Edit user info

Func = 5

##### Request args

| argument  | type    | description                        |
|-----------|---------|------------------------------------|
| Token     | string  | session token                      |
| Login     | *string | new login (null if unchanged)      |
| Password  | *string | new password (null if unchanged)   |
| FirstName | *string | new first name (null if unchanged) |
| LastName  | *string | new last name (null if unchanged)  |

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
