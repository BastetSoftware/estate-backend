package main

import (
	"BastetSoftware/backend/api"
	"BastetSoftware/backend/database"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/vmihailenco/msgpack/v5"
)

var origin string

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
	w.Header().Set("Access-Control-Allow-Origin", origin)

	var buf []byte
	var n int
	handler := apiFHandlers[r.URL.Path[len("/api/"):]]
	if handler == nil {
		// unknown API function, no arguments
		handler = api.UnknownFPlug
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

var apiFHandlers map[string]api.RequestHandler

func main() {
	var err error

	if o := os.Getenv("RESPONSE_ORIGIN"); o != "" {
		origin = o
	} else {
		origin = "*"
	}

	/* setup handlers */

	apiFHandlers = make(map[string]api.RequestHandler)

	apiFHandlers["ping"] = api.HandleFPing

	apiFHandlers["user_create"] = api.HandleFUserCreate
	apiFHandlers["user_log_in"] = api.HandleFLogIn
	apiFHandlers["user_log_out"] = api.HandleFLogOut
	apiFHandlers["user_get_info"] = api.HandleFUserInfo
	apiFHandlers["user_edit"] = api.HandleFUserEdit
	apiFHandlers["user_set_manages_groups"] = api.HandleFUserSetManagesGroups
	apiFHandlers["user_list_groups"] = api.HandleFUserListGroups

	apiFHandlers["group_create"] = api.HandleFGroupCreate
	apiFHandlers["group_remove"] = api.HandleFGroupRemove
	apiFHandlers["group_add_remove_user"] = api.HandleFGroupAddRemoveUser
	apiFHandlers["group_get_info"] = api.HandleFGroupGetInfo

	apiFHandlers["object_create"] = api.HandleFStructCreate
	apiFHandlers["object_get_info"] = api.HandleFStructInfo
	apiFHandlers["find_object"] = api.HandleFStructFind
	apiFHandlers["object_delete"] = api.HandleFDeleteStruct

	apiFHandlers["task_create"] = api.HandleFTaskCreate
	apiFHandlers["task_remove"] = api.HandleFTaskRemove
	apiFHandlers["task_get_info"] = api.HandleFTaskGetInfo
	
	/* =(setup handlers)= */

	api.Db, err = database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/", apiHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
