package util

import (
	"bufio"
	"fmt"

	"github.com/cyralinc/approzium/authenticator/server/credmgrs"
)

type Object string

const (
	// When adding objects, please also add them to the AllObjects
	// and ParseObject functions below.
	ObjPasswords Object = "passwords"
)

func AllObjects() []Object {
	return []Object{ObjPasswords}
}

func ParseObject(s string) (Object, error) {
	switch s {
	case "passwords":
		return ObjPasswords, nil
	default:
		return "", fmt.Errorf("unrecognized object: %s", s)
	}
}

type Verb string

const (
	VerbRead   Verb = "read"
	VerbWrite       = "write"
	VerbDelete      = "delete"
	VerbList        = "list"
)

func ParseVerb(s string) (Verb, error) {
	switch s {
	case "read":
		return VerbRead, nil
	case "write":
		return VerbWrite, nil
	case "delete":
		return VerbDelete, nil
	case "list":
		return VerbList, nil
	default:
		return "", fmt.Errorf("unrecognized verb: %s", s)
	}
}

type Command struct {
	// Required.
	Object Object

	// Required.
	Verb Verb

	// How to handle this command if it's seen.
	Handler func(flags *Flags, credMgr credmgrs.CredentialManager, userInputReader *bufio.Reader) error

	// For user help output, SupportedFlags is a list of the flags processed
	// by this command. No need to list the "help" flag, it's support everywhere.
	SupportedFlags []*FlagInfo
}

type FlagInfo struct {
	Name       string
	IsRequired bool
}
