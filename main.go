package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"reflect"

	"github.com/dmfrank/egaas/cache"
	"github.com/dmfrank/egaas/model"
	"github.com/gorilla/mux"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	auth cache.Auth
	work cache.Work
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", mainPage)
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/login/pass", changePass).Methods("POST")
	r.HandleFunc("/login/job", doWork).Methods("POST")

	http.Handle("/", r)
}

func init() {
	auth = cache.Auth{
		Values: make(map[string]string, 0)}
	work = cache.Work{
		Values: make(map[string]int32, 0)}
	work.Values["admin"] = int32(10000000)
	model.GormInit()
}

func verifyUser(login, pass string) bool {
	if !auth.IsExist(login, pass) {
		u := &model.User{}
		if e := u.Get(login, pass); e != nil {
			return false
		}
		auth.Push(login, pass)
	}
	return true
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte(`
		<!DOCTYPE html>
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if !verifyUser(
			r.FormValue("login"),
			r.FormValue("pass")) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func changePass(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		login := r.FormValue("login")
		if !verifyUser(login, r.FormValue("pass")) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u := &model.User{
			Login: login,
			Pass:  r.FormValue("new_pass"),
		}

		if err := u.Update(); err != nil {
			auth.Push(u.Login, u.Pass)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// DTO data transfer obj
type DTO struct {
	BigNumber int64
	Text      string
}

func doWork(w http.ResponseWriter, r *http.Request) {
	var value DTO

	login := r.FormValue("login")
	if !verifyUser(login, r.FormValue("pass")) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if work.Values[login] < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.Unmarshal([]byte(r.FormValue("value")), &value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	v := reflect.ValueOf(&value).Elem()

	for i := 0; i < v.NumField(); i++ {
		w.Write(reverse(v.Field(i)))
	}
}

func reverse(val reflect.Value) []byte {
	switch val.Kind().String() {
	case "int64":
		result := make([]byte, 8)
		binary.LittleEndian.PutUint64(result, uint64(math.MaxInt64-val.Interface().(int64)))
		return []byte(fmt.Sprintf("%v", binary.BigEndian.Uint64(result)))
	case "string":
		result := []rune(val.Interface().(string))
		for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
			result[i], result[j] = result[j], result[i]
		}
		return []byte(string(result))
	}
	return nil
}
