package models

import "encoding/json"

type PassengersDB interface {
	GetPassengerByID(id int) (Passenger, error)
	GetPassengerByAttributes(id int, attributes []string) (Passenger, error)
	GetAllPassengers() (Passengers, error)
	GetHistogram() error
	BuildTableData()
}

type Passenger struct {
	PassengerID int     `csv:"PassengerId" json:"passenger_id,omitempty"`
	Survived    int     `csv:"Survived" json:"survived,omitempty"`
	Pclass      int     `csv:"Pclass" json:"pclass,omitempty"`
	Name        string  `csv:"Name" json:"name,omitempty"`
	Sex         string  `csv:"Sex" json:"sex,omitempty"`
	Age         float64 `csv:"Age" json:"age,omitempty"`
	SibSp       int     `csv:"SibSp" json:"sibsp,omitempty"`
	Parch       int     `csv:"Parch" json:"parch,omitempty"`
	Ticket      string  `csv:"Ticket" json:"ticket,omitempty"`
	Fare        float64 `csv:"Fare" json:"fare,omitempty"`
	Cabin       string  `csv:"Cabin" json:"cabin,omitempty"`
	Embarked    string  `csv:"Embarked" json:"embarked,omitempty"`
}

func (p Passenger) JsonMarshal() (string, error) {
	jsonData, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

type Passengers []Passenger

func (p Passengers) JsonMarshal() (string, error) {
	jsonData, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

type PassengersDatabase struct {
	passengers Passengers
	tableData  map[int]Passenger
	histogram  map[int]int
}
