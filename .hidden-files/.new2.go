package main

import (
	"encoding/json"
	"fmt"
	"os"
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

	// CreatedAt string `json:"created_at"`
}

func main() {
	// Open the JSON file
	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	// Get the file size
	fileInfo, err := jsonFile.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fileInfo.Size()

	// Create a byte slice of the appropriate size
	byteValue := make([]byte, fileSize)

	// Read the file into the byte slice
	_, err = jsonFile.Read(byteValue)
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

	var filterData []Data
	// Print the data to the terminal
	for _, item := range data {
		if strings.Contains(item.Info.RaspiIssueMsg, "Error! Captured 0 frames") {
			filterData = append(filterData, item)
		}
		// fmt.Printf("RaspiIssueMsg: %s\n", item.Info.RaspiIssueMsg)
		// Print other fields if needed
	}

	filteredDataJSON, err := json.MarshalIndent(filterData, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Write the filtered data to a new JSON file
	newFile, err := os.Create("filtered_data.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer newFile.Close()

	_, err = newFile.Write(filteredDataJSON)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Filtered data written to filtered_data.json")
}
