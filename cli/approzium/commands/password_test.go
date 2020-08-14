package commands

import (
	"github.com/cyralinc/approzium/authenticator/server/credmgrs"
	"github.com/cyralinc/approzium/cli/approzium/util"
	"testing"
)

func TestFilterByFlags(t *testing.T) {
	dbCreds := []*credmgrs.DBCred{
		{
			AccessGrantedTo: []string{"steve"},
		},
		{
			AccessGrantedTo: []string{"page"},
		},
		{
			AccessGrantedTo: []string{"steve", "ellen"},
		},
	}

	result, err := filterByUserFlags(dbCreds, &util.Flags{
		DBHost: &util.StrFlag{},
		DBPort: &util.StrFlag{},
		DBUser: &util.StrFlag{},
		GrantAccessTo: &util.StrFlag{
			Value: "ellen",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatal("expected 1 matches")
	}
}
