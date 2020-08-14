package commands

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"github.com/cyralinc/approzium/authenticator/server/credmgrs"
	"github.com/cyralinc/approzium/cli/approzium/util"
)

func cmdPasswordsRead() *util.Command {
	return &util.Command{
		Object: util.ObjPasswords,
		Verb:   util.VerbRead,
		SupportedFlags: []*util.FlagInfo{
			{Name: util.Host.Name, IsRequired: false},
			{Name: util.Port.Name, IsRequired: false},
			{Name: util.User.Name, IsRequired: false},
			{Name: util.GrantAccessTo.Name, IsRequired: false},
		},
		Handler: func(flags *util.Flags, credMgr credmgrs.CredentialManager, _ *bufio.Reader) error {
			dbCreds, err := credMgr.List()
			if err != nil {
				return err
			}
			dbCreds, err = filterByUserFlags(dbCreds, flags)
			if err != nil {
				return err
			}
			credmgrs.Sort(dbCreds)
			for _, dbCred := range dbCreds {
				printCred(dbCred)
			}
			fmt.Println("")
			return nil
		},
	}
}

func cmdPasswordsWrite() *util.Command {
	return &util.Command{
		Object: util.ObjPasswords,
		Verb:   util.VerbWrite,
		SupportedFlags: []*util.FlagInfo{
			{Name: util.Host.Name, IsRequired: true},
			{Name: util.Port.Name, IsRequired: true},
			{Name: util.User.Name, IsRequired: true},
			{Name: util.Password.Name, IsRequired: true},
			{Name: util.GrantAccessTo.Name, IsRequired: true},
		},
		Handler: func(flags *util.Flags, credMgr credmgrs.CredentialManager, userInputReader *bufio.Reader) error {
			dbCred := &credmgrs.DBCred{}
			var err error

			dbCred.Host, err = util.HandleRequiredFlag(flags.DBHost.Value, flags.DBHost.Name, userInputReader)
			if err != nil {
				return err
			}
			dbCred.Port, err = util.HandleRequiredFlag(flags.DBPort.Value, flags.DBPort.Name, userInputReader)
			if err != nil {
				return err
			}
			dbCred.User, err = util.HandleRequiredFlag(flags.DBUser.Value, flags.DBUser.Name, userInputReader)
			if err != nil {
				return err
			}
			dbCred.Password, err = util.HandleRequiredFlag(flags.DBPassword.Value, flags.DBPassword.Name, userInputReader)
			if err != nil {
				return err
			}
			raw, err := util.HandleRequiredFlag(flags.GrantAccessTo.Value, flags.GrantAccessTo.Name, userInputReader)
			if err != nil {
				return err
			}
			// It's natural to want to place a space after a comma. If folks do, we'll
			// just correct it for them on the way in.
			raw = strings.ReplaceAll(raw, " ", "")
			dbCred.AccessGrantedTo = strings.Split(raw, ",")

			if err := credMgr.Write(dbCred); err != nil {
				return err
			}
			fmt.Println("Write successful!")
			return nil
		},
	}
}

func cmdPasswordsDelete() *util.Command {
	return &util.Command{
		Object: util.ObjPasswords,
		Verb:   util.VerbDelete,
		SupportedFlags: []*util.FlagInfo{
			{Name: util.Host.Name, IsRequired: false},
			{Name: util.Port.Name, IsRequired: false},
			{Name: util.User.Name, IsRequired: false},
			{Name: util.GrantAccessTo.Name, IsRequired: false},
			{Name: util.Force.Name, IsRequired: false},
		},
		Handler: func(flags *util.Flags, credMgr credmgrs.CredentialManager, _ *bufio.Reader) error {
			dbCreds, err := credMgr.List()
			if err != nil {
				return err
			}
			dbCreds, err = filterByUserFlags(dbCreds, flags)
			if err != nil {
				return err
			}
			credmgrs.Sort(dbCreds)

			if flags.Force.IsTrue() {
				for _, dbCred := range dbCreds {
					if err := credMgr.Delete(dbCred); err != nil {
						return err
					}
				}
				fmt.Println("Successfully deleted matches, if they existed!")
			} else {
				fmt.Println("This command would delete the following records IN ENTIRETY:")
				for _, dbCred := range dbCreds {
					printCred(dbCred)
				}
				fmt.Printf("\nRun with --force to delete.\n")
			}
			return nil
		},
	}
}

func cmdPasswordsList() *util.Command {
	return &util.Command{
		Object: util.ObjPasswords,
		Verb:   util.VerbList,
		Handler: func(_ *util.Flags, credMgr credmgrs.CredentialManager, _ *bufio.Reader) error {
			dbCreds, err := credMgr.List()
			if err != nil {
				return err
			}
			credmgrs.Sort(dbCreds)
			for _, dbCred := range dbCreds {
				printCred(dbCred)
			}
			fmt.Println("")
			return nil
		},
	}
}

func filterByUserFlags(dbCreds []*credmgrs.DBCred, flags *util.Flags) ([]*credmgrs.DBCred, error) {
	var filtered []*credmgrs.DBCred
	for _, dbCred := range dbCreds {
		dbHost := flags.DBHost.String()
		if dbHost != "" && dbHost != dbCred.Host {
			continue
		}

		dbPort := flags.DBPort.String()
		if dbPort != "" && dbPort != dbCred.Port {
			continue
		}

		dbUser := flags.DBUser.String()
		if dbUser != "" && dbUser != dbCred.User {
			continue
		}

		// We don't support searching by password because for security, we don't
		// read out passwords at the CLI. Matching on them would allow people
		// to derive information from them as well.

		grantAccessTo := flags.GrantAccessTo.String()
		if grantAccessTo != "" {
			grantsSought := strings.Split(grantAccessTo, ",")
			if len(grantsSought) != 1 {
				// We can't search for multiple users at once because the logic
				// becomes murky. Does the user want us to pull out only matches
				// with all of the access grants? Or any of the access grants?
				// Instead, let's keep this simple for everyone involved :-).
				return nil, errors.New("we must search for one access grant at a time, please update your -grant-access-to flag")
			}
			grantToFilterBy := grantsSought[0]
			match := false
			for _, grant := range dbCred.AccessGrantedTo {
				if grant != grantToFilterBy {
					continue
				}
				match = true
				break
			}
			if !match {
				continue
			}
		}
		filtered = append(filtered, dbCred)
	}
	return filtered, nil
}

func printCred(dbCred *credmgrs.DBCred) {
	output := fmt.Sprintf(`  host: %s
  port: %s
  user: %s
  grant-access-to: %s`, dbCred.Host, dbCred.Port, dbCred.User, strings.Join(dbCred.AccessGrantedTo, ","))
	fmt.Printf("\n%s\n", output)
}
