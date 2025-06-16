package main

import (
	"bmc/src/controllers"
	"bmc/src/models"
	"bmc/src/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	var store models.PassengersDB
	if os.Getenv("USE_SQL") == "true" {
		db, err := models.NewSQLitePassengerDatabase("data/titanic.db", "data/titanic.csv")
		if err != nil {
			log.Fatal("Error loading SQLite database:", err)
		}
		store = db
	} else {
		db, err := models.NewPassengerDatabaseFromCSV("data/titanic.csv")
		if err != nil {
			log.Fatalf("Error loading CSV: %v", err)
		}
		store = db
	}

	controllers.InitController(store)
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
