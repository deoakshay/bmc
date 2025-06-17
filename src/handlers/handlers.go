package handlers

import (
	"bmc/src/controllers"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// @Summary Get all passengers
// @Description Retrieves all passengers from the database
// @Tags Passengers
// @Accept json
// @Produce json
// @Success 200 {array} models.Passenger
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/passengers [get]
func GetPassengersHandler(w http.ResponseWriter, r *http.Request) {
	passengers, err := controllers.GetAllPassengers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(passengers)
}

// @Summary Get passenger by ID
// @Description Retrieves a passenger by their ID
// @Tags Passengers
// @Param id path int true "Passenger ID"
// @Accept json
// @Produce json
// @Success 200 {object} models.Passenger
// @Failure 400 {string} string "Invalid ID"
// @Router /api/passengers/{id} [get]
func GetPassengerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	passenger, err := controllers.GetPassengerByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(passenger)
}

// @Summary Get passenger by attributes
// @Description Retrieves a passenger by their ID and specified attributes
// @Tags Passengers
// @Param id path int true "Passenger ID"
// @Param attributes query string true "Comma-separated list of attributes e.g attributes=Name,Age"
// @Accept json
// @Produce json
// @Success 200 {object} models.Passenger
// @Failure 400 {string} string "Invalid ID or attributes"
// @Router /api/passenger/{id}/attributes  [get]
func GetPassengerByAttributes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	attributes := r.URL.Query()["attributes"]
	fmt.Println("Attributes received:", attributes)
	split := strings.Split(attributes[0], ",")
	if len(attributes) == 0 {
		http.Error(w, "No attributes provided", http.StatusBadRequest)
		return
	}

	passenger, err := controllers.GetPassengerByAttributes(id, split)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(passenger)
}

// @Summary Get histogram of passengers
// @Description Generates and retrieves a histogram of passengers
// @Tags Passengers
// @Accept json
// @Produce html
// @Success 200 {string} string "Histogram image in HTML format"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/histogram  [get]
func GetHistogram(w http.ResponseWriter, r *http.Request) {

	err := controllers.GetHistogram()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imgPath := "histogram.png"

	imgBytes, err := os.ReadFile(imgPath)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	encoded := base64.StdEncoding.EncodeToString(imgBytes)
	html := fmt.Sprintf(`
        <html>
            <body>
                <h2>Histogram</h2>
                <img src="data:image/png;base64,%s" alt="Histogram"/>
            </body>
        </html>
    `, encoded)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
