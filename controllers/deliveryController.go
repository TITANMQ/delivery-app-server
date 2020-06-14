package controllers

import (
	"backend/models"
	u "backend/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var CreateDelivery = func(w http.ResponseWriter, r *http.Request) {
	delivery := &models.Delivery{}
	err := json.NewDecoder(r.Body).Decode(delivery)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
	}

	resp := delivery.Create()
	u.Respond(w, resp)
}

var GetDeliveriesFor = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}

	data := models.GetDeliveries(uint(id))
	resp := u.Message(true, "Success")
	resp["data"] = data
	u.Respond(w, resp)
}

var SearchDeliveries = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	radius, err := strconv.ParseFloat(params["radius"], 32)

	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}

	searchJourney := &models.Journey{}
	err = json.NewDecoder(r.Body).Decode(searchJourney)

	fmt.Println(searchJourney.JourneyStart)

	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request body"))
		fmt.Println(err)
		return
	}

	data := models.GetDeliveriesByRadius(searchJourney, float32(radius))
	if data == nil {
		u.Respond(w, u.Message(false, "No results found"))
		return
	}
	resp := u.Message(true, "Success")
	resp["data"] = data
	u.Respond(w, resp)
}

var AcceptDelivery = func(w http.ResponseWriter, r *http.Request) {
	delivery := &models.Accepted{}
	err := json.NewDecoder(r.Body).Decode(delivery)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
	}
	resp := delivery.AcceptDelivery()
	u.Respond(w, resp)
}
