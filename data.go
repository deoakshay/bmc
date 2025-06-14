package main

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
)

func readFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
	defer file.Close()

	var p PassengersDatabase
	if err := gocsv.UnmarshalFile(file, &p.passengers); err != nil {
		return err
	}

	p.BuildTableData()
	allPassengers, err := p.GetAllPassengers()
	if err != nil {
		return fmt.Errorf("error getting all passengers: %v", err)

	}
	jsonData, err := allPassengers.JsonMarshal()
	if err != nil {
		return fmt.Errorf("error marshalling passengers data to JSON: %v", err)
	}
	fmt.Println("All passengers in JSON format:", string(jsonData))
	attributes := []string{"Name", "Cabin", "Survived", "Pclass", "Fare"}
	val, err := p.GetPassengerByAttributes(15, attributes)
	if err != nil {
		return fmt.Errorf("error getting passenger by attributes: %v", err)
	}
	jsonData, err = val.JsonMarshal()
	if err != nil {
		return fmt.Errorf("error marshalling passenger data to JSON: %v", err)
	}
	fmt.Println("Passenger data in JSON format:", string(jsonData))

	err = p.GetHistogram()
	if err != nil {
		return fmt.Errorf("error getting histogram: %v", err)
	}

	return nil
}
