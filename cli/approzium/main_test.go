package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/cyralinc/approzium/cli/approzium/util"
	vault "github.com/hashicorp/vault/api"
)

var origStdout *os.File

type testCase struct {
	Name     string
	Args     []string
	Expected string
}

func TestLocalFile(t *testing.T) {
	os.Setenv("APPROZIUM_SECRETS_MANAGER", "local")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	project := "github.com/cyralinc/approzium"
	i := strings.Index(wd, project)
	pathToSecretsYaml := wd[:i] + project + "/authenticator/server/testing/secrets.yaml"
	os.Setenv("APPROZIUM_LOCAL_FILE_PATH", pathToSecretsYaml)

	testCases := []*testCase{
		{
			Name:     "first list",
			Args:     []string{"approzium", "passwords", "list"},
			Expected: "passwords-list.txt",
		},
		{
			Name:     "read with no flags is the same as list",
			Args:     []string{"approzium", "passwords", "read"},
			Expected: "passwords-list.txt",
		},
		{
			Name:     "read with a flag filters results",
			Args:     []string{"approzium", "passwords", "read", "-port=3306"},
			Expected: "passwords-read.txt",
		},
		{
			Name:     "deleting one record works",
			Args:     []string{"approzium", "passwords", "delete", "-host=foo", "--force"},
			Expected: "passwords-delete.txt",
		},
		{
			Name:     "the deleted record is really gone",
			Args:     []string{"approzium", "passwords", "list"},
			Expected: "passwords-list-after-delete.txt",
		},
		{
			Name: "writing the original record works",
			Args: []string{"approzium", "passwords", "write",
				"-host=foo", "-port=5432", "-user=bob", "-password=foo",
				"-grant-access-to=arn:partition:service:region:account-id:arn-thats-not-mine"},
			Expected: "passwords-write.txt",
		},
		{
			Name:     "we have the same records as at the start",
			Args:     []string{"approzium", "passwords", "list"},
			Expected: "passwords-list.txt",
		},
	}

	for _, tstCase := range testCases {
		t.Run(tstCase.Name, func(t *testing.T) {
			ensureExpectedOutputForArgs(t, tstCase)
		})
	}
}

func TestVault(t *testing.T) {
	// Ensure we have the necessary test environment.
	addr := os.Getenv(vault.EnvVaultAddress)
	if addr == "" {
		t.Skip("skipping because VAULT_ADDR is unset")
	}

	os.Setenv("APPROZIUM_SECRETS_MANAGER", "vault")
	testCases := []*testCase{
		{
			Name: "writing a record works",
			Args: []string{"approzium", "passwords", "write",
				"-host=foo", "-port=5432", "-user=bob", "-password=foo",
				"-grant-access-to=buzz"},
			Expected: "passwords-write.txt",
		},
		{
			Name:     "it was really written",
			Args:     []string{"approzium", "passwords", "list"},
			Expected: "vault-passwords-list.txt",
		},
		{
			Name:     "deleting a record works",
			Args:     []string{"approzium", "passwords", "delete", "--f"},
			Expected: "passwords-delete.txt",
		},
		{
			Name:     "the record is gone",
			Args:     []string{"approzium", "passwords", "list"},
			Expected: "passwords-empty-list.txt",
		},
	}

	for _, tstCase := range testCases {
		t.Run(tstCase.Name, func(t *testing.T) {
			ensureExpectedOutputForArgs(t, tstCase)
		})
	}
}

func TestTopLevelHelp(t *testing.T) {
	testCases := []*testCase{
		{
			Name:     "alone",
			Args:     []string{"approzium"},
			Expected: "top-level-help.txt",
		},
		{
			Name:     "with --h",
			Args:     []string{"approzium", "--h"},
			Expected: "top-level-help.txt",
		},
		{
			Name:     "with -help=true",
			Args:     []string{"approzium", "-help=true"},
			Expected: "top-level-help.txt",
		},
	}

	for _, tstCase := range testCases {
		t.Run(tstCase.Name, func(t *testing.T) {
			ensureExpectedOutputForArgs(t, tstCase)
		})
	}
}

func TestObjectLevelHelp(t *testing.T) {
	testCases := []*testCase{
		{
			Name:     "alone",
			Args:     []string{"approzium", "passwords"},
			Expected: "passwords-level-help.txt",
		},
		{
			Name:     "with --h",
			Args:     []string{"approzium", "passwords", "--h"},
			Expected: "passwords-level-help.txt",
		},
		{
			Name:     "with -help=true",
			Args:     []string{"approzium", "passwords", "-help=true"},
			Expected: "passwords-level-help.txt",
		},
	}

	for _, tstCase := range testCases {
		t.Run(tstCase.Name, func(t *testing.T) {
			ensureExpectedOutputForArgs(t, tstCase)
		})
	}
}

func TestCommandLevelHelp(t *testing.T) {
	testCases := []*testCase{
		{
			Name:     "with --h",
			Args:     []string{"approzium", "passwords", "read", "--h"},
			Expected: "passwords-read-help.txt",
		},
		{
			Name:     "with -help=true",
			Args:     []string{"approzium", "passwords", "read", "-help=true"},
			Expected: "passwords-read-help.txt",
		},
	}

	for _, tstCase := range testCases {
		t.Run(tstCase.Name, func(t *testing.T) {
			ensureExpectedOutputForArgs(t, tstCase)
		})
	}
}

func ensureExpectedOutputForArgs(t *testing.T, tstCase *testCase) {
	outputFile, done := redirectStdout(t)
	os.Args = tstCase.Args
	main()
	expected, err := ioutil.ReadFile("fixtures/" + tstCase.Expected)
	if err != nil {
		t.Fatal(err)
	}
	result, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	if clean(result) != clean(expected) {
		t.Fatalf("unexpected response: \n%s", result)
	}
	done()
}

// redirectStdout should be called before each run of main()
func redirectStdout(t *testing.T) (fileName string, done func()) {
	origStdout = os.Stdout
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = tmpFile
	return tmpFile.Name(), func() {
		// Close and delete the file we were using to capture
		// stdout. Restore to normal stdout.
		if err := tmpFile.Close(); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Fatal(err)
		}
		os.Stdout = origStdout

		// Reset the values of all the flags so they won't
		// interact between tests.
		util.Host.Value = ""
		util.Port.Value = ""
		util.User.Value = ""
		util.Password.Value = ""
		util.GrantAccessTo.Value = ""
		util.Force.Value = false
		util.Help.Value = false
	}
}

func clean(b []byte) string {
	cleaned := strings.ReplaceAll(string(b), " ", "")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	return cleaned
}
