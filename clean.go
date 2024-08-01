package main

import (
	"encoding/json"
	"fmt"
	"io"
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

var inputDir = "json-files"           // Directory containing the JSON files
var outputDir = "filtered-json-files" // Directory to save the filtered JSON files
var txtDir = "txt-files/"

func main() {
	FilterData()
	CreateTxtFile()
}

func FilterData() {
	// List all JSON files in the input directory
	files, err := os.ReadDir(inputDir)
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
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Create a byte slice of the appropriate size
	byteValue := make([]byte, fileInfo.Size())

	// Read the file into the byte slice
	_, err = file.Read(byteValue)
	if err != nil && err != io.EOF {
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
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(filteredDataJSON)
	if err != nil {
		return err
	}

	return nil
}

func CreateTxtFile() {
	uniqueXcallIDs := make(map[string]map[string]struct{})

	// Read and process each JSON file in the directory
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			processJSONFile(path, uniqueXcallIDs)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Write the unique xcall_id values to corresponding files
	for key, xcallIDs := range uniqueXcallIDs {
		filename := fmt.Sprintf(txtDir+"%s.txt", key)
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer file.Close()

		for xcallID := range xcallIDs {
			_, err := file.WriteString(fmt.Sprintf("%s\n", xcallID))
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}

	fmt.Println("Unique xcall_id values written to respective files.")
}

func processJSONFile(path string, uniqueXcallIDs map[string]map[string]struct{}) {
	// Open the JSON file
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	// Read the file content
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse the JSON data
	var data []Data
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Populate the map with unique xcall_id values
	for _, item := range data {
		key := fmt.Sprintf("%s-%s-%s", item.Info.StoreId, item.Info.CartId, item.Info.BoardRole)
		if _, exists := uniqueXcallIDs[key]; !exists {
			uniqueXcallIDs[key] = make(map[string]struct{})
		}
		uniqueXcallIDs[key][item.Info.XcallId] = struct{}{}
	}
}
