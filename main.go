package main

import (
	"encoding/binary"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/TestTask/model"
	"github.com/gorilla/mux"
)

var (
	Auth map[string]string
	Work map[string]int32
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", mainPage)
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/login/pass", changePass).Methods("POST")

	http.Handle("/", r)
}

func init() {
	Work := make(map[string]int32, 0)
	Work["admin"] = 1000000
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>
		<html>
		<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<meta name="theme-color" content="#375EAB">
		
			<title>main page</title>
		</head>
		<body>
			Page body and some more content
		</body>
		</html>`))
}

func login(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	pass := r.FormValue("pass")

	if Auth[login] == pass {
		w.WriteHeader(http.StatusOK)
	}

	user := &model.User{}
	err := user.Get(login, pass)
	if err == nil {
		Auth[login] = pass
		Work[login] = user.WorkNumber
	}

	w.WriteHeader(http.StatusBadRequest)
}

func changePass(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	pass := r.FormValue("pass")

	newPass := r.FormValue("newPass")

	if Auth[login] != pass {
		w.WriteHeader(http.StatusBadRequest)
	}

	user := &model.User{}
	user.Pass = newPass
	err := user.Save()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type DTO struct {
	bigNumber int64
	text      string
}

func doWork(w http.ResponseWriter, r *http.Request) {
	var value DTO
	login := r.FormValue("login")
	if Work[login] <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.Unmarshal([]byte(r.FormValue("value")), &value)

	v := reflect.ValueOf(value)
	for i := 0; i < v.NumField(); i++ {
		w.Write(reverse(v.Elem().Field(i)))
	}
}

func reverse(val reflect.Value) []byte {
	switch val.Kind().String() {
	case "int64":
		fallthrough
	case "int32":
		result := make([]byte, 4)
		binary.LittleEndian.PutUint32(result, uint32(2147483647-val.Interface().(int32)))
		return result
	case "string":
		var result string
		for i := len(val.Interface().(string)); i > 0; i++ {
			result += string(val.Interface().(string)[i])
		}
		return []byte(result)
	}
	return nil
}
