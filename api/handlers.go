package api

import (
	"encoding/json"
	"github.com/Russiancold/testApp/service"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type createRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeInternal(w, err)
		return
	}
	req := createRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		writeBadReq(w, err)
		return
	}
	if !isEmailValid(req.Email) {
		writeBadReq(w, service.InvalidEmail)
		return
	}
	if len(req.Username) < 1 || len(req.Username) > 30 {
		writeBadReq(w, service.InvalidName)
		return
	}
	if err := service.GetService().CreateAccount(req.Username, req.Email); err != nil {
		if err == service.AlreadyExist {
			w.WriteHeader(http.StatusConflict)
			if _, e := w.Write(marshalError(err)); e != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
		writeInternal(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getToken(w http.ResponseWriter, r * http.Request) {
	name := mux.Vars(r)["name"]
	token, err := service.GetService().GetToken(name)
	if err != nil {
		if err == service.NoUser {
			writeBadReq(w, err)
			return
		}
		if err == service.NoContent {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		writeInternal(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"token":"` + token + `"}"`)); err != nil {
		writeInternal(w, err)
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if err := service.GetService().DeleteAccount(name); err != nil {
		writeInternal(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func updateToken(w http.ResponseWriter, r * http.Request) {
	name := mux.Vars(r)["name"]
	token := r.FormValue("token")
	if err := service.GetService().UpdateToken(name, token); err != nil {
		if err == service.NoUser {
			writeBadReq(w, err)
			return
		}
		writeInternal(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(email string) bool {
	return emailRegexp.MatchString(email)
}

func marshalError(err error) []byte {
	return []byte(`{"error":"` + err.Error() +`"}`)
}

func writeBadReq(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusBadRequest)
	log.Println("bad request", e)
	if _, err := w.Write(marshalError(e)); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func writeInternal(w http.ResponseWriter, e error) {
	log.Println("internal err", e)
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := w.Write(marshalError(e)); err != nil {
		log.Println(err)
	}
}