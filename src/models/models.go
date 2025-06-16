package models

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/gocarina/gocsv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type PassengersDB interface {
	GetPassengerByID(id int) (Passenger, error)
	GetPassengerByAttributes(id int, attributes []string) (Passenger, error)
	GetAllPassengers() (Passengers, error)
	GetHistogram() error
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
	Passengers Passengers
	TableData  map[int]Passenger
	Histogram  map[int]int
}

func NewPassengerDatabaseFromCSV(fileName string) (*PassengersDatabase, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var p Passengers
	if err := gocsv.UnmarshalFile(file, &p); err != nil {
		return nil, err
	}
	data := make(map[int]Passenger)
	for _, passenger := range p {
		data[passenger.PassengerID] = passenger
	}
	return &PassengersDatabase{
		Passengers: p,
		TableData:  data,
		Histogram:  make(map[int]int),
	}, nil
}
func (p *PassengersDatabase) GetAllPassengers() (Passengers, error) {
	if len(p.Passengers) == 0 {
		return nil, fmt.Errorf("no passengers found")
	}
	return p.Passengers, nil
}

func (p *PassengersDatabase) GetPassengerByID(id int) (Passenger, error) {
	passengers := p.TableData
	fmt.Println("Getting Passenger by ID:", id)
	if passenger, ok := passengers[id]; ok {
		return passenger, nil
	}
	return Passenger{}, fmt.Errorf("passenger with ID %d not found", id)
}

func (p *PassengersDatabase) GetPassengerByAttributes(id int, attributes []string) (Passenger, error) {
	passenger, err := p.GetPassengerByID(id)
	if err != nil {
		return Passenger{}, err
	}
	fmt.Println("Getting Passenger by Attributes:", attributes, " for Passenger ID:", id, " Passenger:", passenger)

	attributesMap := getOnePassengerMap(passenger)
	filteredPassengerWithAttributes := Passenger{}
	for attr, _ := range attributes {
		val, exists := attributesMap[attributes[attr]]
		fmt.Println("val:", val, " exists:", exists, " attr:", attributes[attr])
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
	fmt.Println("Filtered Passenger with Attributes:", filteredPassengerWithAttributes, " Attributes:", attributes)
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
	p.Histogram = buildHistogram(p.Passengers, 10)
	if len(p.Histogram) == 0 {
		return fmt.Errorf("histogram not built")
	}
	return plotHistogram(p.Histogram)
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
