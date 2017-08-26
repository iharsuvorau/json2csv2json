// JSON to CSV and CSV to JSON converter.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

func main() {
	// get data

	// TODO: walk through all json files in a directory recursively
	f, err := os.Open("samples/biodiversity.ru.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data := make(map[string]interface{})

	d := json.NewDecoder(f)
	if err = d.Decode(&data); err != nil {
		log.Fatal(err)
	}

	// convert to plain id:value map
	out := make(map[string]interface{})
	for k, v := range data {
		iterate(k, v, out)
	}

	// save to csv
	csvFile, err := os.Create("output.tsv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	if _, err = csvFile.WriteString("id\t value\n"); err != nil {
		log.Println(err)
	} // headers

	for id, value := range out {
		// TODO: strip special symbols

		// write only fields: title, text, lead, teaser, description
		if strings.Contains(id, "title") || strings.Contains(id, "text") || strings.Contains(id, "lead") || strings.Contains(id, "teaser") || strings.Contains(id, "description") {
			if _, err = csvFile.WriteString(fmt.Sprintf("%s\t \"%v\"\n", id, value)); err != nil {
				log.Println(err)
			} // rows
		}
	}
}

func iterate(parent string, data interface{}, output map[string]interface{}) {
	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.String:
		// fmt.Printf("%s (string)\n", parent)
		output[parent] = v.Interface()
	case reflect.Float64:
		// fmt.Printf("%s (float64)\n", parent)
		output[parent] = v.Interface()
	case reflect.Slice:
		// fmt.Printf("%s (slice of %v)\n", parent, v.Len())
		for i := 0; i < v.Len(); i++ {
			iterate(
				fmt.Sprintf("%s-%d", parent, i),
				reflect.ValueOf(v.Index(i).Interface()).Interface(),
				output,
			)
		}
	case reflect.Map:
		// fmt.Printf("%s (map with %v)\n", parent, v.MapKeys())
		for _, k := range v.MapKeys() {
			iterate(
				parent+"/"+k.String(),
				v.MapIndex(k).Interface(),
				output,
			)
		}
	}
}
