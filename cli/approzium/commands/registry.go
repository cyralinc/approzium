package commands

import "github.com/cyralinc/approzium/cli/approzium/util"

var Registry = []*util.Command{
	cmdPasswordsRead(),
	cmdPasswordsWrite(),
	cmdPasswordsDelete(),
	cmdPasswordsList(),
}
