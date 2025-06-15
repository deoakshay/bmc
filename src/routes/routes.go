package routes

import (
 "net/http"

 "github.com/gorilla/mux" // One of the popular routing frameworks
)

func SetupRoutes() *mux.Router {
 router := mux.NewRouter()

 router.HandleFunc("/api/getPassengers", GetUsers).Methods("GET")
 router.HandleFunc("/api/getPassengers/{id}", GetUserByID).Methods("GET")
 router.HandleFunc("/api/getPassengersByAttributes/{id}",GetPassengerByAttributes).Methods("GET")
 router.HandleFunc("/api/getPassengersHistogram/{id}", GetHistogram).Methods("GET")


 return router
}
