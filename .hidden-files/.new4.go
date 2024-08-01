package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	// Directory containing JSON files
	jsonDir := "filtered-json-files" // Replace with your directory name

	// Create a map to hold unique xcall_id for each store_id-cart_id-boardRole
	uniqueXcallIDs := make(map[string]map[string]struct{})

	// Read and process each JSON file in the directory
	err := filepath.Walk(jsonDir, func(path string, info os.FileInfo, err error) error {
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
		filename := fmt.Sprintf("txt-files/%s.txt", key)
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
