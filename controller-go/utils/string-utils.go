package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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

func ConvertToUint(src string) uint {
	i, err := strconv.ParseUint(src, 10, 64)
	if err == nil {
		return uint(i)
	}

	return 0
}
