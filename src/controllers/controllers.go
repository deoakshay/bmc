package controllers

import (
	"encoding/json"
	"fmt"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func (p *PassengersDatabase) GetAllPassengers() (Passengers, error) {
	if len(p.passengers) == 0 {
		return nil, fmt.Errorf("no passengers found")
	}
	return p.passengers, nil
}

func (p *PassengersDatabase) GetPassengerByID(id int) (Passenger, error) {
	passengers := p.tableData
	if passenger, ok := passengers[id]; ok {
		return passenger, nil
	}
	return Passenger{}, fmt.Errorf("passenger with ID %d not found", id)
}

func (p *PassengersDatabase) BuildTableData() {
	data := make(map[int]Passenger)
	for _, passenger := range p.passengers {
		data[passenger.PassengerID] = passenger
	}
	p.tableData = data
}

func (p *PassengersDatabase) GetPassengerByAttributes(id int, attributes []string) (Passenger, error) {
	passenger, err := p.GetPassengerByID(id)
	if err != nil {
		return Passenger{}, err
	}

	attributesMap := getOnePassengerMap(passenger)
	filteredPassengerWithAttributes := Passenger{}
	for attr, _ := range attributes {
		val, exists := attributesMap[attributes[attr]]
		if exists {
			switch attributes[attr] {
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

func getOnePassengerMap(p Passenger) map[string]interface{} {
	return map[string]interface{}{
		"PassengerId": p.PassengerID,
		"Survived":    p.Survived,
		"Pclass":      p.Pclass,
		"Name":        p.Name,
		"Sex":         p.Sex,
		"Age":         p.Age,
		"SibSp":       p.SibSp,
		"Parch":       p.Parch,
		"Ticket":      p.Ticket,
		"Fare":        p.Fare,
		"Cabin":       p.Cabin,
		"Embarked":    p.Embarked,
	}
}

func (p *PassengersDatabase) GetHistogram() error {
	p.histogram = buildHistogram(p.passengers, 10)
	if len(p.histogram) == 0 {
		return fmt.Errorf("histogram not built")
	}
	return plotHistogram(p.histogram)
}

func buildHistogram(passengers Passengers, level int) map[int]int {
	fares := []float64{}
	for _, passenger := range passengers {
		if passenger.Fare > 0 {
			fares = append(fares, passenger.Fare)
		}
	}
	sort.Float64s(fares)
	length := len(fares)
	histogram := make(map[int]int)
	for p := 0; p <= 100; p += level {
		if length == 0 {
			histogram[p] = 0
			continue
		}
		i := int(float64(p) / 100.0 * float64(length-1))
		fare := fares[i]
		count := 0
		for _, fareValue := range fares {
			if fareValue >= fare && fareValue < fare+float64(level) {
				count++
			}
		}
		histogram[p] = count
	}
	return histogram
}

func plotHistogram(histogram map[int]int) error {
	fmt.Println("Plotting")
	p := plot.New()
	p.Title.Text = "Passenger Fare Histogram"
	p.X.Label.Text = "Fare Range"
	p.Y.Label.Text = "Count fares <= percentile"

	coordinates := make(plotter.XYs, 0, len(histogram))
	for percentile, count := range histogram {
		coordinates = append(coordinates, plotter.XY{X: float64(percentile), Y: float64(count)})
	}
	sort.Slice(coordinates, func(i, j int) bool {
		return coordinates[i].X < coordinates[j].X
	})
	line, err := plotter.NewLine(coordinates)
	if err != nil {
		return fmt.Errorf("Error plotting graph:", err)
	}
	p.Add(line)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "histogram.png"); err != nil {
		return fmt.Errorf("Error saving histogram plot:", err)
	}
	fmt.Println("Histogram plot saved as histogram.png")
	return nil
}
