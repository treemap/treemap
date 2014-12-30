package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func GenerateZipcodes() {
	log.Println("Gathering Zipcodes")
	zipcodes := AllZipcodes()

	log.Println("Got Zipcodes")
	for i := range zipcodes {
		log.Println("Writing", zipcodes[i].Number)
		zipcodes[i].GetInfo()

		os.MkdirAll("static/data/zipcodes", os.ModeDir|os.ModePerm)

		b, err := json.Marshal(zipcodes[i])
		if err != nil {
			log.Println("error:", err)
		}

		// Write the file
		err = ioutil.WriteFile("static/data/zipcodes/"+zipcodes[i].Number+".json", b, 0644)

	}
}

func GenerateZipcodeTable(tableName string, distance uint) {
	log.Println("Gathering Zipcodes")
	zipcodes := AllZipcodes()

	dirName := fmt.Sprintf("static/data/%s/%d", tableName, distance)
	fmt.Println(dirName)

	// Mkdir if it doesn't exist
	os.MkdirAll(dirName, os.ModeDir|os.ModePerm)

	for i := range zipcodes {
		// Get the Tree information
		data := zipcodes[i].TableData(tableName, distance)

		b, err := json.Marshal(data)
		if err != nil {
			log.Println("error:", err)
		}

		// Write the file to disk.
		err = ioutil.WriteFile(dirName+"/"+zipcodes[i].Number+".json", b, 0644)
	}
}
