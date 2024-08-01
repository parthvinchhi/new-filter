package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

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
	jsonFile, err := os.Open("json-files/10-days.json")
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

	uniqueXcallIDs := make(map[string]map[string]struct{})

	// Populate the map with unique xcall_id values
	for _, item := range data {
		key := fmt.Sprintf("%s-%s-%s", item.Info.StoreId, item.Info.CartId, item.Info.BoardRole)
		if _, exists := uniqueXcallIDs[key]; !exists {
			uniqueXcallIDs[key] = make(map[string]struct{})
		}
		uniqueXcallIDs[key][item.Info.XcallId] = struct{}{}
	}

	// Write the unique xcall_id values to corresponding files
	for key, xcallIDs := range uniqueXcallIDs {
		filename := fmt.Sprintf("txt-files/%s.txt", key)
		file, err := os.Create(filename)
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
