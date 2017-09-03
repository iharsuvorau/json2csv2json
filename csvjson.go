package jsoncsv

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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

	for _, rec := range records {
		fmt.Println(rec[0])
	}

	// extract json paths (ids)
	// ids := idsFromRecords(records)
	// fmt.Printf("%+v\n", ids)

	// m := mapFromIds(ids)
	// fmt.Printf("%+v\n", m)

	return nil
}

func treeFrom2DMatrix(ids [][]string) map[string]interface{} {
	m := make(map[string]interface{})

	// keys with index 0 are all from the same node, so add them as keys and check each for final data or further run
	// go throuhg all keys with index 0 and find slices or maps to go further
	//   if found a slice, then find out the amount of elements, create a slice with len(elements)
	//     for the length of the slice run through each index and collect keys with the keyname + index
	//     remove first elements and pass it the func, returned values added to the slice at index

	for _, v := range ids {
		if isData(v) {
			m[v[0]] = ""
		} else if isMap(v) {
			cm := make(map[string]interface{})
			m[v[0]] = cm
		} else if isSlice(v) {
			key := cleanName(v[0])

			// create a slice if it doesnt' exist yet
			if _, ok := m[key]; !ok {
				// fmt.Println("slice create:", key)
				size, err := findSliceLength(idsByColumn(ids, 0), key)
				if err != nil {
					log.Fatal(err)
				}

				sl := make([]map[string]interface{}, size)
				m[key] = sl

				// populate each element of the slice with a map
				for i := 0; i < size; i++ {
					// fmt.Println("slice index", i)
					nids := [][]string{}
					for _, id := range ids {
						if isSlice(id) {
							_, idx := getKeyIndex(id[0])
							if i == idx {
								nids = append(nids, id[1:]) // delete the root key
							}
						}
					}
					// fmt.Println("found", nids)
					if len(nids) > 0 {
						sl[i] = treeFrom2DMatrix(nids)
					}
				}
			} else {
				// fmt.Println("slice already created:", key, "skip")
			}

		}
	}

	return m
}

func idsFromRecords(records [][]string) [][]string {
	ids := make([][]string, len(records))

	for i := range records {
		ids[i] = strings.Split(records[i][0], "/")
	}

	return ids
}

func cleanName(s string) string {
	if strings.Contains(s, "-") {
		return strings.Split(s, "-")[0]
	}

	return s
}

func getKeyIndex(s string) (key string, idx int) {
	parts := strings.Split(s, "-")
	key = parts[0]
	idx, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal(err)
	}
	return
}

func keysFromColumn(col []string) map[string]string {
	keys := make(map[string]string)
	var t string
	for _, v := range col {
		t = "data"
		if strings.Contains(v, "-") {
			size, err := findSliceLength(col, cleanName(v))
			if err != nil {
				log.Fatal(err)
			}
			t = fmt.Sprintf("slice:%d", size)
		}
		keys[cleanName(v)] = t
	}

	return keys
}

// isData return true if the ID has no more children and it most probably contains final data.
func isData(id []string) bool {
	if len(id) == 1 {
		return true
	}
	return false
}

// isMap returns true if IDs of a certain depth contain non-slice IDs (without '-'),
// that means that it will be a map and not a slice.
func isMap(id []string) bool {
	if len(id) > 1 {
		if !strings.Contains(id[0], "-") {
			return true
		}
	}
	return false
}

func isSlice(id []string) bool {
	if strings.Contains(id[0], "-") {
		return true
	}
	return false
}

func findSliceLength(column []string, key string) (int, error) {
	var maxIndex int

	for _, idPart := range column {
		if strings.Contains(idPart, key) {
			index := strings.Split(idPart, "-")[1]
			idx, err := strconv.Atoi(index)
			if err != nil {
				return maxIndex + 1, err
			}

			if idx > maxIndex {
				maxIndex = idx
			}
		}
	}

	// length is bigger than index by 1
	return maxIndex + 1, nil
}

// idsByColumn extracts ids by a column.
func idsByColumn(ids [][]string, colNumber int) (column []string) {
	column = []string{}

	for _, id := range ids {
		if len(id)-1 >= colNumber {
			column = append(column, id[colNumber])
		}
	}

	return
}
