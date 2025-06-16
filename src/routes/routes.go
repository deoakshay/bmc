package routes

import (
	"bmc/src/handlers"

	"github.com/gorilla/mux" // One of the popular routing frameworks
)

func RegisterRoutes(router *mux.Router) {

	router.HandleFunc("/api/getPassengers", handlers.GetPassengersHandler).Methods("GET")
	router.HandleFunc("/api/getPassengers/{id}", handlers.GetPassengerByID).Methods("GET")
	router.HandleFunc("/api/getPassengersByAttributes/{id}", handlers.GetPassengerByAttributes).Methods("GET")
	router.HandleFunc("/api/getPassengersHistogram", handlers.GetHistogram).Methods("GET")

}
