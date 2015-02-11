package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func GenerateTrees() {
	log.Println("Gathering Trees")
	trees := AllTrees()

	log.Println("Got Trees")

	dirName := "static/data/trees"
	os.MkdirAll(dirName, os.ModeDir|os.ModePerm)

	// Write the tree index file
	b, err := json.Marshal(trees)
	if err != nil {
		log.Println("error:", err)
	}
	err = ioutil.WriteFile(dirName+"/index.json", b, 0644)

	for i := range trees {
		log.Println("Writing", trees[i].LatinName)
		trees[i].GetGeodata()
		trees[i].GetArea()

		b, err := json.Marshal(trees[i])
		if err != nil {
			log.Println("error:", err)
		}

		// Write the file
		err = ioutil.WriteFile(fmt.Sprintf("%s/%d.json", dirName, trees[i].Id), b, 0644)

	}
}

func GenerateParks() {
	log.Println("Gathering Parks")
	parks := AllNationalParks()

	log.Println("Got Parks")

	dirName := "static/data/parks"
	os.MkdirAll(dirName, os.ModeDir|os.ModePerm)

	// Write the tree index file
	b, err := json.Marshal(parks)
	if err != nil {
		log.Println("error:", err)
	}
	err = ioutil.WriteFile(dirName+"/index.json", b, 0644)

	for i := range parks {
		log.Println("Writing", parks[i].UnitName)

		b, err := json.Marshal(parks[i])
		if err != nil {
			log.Println("error:", err)
		}

		// Write the file
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", dirName, parks[i].UnitCode), b, 0644)

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
