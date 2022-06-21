package malapi

import (
	"testing"
)

func TestExportToCsv(t *testing.T) {
	apis := []Api{
		{
			Name:              "Test1",
			Description:       "Description 1",
			Library:           "lib1",
			AssociatedAttacks: nil,
			Documentation:     "Doc 1",
			Created:           "2020-02-01",
			LastUpdate:        "2020-02-01",
			Credits:           "jason",
		},
		{
			Name:              "Test2",
			Description:       "Description \"jason\" 1",
			Library:           "lib2",
			AssociatedAttacks: []string{"attack1", "attack2"},
			Documentation:     "Doc 2",
			Created:           "2020-02-01",
			LastUpdate:        "2020-02-01",
			Credits:           "james",
		}}

	err := ExportToCsv(&apis)

	if err != nil {
		t.Error(err)
	}

}

func TestGet(t *testing.T) {
	
}
