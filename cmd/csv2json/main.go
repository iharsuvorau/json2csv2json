// CSV to JSON converter.
// NOTE: JSON keys must not contain '-' because it's used as a slice indicator.

package main

import (
	"log"

	"github.com/iharsuvorau/jsoncsv"
)

func main() {
	err := jsoncsv.Read("samples/biodiversity.ru.tsv", '\t')
	if err != nil {
		log.Fatal(err)
	}
}
