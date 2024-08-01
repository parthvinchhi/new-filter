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

	// Parse the JSON data
	var data []Data

	// Filter the data
	var filteredData []Data
	for _, item := range data {
		// if item.RaspiIssueMsg == "captured 0 frames" {
		// 	filteredData = append(filteredData, item)
		// }

		if item.RaspiIssueMsg != nil {

		}

		if strings.Contains(item.RaspiIssueMsg, "Error! Captured 0 frames") {
			filteredData = append(filteredData, item)
		}

		if filteredData != nil {
			fmt.Println("Not nil")
		}
	}

	// Convert the filtered data to JSON
	filteredDataJSON, err := json.MarshalIndent(filteredData, "", "  ")
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
