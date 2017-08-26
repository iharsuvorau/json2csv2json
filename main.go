// JSON to CSV and CSV to JSON converter.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

func main() {
	// get data

	// TODO: walk through all json files in a directory recursively
	fmt.Printf("%+v\n", os.Args[1:])
}

// work is a main wrapper function which controls the whole workflow.
func work(input, output, delimeter string) (err error) {
	data := make(map[string]interface{})
	if err = build(input, data); err != nil {
		return
	}

	if err = save(output, delimeter, data); err != nil {
		return
	}

	return nil
}

// build reads a source file and reflects on its content to build a plain map of id:value pairs.
func build(filename string, output map[string]interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	data := make(map[string]interface{})

	d := json.NewDecoder(f)
	if err = d.Decode(&data); err != nil {
		return err
	}

	// convert to plain id:value map
	for k, v := range data {
		iterate(k, v, output)
	}

	return nil
}

// iterate goes recursively through an interface value, composes ids and values and adds them to the output map.
func iterate(parent string, data interface{}, output map[string]interface{}) {
	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.String:
		// fmt.Printf("%s (string)\n", parent)
		output[parent] = v.Interface()
		return
	case reflect.Float64:
		// fmt.Printf("%s (float64)\n", parent)
		output[parent] = v.Interface()
		return
	case reflect.Slice:
		// fmt.Printf("%s (slice of %v)\n", parent, v.Len())
		for i := 0; i < v.Len(); i++ {
			iterate(
				fmt.Sprintf("%s-%d", parent, i),
				reflect.ValueOf(v.Index(i).Interface()).Interface(),
				output,
			)
		}
		return
	case reflect.Map:
		// fmt.Printf("%s (map with %v)\n", parent, v.MapKeys())
		for _, k := range v.MapKeys() {
			iterate(
				parent+"/"+k.String(),
				v.MapIndex(k).Interface(),
				output,
			)
		}
		return
	}
}

// save saves an output map to a CSV-, TSV-file. The format of the file is dependent on the delimeter
// provided to the function: "\t" or ",".
func save(filename, delimiter string, output map[string]interface{}) error {
	csvFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	if _, err = csvFile.WriteString("id\t value\n"); err != nil {
		return err
	} // headers

	for id, value := range output {
		// TODO: strip special symbols

		// NOTE: specific requirements, can be ignored for more general usage
		// write only fields: title, text, lead, teaser, description
		if strings.Contains(id, "title") || strings.Contains(id, "text") || strings.Contains(id, "lead") || strings.Contains(id, "teaser") || strings.Contains(id, "description") {
			if _, err = csvFile.WriteString(fmt.Sprintf("%s%s \"%v\"\n", id, delimiter, value)); err != nil {
				return err
			} // rows
		}
	}

	return nil
}
