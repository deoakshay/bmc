package handlers

import (
	"bmc/src/controllers"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func GetPassengersHandler(w http.ResponseWriter, r *http.Request) {
	passengers, err := controllers.GetAllPassengers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(passengers)
}

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

func GetHistogram(w http.ResponseWriter, r *http.Request) {

	err := controllers.GetHistogram()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imgPath := "histogram.png"

	imgBytes, err := ioutil.ReadFile(imgPath)
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
