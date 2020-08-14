package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	Host = &StrFlag{
		Name:  "host",
		Usage: "database host",
	}
	Port = &StrFlag{
		Name:  "port",
		Usage: "database port",
	}
	User = &StrFlag{
		Name:  "user",
		Usage: "database user",
	}
	Password = &StrFlag{
		Name:  "password",
		Usage: "database password",
	}
	GrantAccessTo = &StrFlag{
		Name:  "grant-access-to",
		Usage: "the comma-separated ARNs of users who should be able to access this record",
	}
	Force = &BoolFlag{
		Name:  "force",
		Usage: "force command to be accepted",
	}
	Help = &BoolFlag{
		Name:  "help",
		Usage: "gives tips about how to use the CLI",
	}
)

func IsHelpFlag(s string) bool {
	switch s {
	// Note that these should correspond with the help flags in the switch below.
	case "-help=true", "-h=true", "--h", "--help":
		return true
	default:
		return false
	}
}

func AddHelpTip() {
	fmt.Println("")
	fmt.Println("Get help with any command by adding '--h' or '--help'.")
}

func ParseFlags() (*Flags, error) {
	// os.Args[0] is a bunch of gobbledygook.
	// os.Args[1] is an object like "passwords".
	// os.Args[2] is a verb like "read".
	for i := 3; i < len(os.Args); i++ {
		field := os.Args[i]

		fieldName := ""
		fieldValue := ""
		if strings.HasPrefix(field, "--") {
			// Attempt to interpret commands formatted like "--force".
			fieldName = strings.TrimPrefix(field, "--")
			fieldValue = "true"
		} else {
			// Attempt to interpret commands formatted like "-host=foo".
			fieldNameValue := strings.Split(field, "=")
			if len(fieldNameValue) != 2 {
				return nil, fmt.Errorf("unexpected field: %s", field)
			}
			fieldName = strings.TrimPrefix(fieldNameValue[0], "-")
			fieldValue = fieldNameValue[1]
		}

		switch fieldName {
		case "host":
			Host.Value = fieldValue
		case "port":
			Port.Value = fieldValue
		case "user":
			User.Value = fieldValue
		case "password":
			Password.Value = fieldValue
		case "grant-access-to":
			GrantAccessTo.Value = fieldValue
		case "f", "force":
			v, err := strconv.ParseBool(fieldValue)
			if err != nil {
				return nil, err
			}
			Force.Value = v
		case "h", "help":
			v, err := strconv.ParseBool(fieldValue)
			if err != nil {
				return nil, err
			}
			Help.Value = v
		default:
			return nil, fmt.Errorf("unrecognized field %s", fieldName)
		}
	}
	return &Flags{
		DBHost:        Host,
		DBPort:        Port,
		DBUser:        User,
		DBPassword:    Password,
		GrantAccessTo: GrantAccessTo,
		Force:         Force,
		Help:          Help,
	}, nil
}

// Flags exist so the actual flags that are processed at runtime
// can be passed around.
type Flags struct {
	DBHost, DBPort, DBUser, DBPassword, GrantAccessTo *StrFlag
	Force, Help                                       *BoolFlag
}

type StrFlag struct {
	Name, Usage, Value string
}

func (s *StrFlag) String() string {
	return s.Value
}

type BoolFlag struct {
	Name  string
	Usage string

	Value bool
}

func (b *BoolFlag) IsTrue() bool {
	return b.Value
}
