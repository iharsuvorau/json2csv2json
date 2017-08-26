package jsoncsv

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// WalkWokring does the work recursively on a bunch of files.
func WalkWokring(outputDir, outputExt string, delimeter rune) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var outputFilename string

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			outputFilename = strings.Replace(filepath.Base(path), filepath.Ext(path), "", 1) + outputExt
			if err = Work(path, filepath.Join(outputDir, outputFilename), delimeter); err != nil {
				return err
			}
		}

		return nil
	}
}

// Work is a main wrapper function which controls the whole workflow.
func Work(input, output string, delimeter rune) (err error) {
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
// TODO: use encoding/csv
func save(filename string, delimeter rune, output map[string]interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = delimeter

	defer w.Flush()

	if err = w.Write([]string{"id", "value"}); err != nil {
		return err
	} // header

	for id, value := range output {
		// TODO: strip special symbols

		// NOTE: specific requirements, can be ignored for more general usage
		// write only fields: title, text, lead, teaser, description
		if strings.Contains(id, "title") || strings.Contains(id, "text") || strings.Contains(id, "lead") || strings.Contains(id, "teaser") || strings.Contains(id, "description") {
			if err = w.Write([]string{id, fmt.Sprintf("%v", value)}); err != nil {
				return err
			} // rows
		}
	}

	return nil
}

//

func Read(filename string, delimeter rune) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}

	r := csv.NewReader(f)
	r.Comma = delimeter

	records, err := r.ReadAll()
	if err != nil {
		return
	}

	m := make(map[string]interface{})
	for _, rec := range records {
		// extract json path from id
		fmt.Println(rec[0])

		Analyze(rec[0], m)

		// add the value
		// rec[1]
	}
	fmt.Printf("%+v\n", m)

	return nil
}

func Analyze(id string, m map[string]interface{}) {
	parts := strings.Split(id, "/")

	if len(parts) == 1 {
		m[parts[0]] = "data"
		return
	}

	parts1 := strings.Split(parts[0], "-")
	if len(parts1) > 1 {
		m[parts1[0]] = []map[string]interface{}{}
		Analyze(strings.Join(parts[1:], "/"), m[parts1[0]])
	} else {
		m[parts[0]] = parts[1:]
	}
}
