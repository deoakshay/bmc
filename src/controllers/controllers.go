package controllers

import (
	"bmc/src/models"
)

var store *models.PassengersDatabase

func InitController(passengers models.PassengersDatabase) {
	store = &passengers
}

func GetAllPassengers() (models.Passengers, error) {
	passengers, err := store.GetAllPassengers()
	if err != nil {
		return nil, err
	}
	return passengers, nil
}
func GetPassengerByID(id int) (models.Passenger, error) {
	passenger, err := store.GetPassengerByID(id)
	if err != nil {
		return models.Passenger{}, err
	}
	return passenger, nil
}
func GetPassengerByAttributes(id int, attributes []string) (models.Passenger, error) {
	passenger, err := store.GetPassengerByAttributes(id, attributes)
	if err != nil {
		return models.Passenger{}, err
	}
	return passenger, nil
}
func GetHistogram() error {
	err := store.GetHistogram()
	if err != nil {
		return err
	}
	return nil
}
