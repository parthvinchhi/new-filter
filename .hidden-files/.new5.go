package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Define your struct according to the JSON structure
type Info struct {
	Source        string `json:"source"`
	CartId        string `json:"cart_id"`
	StoreId       string `json:"store_id"`
	XcallId       string `json:"xcall_id"`
	BoardRole     string `json:"board_role"`
	ErrorCode     int    `json:"error_code"`
	RaspiIssueId  string `json:"raspi_issue_id"`
	RaspiIssueMsg string `json:"raspi_issue_msg"`
}

type Data struct {
	Id        string    `json:"id"`
	Source    string    `json:"source"`
	Info      Info      `json:"info"`
	StoreId   string    `json:"store_id"`
	CartId    string    `json:"cart_id"`
	CreatedAt time.Time `json:"created_at"`
	ErrorCode string    `json:"error_code"`
}

func main() {
	inputDir := "json-files"           // Directory containing the JSON files
	outputDir := "filtered-json-files" // Directory to save the filtered JSON files

	// Ensure the output directory exists
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	// List all JSON files in the input directory
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		fmt.Println("Error reading input directory:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			inputFilePath := filepath.Join(inputDir, file.Name())
			outputFilePath := filepath.Join(outputDir, file.Name())

			data, err := loadJSON(inputFilePath)
			if err != nil {
				fmt.Printf("Error loading JSON file %s: %v\n", inputFilePath, err)
				continue
			}

			filteredData := filterData(data)

			err = saveJSON(filteredData, outputFilePath)
			if err != nil {
				fmt.Printf("Error saving filtered data to %s: %v\n", outputFilePath, err)
				continue
			}

			fmt.Printf("Filtered data written to %s\n", outputFilePath)
		}
	}
}

// loadJSON reads a JSON file and unmarshals it into a slice of Data
func loadJSON(filename string) ([]Data, error) {
	// Read the file
	byteValue, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Parse the JSON data
	var data []Data
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// filterData filters the data based on the RaspiIssueMsg condition
func filterData(data []Data) []Data {
	var filterData []Data
	for _, item := range data {
		if strings.Contains(item.Info.RaspiIssueMsg, "Error! Captured 0 frames") {
			filterData = append(filterData, item)
		}
	}
	return filterData
}

// saveJSON marshals the data and writes it to a new JSON file
func saveJSON(data []Data, filename string) error {
	// Marshal the data to JSON
	filteredDataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Write the filtered data to a new JSON file
	err = ioutil.WriteFile(filename, filteredDataJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}
