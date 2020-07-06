package testing

import "os"

const EnvVarAcceptanceTests = "APPROZIUM_ACC"

func ShouldRunAcceptanceTests() bool {
	return os.Getenv(EnvVarAcceptanceTests) != ""
}
