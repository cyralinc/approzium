package testing

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

const envVarTestRole = "TEST_ASSUMABLE_ARN"

// This allows us to only get the signedGetCallerIdentity string once, but
// to reuse it throughout tests through the testEnv variable, reducing load
// on AWS.
type AwsEnv struct {
	signedGetCallerIdentity string
}

func (e *AwsEnv) ClaimedArn() string {
	return os.Getenv(envVarTestRole)
}

func (e *AwsEnv) SignedGetCallerIdentity(t *testing.T) (string, error) {

	if os.Getenv(envVarTestRole) == "" {
		t.Skip(fmt.Sprintf("skipping because %s is unset", envVarTestRole))
	}

	// If it's cached, return it.
	if e.signedGetCallerIdentity != "" {
		return e.signedGetCallerIdentity, nil
	}

	// If it's uncached, get it, cache it, and return it.
	sess, err := session.NewSession()
	if err != nil {
		return "", err
	}
	creds := stscreds.NewCredentials(sess, os.Getenv(envVarTestRole))

	// Create service client value configured for credentials
	// from assumed role.
	svc := sts.New(sess, &aws.Config{Credentials: creds})

	req, _ := svc.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{})
	signedGetCallerIdentity, err := req.Presign(time.Minute * 15)
	if err != nil {
		return "", err
	}
	e.signedGetCallerIdentity = signedGetCallerIdentity
	return e.signedGetCallerIdentity, nil
}
