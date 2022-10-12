package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintPrettyFormatStruct(obj interface{}) {
	empJSON, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%s\n", string(empJSON))
}

func PrintFormatStruct(obj interface{}) {
	empJSON, err := json.Marshal(obj)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%s\n", string(empJSON))
}
