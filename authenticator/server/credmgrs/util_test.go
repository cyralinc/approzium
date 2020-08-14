package credmgrs

import (
	"testing"
)

func TestDeleteIfExists(t *testing.T) {
	toDel := &DBCred{
		Host:            "localhost",
		Port:            "5432",
		User:            "engineering",
		Password:        "(*(*#&(*$&(*#&",
		AccessGrantedTo: []string{"bob", "steve"},
	}
	all := []*DBCred{
		{
			Host:            "localhost",
			Port:            "5432",
			User:            "engineering",
			Password:        "(*(*#&(*$&(*#&",
			AccessGrantedTo: []string{"bob", "steve"},
		},
	}
	all = deleteIfExists(toDel, all)
	if len(all) != 0 {
		t.Fatal("expected 0 records")
	}
}
