package models

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	_ "github.com/mattn/go-sqlite3"
)

type SQLitePassengerDatabase struct {
	db *sql.DB
}

func NewSQLitePassengerDatabase(path, fileName string) (*SQLitePassengerDatabase, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS passengers (
		PassengerId INTEGER PRIMARY KEY,
		Survived INTEGER,
		Pclass INTEGER,
		Name TEXT,
		Sex TEXT,
		Age REAL,
		SibSp INTEGER,
		Parch INTEGER,
		Ticket TEXT,
		Fare REAL,
		Cabin TEXT,
		Embarked TEXT
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(`
		INSERT OR IGNORE INTO passengers (
			PassengerId, Survived, Pclass, Name, Sex, Age,
			SibSp, Parch, Ticket, Fare, Cabin, Embarked
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	csvFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer csvFile.Close()
	var csvPassengers Passengers
	if err := gocsv.UnmarshalFile(csvFile, &csvPassengers); err != nil {
		return nil, err
	}

	for _, p := range csvPassengers {
		_, err := stmt.Exec(
			p.PassengerID, p.Survived, p.Pclass, p.Name, p.Sex, p.Age,
			p.SibSp, p.Parch, p.Ticket, p.Fare, p.Cabin, p.Embarked,
		)
		if err != nil {
			return nil, err
		}
	}

	return &SQLitePassengerDatabase{db: db}, nil
}

func (s *SQLitePassengerDatabase) Close() error {
	return s.db.Close()
}

func (s *SQLitePassengerDatabase) GetAllPassengers() (Passengers, error) {
	rows, err := s.db.Query("SELECT * FROM passengers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passengers Passengers
	for rows.Next() {
		var p Passenger
		if err := rows.Scan(&p.PassengerID, &p.Survived, &p.Pclass, &p.Name, &p.Sex, &p.Age, &p.SibSp, &p.Parch, &p.Ticket, &p.Fare, &p.Cabin, &p.Embarked); err != nil {
			return nil, err
		}
		passengers = append(passengers, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return passengers, nil
}

func (s *SQLitePassengerDatabase) GetPassengerByID(id int) (Passenger, error) {
	row := s.db.QueryRow("SELECT * FROM passengers WHERE PassengerId = ?", id)
	var p Passenger
	if err := row.Scan(&p.PassengerID, &p.Survived, &p.Pclass, &p.Name, &p.Sex, &p.Age, &p.SibSp, &p.Parch, &p.Ticket, &p.Fare, &p.Cabin, &p.Embarked); err != nil {
		if err == sql.ErrNoRows {
			return Passenger{}, nil
		}
		return Passenger{}, err
	}
	return p, nil
}

func (s *SQLitePassengerDatabase) GetPassengerByAttributes(id int, attributes []string) (Passenger, error) {
	passenger, err := s.GetPassengerByID(id)
	if err != nil {
		return Passenger{}, err
	}

	attributesMap := getOnePassengerMap(passenger)
	filteredPassengerWithAttributes := Passenger{}
	for _, attr := range attributes {
		val, exists := attributesMap[attr]
		if exists {
			switch attr {
			case "PassengerId":
				filteredPassengerWithAttributes.PassengerID = val.(int)
			case "Survived":
				filteredPassengerWithAttributes.Survived = val.(int)
			case "Pclass":
				filteredPassengerWithAttributes.Pclass = val.(int)
			case "Name":
				filteredPassengerWithAttributes.Name = val.(string)
			case "Sex":
				filteredPassengerWithAttributes.Sex = val.(string)
			case "Age":
				filteredPassengerWithAttributes.Age = val.(float64)
			case "SibSp":
				filteredPassengerWithAttributes.SibSp = val.(int)

			case "Parch":
				filteredPassengerWithAttributes.Parch = val.(int)
			case "Ticket":
				filteredPassengerWithAttributes.Ticket = val.(string)
			case "Fare":
				filteredPassengerWithAttributes.Fare = val.(float64)
			case "Cabin":
				filteredPassengerWithAttributes.Cabin = val.(string)
			case "Embarked":
				filteredPassengerWithAttributes.Embarked = val.(string)
			}

		}
	}
	return filteredPassengerWithAttributes, nil

}

func (s *SQLitePassengerDatabase) GetHistogram() error {
	p, err := s.GetAllPassengers()
	if err != nil {
		return err
	}
	h := buildHistogram(p, 10)
	if len(h) == 0 {
		return fmt.Errorf("histogram not built")
	}
	return plotHistogram(h)
}
