package jsoncsv

import (
	"encoding/json"
	"testing"
)

// func Test_mapFromID(t *testing.T) {
// 	ids := [][]string{
// 		{"id"},
// 		{"sections-0", "sections-3", "paragraphs-1", "text"},
// 		{"sections-0", "sections-4", "title"},
// 		{"sections-2", "title"},
// 	}
// 	id := ids[2]

// Output:
// id: "",
// sections: [
// 	{
// 		"title": "",
// 		section: [
// 			{},
// 			{},
// 			{},
// 			{
// 				paragraphs: [
// 					{},
// 					{
// 						text: ""
// 					}
// 				]
// 			},
// 			{
// 				title: ""
// 			}
// 		]
// 	},
// 	{},
// 	{
// 		title: ""
// 	},
// ]

// 	m := mapFromID(id, ids, 0)

// 	t.Logf("%+v\n", m)
// 	t.Fail()
// }

// func Test_iterByColumn(t *testing.T) {
// 	ids := [][]string{
// 		{"id"},
// 		{"sections-0", "sections-3", "paragraphs-1", "text"},
// 		{"sections-0", "sections-4", "title"},
// 		{"sections-2", "title"},
// 	}

// 	m := iterByColumn(ids, 0)

// 	b, err := json.MarshalIndent(m, "", "  ")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("json: %+s\n", b)

// 	// t.Logf("map: %+v\n", m)
// 	t.Fail()
// }

func TestTreeFrom2DMatrix(t *testing.T) {
	records := [][]string{
		{"id", "0123"},
		{"tag/title", "tag title"},
		{"tag/name", "tag name"},
		{"sections-0/sections-3/paragraphs-1/text", "hi there"},
		{"sections-0/sections-4/title", "main title"},
		{"sections-1/paragraphs-0/text", "paragraph of something"},
	}

	wantm := map[string]interface{}{
		"id": "0123",
		"tag": map[string]interface{}{
			"title": "tag title",
			"name":  "tag name",
		},
		"sections": []map[string]interface{}{
			{"sections": []map[string]interface{}{
				nil,
				nil,
				nil,
				{"paragraphs": []map[string]interface{}{
					nil,
					{"text": "hi there"},
				}},
				{"title": "main title"},
			}},
			{"paragraphs": []map[string]interface{}{
				{"text": "paragraph of something"},
			}},
		},
	}
	want, err := json.Marshal(wantm)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("want: %+v\n", string(want))

	ids := idsFromRecords(records)
	t.Logf("ids: %+v\n", ids)

	gotMap := treeFrom2DMatrix(records, ids, "")
	got, err := json.Marshal(gotMap)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("got: %+s\n", got)

	if string(want) != string(got) {
		t.Fail()
	}
}
