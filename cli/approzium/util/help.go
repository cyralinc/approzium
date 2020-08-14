package util

import (
	"fmt"
	"os"
)

func ShouldGiveTopLevelHelp() bool {
	if len(os.Args) == 1 {
		// They only entered "approzium" with nothing following it.
		return true
	}
	if len(os.Args) > 1 && IsHelpFlag(os.Args[1]) {
		// They entered something like "approzium --h".
		return true
	}
	return false
}

func GiveTopLevelHelp() {
	fmt.Println("approzium supports the following sub-commands:")
	for _, obj := range AllObjects() {
		fmt.Println("  " + obj)
	}

	fmt.Println("")
	fmt.Println("Example commands:")
	fmt.Println("  approzium passwords list")
	fmt.Println("  approzium passwords read --h")
	fmt.Println("  approzium passwords read -port=5432")
}

func ShouldGiveObjectLevelHelp() bool {
	if len(os.Args) == 2 {
		// They only entered something like "approzium passwords" with nothing following it.
		return true
	}
	if len(os.Args) > 2 && IsHelpFlag(os.Args[2]) {
		// They entered something like "approzium passwords --h".
		return true
	}
	return false
}

func GiveObjectLevelHelp(commands []*Command) error {
	obj, err := ParseObject(os.Args[1])
	if err != nil {
		return err
	}

	fmt.Printf("approzium %s supports the following methods:\n", obj)
	for _, command := range commands {
		if command.Object != obj {
			continue
		}
		fmt.Println("  " + command.Verb)
	}
	return nil
}

func GiveCommandLevelHelp(command *Command) {
	fmt.Printf("approzium %s %s supports the following flags:\n", command.Object, command.Verb)
	for _, supportedFlag := range command.SupportedFlags {
		if supportedFlag.IsRequired {
			fmt.Printf("  -%s (required)\n", supportedFlag.Name)
		} else {
			fmt.Printf("  -%s (optional)\n", supportedFlag.Name)
		}
	}
}
