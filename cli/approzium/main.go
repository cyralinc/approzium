package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cyralinc/approzium/authenticator/server/config"
	"github.com/cyralinc/approzium/authenticator/server/credmgrs"
	"github.com/cyralinc/approzium/cli/approzium/commands"
	"github.com/cyralinc/approzium/cli/approzium/util"
)

// Generally, we should add a help tip:
// - Anytime help is shown
// - Anytime we error exit
func main() {
	if util.ShouldGiveTopLevelHelp() {
		util.GiveTopLevelHelp()
		util.AddHelpTip()
		return
	}

	if util.ShouldGiveObjectLevelHelp() {
		if err := util.GiveObjectLevelHelp(commands.Registry); err != nil {
			fmt.Println("Error giving help for object: " + err.Error())
		}
		util.AddHelpTip()
		return
	}

	conf, err := config.Parse()
	if err != nil {
		fmt.Println("Error parsing config: " + err.Error())
		util.AddHelpTip()
		return
	}
	credMgr, err := credmgrs.RetrieveConfigured(util.SilentLogger, conf)
	if err != nil {
		fmt.Println("Unable to select credential manager: " + err.Error())
		util.AddHelpTip()
		return
	}

	obj, err := util.ParseObject(os.Args[1])
	if err != nil {
		fmt.Println("Unable to parse object: " + err.Error())
		util.AddHelpTip()
		return
	}

	verb, err := util.ParseVerb(os.Args[2])
	if err != nil {
		fmt.Println("Unable to parse verb: " + err.Error())
		util.AddHelpTip()
		return
	}

	flags, err := util.ParseFlags()
	if err != nil {
		fmt.Println("Unable to parse flags: " + err.Error())
		util.AddHelpTip()
		return
	}

	userInputReader := bufio.NewReader(os.Stdin)
	for _, command := range commands.Registry {
		if command.Object != obj {
			continue
		}
		if command.Verb != verb {
			continue
		}
		if flags.Help.IsTrue() {
			util.GiveCommandLevelHelp(command)
			util.AddHelpTip()
			return
		}
		if err := command.Handler(flags, credMgr, userInputReader); err != nil {
			fmt.Println("Error executing command: " + err.Error())
			util.AddHelpTip()
			return
		}
		break
	}
}
