package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Change Paths to load csv and save json to
	csvFilePath := "csv_file/adress.csv"
	jsonfilePath := "csv_file/adress.json"
	// Change CSV to Struct
	adress := loadCsvData(csvFilePath)

	var allAdresses []addressJson

	for _, each := range adress {
		var traderadress addressJson
		trimData(each)
		traderadress.Streetname = each[0]
		traderadress.Number = each[1]
		traderadress.Offer = each[2]
		urlApi := createUrl(traderadress)
		results := fetchNominatimData(urlApi)
		addCoordsToStruct(results, &traderadress)
		time.Sleep(1 * time.Second)
		allAdresses = append(allAdresses, traderadress)
	}

	JSONData, err := json.Marshal(allAdresses)
	if err != nil {
		log.Fatal("Error by creation of JSON", err)
	}
	createJson(JSONData, jsonfilePath)
	fmt.Println("Work completed!")

}

type addressJson struct {
	Streetname string  `json:"streetname"`
	Number     string  `json:"number"`
	Offer      string  `json:"offer"`
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"lon"`
}
type nominatimCoords struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func createUrl(address addressJson) string {
	params := url.Values{
		"street":     {address.Streetname + " " + address.Number},
		"postalcode": {"38124"},
		"city":       {"Brunswick"},
		"country":    {"Germany"},
		"format":     {"jsonv2"},
		"limit":      {"1"},
	}
	url := "https://nominatim.openstreetmap.org/search?" + params.Encode()

	return url

}
func fetchNominatimData(url string) []nominatimCoords {
	var results []nominatimCoords

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Request Error:", err)
	}
	req.Header.Set("User-Agent",
		"FleamarketOrganisation/1.0 mhass@murena.io")
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("API Error:", err)
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		log.Fatal("Could not Decode:", err)
	}
	return results
}
func addCoordsToStruct(res []nominatimCoords, address *addressJson) {
	if len(res) == 0 {
		fmt.Println("No Data! Streetname and Number:", address.Streetname, address.Number)
		return
	}
	lat, err := strconv.ParseFloat(res[0].Lat, 64)
	if err != nil {
		log.Fatal("Could not convert String:", err)
	}

	lon, err := strconv.ParseFloat(res[0].Lon, 64)
	if err != nil {
		log.Fatal("Could not convert String:", err)
	}
	address.Latitude = lat
	address.Longitude = lon
}

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

func createJson(JSONData []byte, fpath string) {
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
