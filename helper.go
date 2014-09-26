package main

import (
	"encoding/json"
	"log"
)

// error handling
func PanicIf(err error) {
	if err != nil {
		logError(err.Error())
		panic(err)
	}
}

// logging
func logError(message string) {
	log.Println("[Error] " + message)
}

func logInfo(message string) {
	log.Println("[Info] " + message)
}

func logDebug(message string) {
	log.Println("[Debug] " + message)
}

//json
func toJSON(target interface{}) (string, error) {
	byteArry, err := json.Marshal(target)
	if err != nil {
		return "", err
	}
	return string(byteArry[:]), nil
}

func fromJSON(target interface{}, jsonString string) error {
	byteArray := []byte(jsonString)
	err := json.Unmarshal(byteArray, target)
	if err != nil {
		return err
	}
	return nil
}
