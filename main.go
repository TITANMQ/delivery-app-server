package main

import (
	"backend/app"
	"backend/controllers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) //middleware Jwt auth

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println(port)

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")

	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	router.HandleFunc("/api/user/new/profile", controllers.CreateProfile).Methods("POST")

	router.HandleFunc("/api/delivery/new", controllers.CreateDelivery).Methods("POST")

	router.HandleFunc("/api/user/{id}", controllers.GetAccountProfile).Methods("GET")

	router.HandleFunc("/api/user/deliveries/{id}", controllers.GetDeliveriesFor).Methods("GET")

	router.HandleFunc("/api/user/deliveries/search/{radius}", controllers.SearchDeliveries).Methods("POST")

	router.HandleFunc("/api/delivery/accept", controllers.AcceptDelivery).Methods("POST")

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
