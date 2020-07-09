package controllers

import (
	"backend/models"
	u "backend/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//CreateAccount is a function used to handle creating account requests
var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //retrieves data from request body and saves in the account struct
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
	}

	resp := account.Create()
	u.Respond(w, resp)
}

//CreateProfile is a function used to handle creating profile requests
var CreateProfile = func(w http.ResponseWriter, r *http.Request) {

	profile := &models.Profile{}
	err := json.NewDecoder(r.Body).Decode(profile) //retrieves data from request body and saves in the account struct
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
	}

	resp := profile.Create()
	u.Respond(w, resp)
}

//Authenticate is a function used to handle authentication
var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //retrieves data from request body and saves in the account struct
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}

//GetAccountProfile is a function used to handle get account profile requests
var GetAccountProfile = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}

	data := models.GetAccountProfile(uint(id))
	resp := u.Message(true, "Success")
	resp["data"] = data
	u.Respond(w, resp)

}
