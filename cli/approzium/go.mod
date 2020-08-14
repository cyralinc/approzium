module github.com/cyralinc/approzium/cli/approzium

go 1.14

replace github.com/cyralinc/approzium/authenticator => ../../authenticator

require (
	github.com/cyralinc/approzium/authenticator v0.0.0-20200812230944-f454ad946be3
	github.com/hashicorp/vault/api v1.0.4
	github.com/sirupsen/logrus v1.6.0
)
