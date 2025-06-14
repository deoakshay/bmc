package main

import "fmt"

func main() {
	err := readFile("./titanic.csv")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
