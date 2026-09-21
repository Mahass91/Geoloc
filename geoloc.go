package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func loadCsvData(fpath string) [][]string {
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatal("Could not open File:", err)
	}
	defer f.Close()
	dataReader := csv.NewReader(f)

	dataSet, err := dataReader.ReadAll()
	if err != nil {
		log.Fatal("Could not get Data: ", err)
	}
	return dataSet
}

func createJSON(JSONData []byte, fpath string) {
	newJSON, err := os.Create(fpath)
	if err != nil {
		log.Fatal("Not able to create JSON:", err)
	}
	fmt.Println("JSON Created!")
	defer newJSON.Close()
	newJSON.Write(JSONData)
	newJSON.Close()
}
func trimData(list []string) {
	for each := range list {
		list[each] = strings.TrimSpace(list[each])
	}
}
func main() {
	// Change Paths to load csv and save json to
	csvFilePath := "csv_file/adress.csv"
	jsonfilePath := "csv_file/adress.json"

	adress := loadCsvData(csvFilePath)

	type adressJSON struct {
		Streetname string `json:"streetname"`
		Number     string `json:"number"`
		Offer      string `json:"offer"`
	}
	var traderadress adressJSON
	var allAdresses []adressJSON

	for _, d := range adress {
		trimData(d)
		traderadress.Streetname = d[0]
		traderadress.Number = d[1]
		traderadress.Offer = d[2]
		allAdresses = append(allAdresses, traderadress)
	}

	JSONData, err := json.Marshal(allAdresses)
	if err != nil {
		log.Fatal("Error by creation of JSON", err)
	}
	createJSON(JSONData, jsonfilePath)
	fmt.Println("Work completed!")

}
