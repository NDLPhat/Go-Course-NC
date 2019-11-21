package main

import (
	"encoding/json"
	"fmt"
)

type StudentInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	ClassName string `json:"class_name"`
}

type StudentInfos []StudentInfo

func main() {
	inputJson := `[
		{"first_name": "Victor", "last_name": "Nguyen", "age": 100, "class_name":"golang", "test": "123", "alo": "alo here"},
		{"first_name": "Anh", "last_name": "Dinh", "age":200, "class_name":"golang"}
	]`
	var studentInfos StudentInfos

	err := studentInfos.parse(inputJson)
	if err != nil {
		// exit and printout error
		panic(err)
	}
	// printout data
	for _, v := range studentInfos {
		fmt.Printf("%#v \r\n", v)
	}
}

func (inf *StudentInfos) parse(inputJson string) error {
	// parse json-encoded data
	err := json.Unmarshal([]byte(inputJson), &inf)
	return err
}
