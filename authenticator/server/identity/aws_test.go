package identity

import (
	"strings"
	"testing"

	pb "github.com/approzium/approzium/authenticator/server/protos"
	testtools "github.com/approzium/approzium/authenticator/server/testing"
)

var testEnv = &testtools.AwsEnv{}

func TestVerifyService(t *testing.T) {
	signedGetCallerIdentity, err := testEnv.SignedGetCallerIdentity(t)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		TestName                string
		SignedGetCallerIdentity string
		ExpectedArn             string
		ExpectErr               bool
	}{
		{
			TestName:                "Sunny path, regular arn",
			SignedGetCallerIdentity: signedGetCallerIdentity,
			ExpectedArn:             testEnv.ClaimedArn(),
			ExpectErr:               false,
		},
		{
			TestName:                "Empty values",
			SignedGetCallerIdentity: "",
			ExpectErr:               true,
		},
		{
			TestName:                "Malicious base URL injected",
			SignedGetCallerIdentity: strings.ReplaceAll(signedGetCallerIdentity, "sts", "somewhere-else"),
			ExpectErr:               true,
		},
		{
			TestName:                "Different call than GetCallerIdentity",
			SignedGetCallerIdentity: strings.ReplaceAll(signedGetCallerIdentity, "GetCallerIdentity", "GetSessionToken"),
			ExpectErr:               true,
		},
	}
	a := &aws{}
	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			verifiedARN, err := a.getAwsIdentity(testCase.SignedGetCallerIdentity, pb.ClientLanguage_GO)
			if testCase.ExpectErr {
				if err == nil {
					t.Fatal("expected err")
				} else {
					// We expected an error and received it, so we've succeeded
					// and there's nothing else to do here.
					return
				}
			}
			if err != nil {
				t.Fatal(err)
			}

			// We don't expect an error. Let's make sure we got the expected response.
			match, err := a.arnsMatch(testCase.ExpectedArn, verifiedARN)
			if err != nil {
				t.Fatal(err)
			}
			if !match {
				t.Fatalf("expected %s but received %s", testCase.ExpectedArn, verifiedARN)
			}
		})
	}
}
