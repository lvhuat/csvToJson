package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var FieldTypeMap = map[string]string{
	"string": "string", "int": "int", "double": "double",
}
var path = flag.String("file", "symbols.csv", "symbol file")

func jsonField(tp, value string) string {
	var s string
	switch tp {
	case "string":
		s = fmt.Sprintf("\"%s\"", value)
	case "int":
		s = value
	case "double":
		s = value
	}

	return s
}

func main() {
	flag.Parse()
	file, err := os.Open(*path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	fieldNames, err := reader.Read()
	if err != nil {
		panic(err)
	}

	fieldTypes, err := reader.Read()
	if err != nil {
		panic(err)
	}

	fieldCount := len(fieldNames)
	if len(fieldTypes) != fieldCount {
		panic("header error")
	}

	for _, t := range fieldTypes {
		_, exist := FieldTypeMap[t]
		if !exist {
			panic("type not exist")
		}
	}

	start := true
	var s string
	var array bool
	s += "["
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if len(fields) != fieldCount {
			panic("bad data row count")
		}

		if !start {
			s += ","
		}

		start = false

		s += "{"

		for i := 0; i < len(fieldTypes); i++ {
			value := fields[i]
			ft := fieldTypes[i]
			name := fieldNames[i]
			arrayItem := strings.Index(name, "[]") == 0
			if arrayItem {
				name = name[2:]
			}
			if !arrayItem {
				// Last field is array
				if array {
					s += "]"
					array = false
				}

				// This field is not the object start
				if i != 0 {
					s += ","
				}

				s += fmt.Sprintf("\"%s\":%s", name, jsonField(ft, value))
			} else {
				// Add array item
				if array {
					s += ","
				} else {

					// This field is not the object start
					if i != 0 {
						s += ","
					}

					// Start the array
					s += fmt.Sprintf("\"%s\":[", name)
					array = true
				}
				s += jsonField(ft, value)
			}
		}

		// Last item is array
		if array {
			s += "]"
			array = false
		}
		//
		s += "}"
	}

	s += "]"

	name := strings.TrimSuffix(*path, ".csv")
	jsonFile, err := os.OpenFile(name+".json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		panic("create json file error " + err.Error())
	}

	defer jsonFile.Close()

	_, err = jsonFile.WriteString(s)
	if err != nil {
		panic("write json file error " + err.Error())
	}
}
