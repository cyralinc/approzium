package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	pb "github.com/approzium/approzium/authenticator/protos"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/google/gofuzz"
	vault "github.com/hashicorp/vault/api"
)

const envVarTestRole = "TEST_IAM_ROLE"

var testEnv = &env{}

// TestAuthenticator_GetPGMD5Hash issues real STS GetCallerIdentity because at the
// time of writing there were no documented limits. Hitting the real API will allow
// us to catch if it changes.
func TestAuthenticator_GetPGMD5Hash(t *testing.T) {
	// These tests rely upon the file back-end, so unset the Vault addr if it exists.
	_ = os.Setenv(vault.EnvVaultAddress, "")
	signedGetCallerIdentity, err := testEnv.SignedGetCallerIdentity(t)
	if err != nil {
		t.Fatal(err)
	}

	authenticator, err := NewAuthenticator()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := authenticator.GetPGMD5Hash(nil, &pb.PGMD5HashRequest{
		Authtype:       pb.AuthType_AWS,
		ClientLanguage: pb.ClientLanguage_GO,
		Dbhost:         "dbmd5",
		Dbport:         "5432",
		Dbuser:         "bob",
		Awsauth: &pb.AWSAuth{
			SignedGetCallerIdentity: signedGetCallerIdentity,
			ClaimedIamArn:           testEnv.ClaimedArn(),
		},
		Salt: []byte{1, 2, 3, 4},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Hash != "d576ce99165615bb3f4331154ca6660c" {
		t.Fatalf("expected %s but received %s", "d576ce99165615bb3f4331154ca6660c", resp.Hash)
	}

	// Now use a bad claimed arn and make sure we fail.
	resp, err = authenticator.GetPGMD5Hash(nil, &pb.PGMD5HashRequest{
		Authtype:       pb.AuthType_AWS,
		ClientLanguage: pb.ClientLanguage_GO,
		Dbhost:         "foo",
		Dbport:         "5432",
		Dbuser:         "bob",
		Awsauth: &pb.AWSAuth{
			SignedGetCallerIdentity: signedGetCallerIdentity,
			ClaimedIamArn:           "arn:partition:service:region:account-id:arn-thats-not-mine",
		},
		Salt: []byte{1, 2, 3, 4},
	})
	if err == nil {
		t.Fatal("using a claimed arn that doesn't belong to me should fail")
	}
}

// TestAuthenticator_GetPGSHA256Hash issues real STS GetCallerIdentity because at the
// time of writing there were no documented limits. Hitting the real API will allow
// us to catch if it changes.
func TestAuthenticator_GetPGSHA256Hash(t *testing.T) {
	// These tests rely upon the file back-end, so unset the Vault addr if it exists.
	_ = os.Setenv(vault.EnvVaultAddress, "")
	signedGetCallerIdentity, err := testEnv.SignedGetCallerIdentity(t)
	if err != nil {
		t.Fatal(err)
	}

	authenticator, err := NewAuthenticator()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := authenticator.GetPGSHA256Hash(nil, &pb.PGSHA256HashRequest{
		Authtype:       pb.AuthType_AWS,
		ClientLanguage: pb.ClientLanguage_GO,
		Dbhost:         "dbsha256",
		Dbport:         "5432",
		Dbuser:         "bob",
		Awsauth: &pb.AWSAuth{
			SignedGetCallerIdentity: signedGetCallerIdentity,
			ClaimedIamArn:           testEnv.ClaimedArn(),
		},
		Salt:              "1234",
		Iterations:        0,
		AuthenticationMsg: "hello, world!",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Cproof != "WCRy25VQ5iCiPGQoLt4srzNvMhDDhVo19p72M8KB+cU=" {
		t.Fatalf("expected %s but received %s", "WCRy25VQ5iCiPGQoLt4srzNvMhDDhVo19p72M8KB+cU=", resp.Cproof)
	}
	if resp.Sproof != "/N4NMwnT+TeFI4Ymbaj0nk5sjJQTCrwnvaXhApkjRYo=" {
		t.Fatalf("expected %s but received %s", "/N4NMwnT+TeFI4Ymbaj0nk5sjJQTCrwnvaXhApkjRYo=", resp.Sproof)
	}

	// Now use a bad claimed arn and make sure we fail.
	resp, err = authenticator.GetPGSHA256Hash(nil, &pb.PGSHA256HashRequest{
		Authtype:       pb.AuthType_AWS,
		ClientLanguage: pb.ClientLanguage_GO,
		Dbhost:         "foo",
		Dbport:         "5432",
		Dbuser:         "bob",
		Awsauth: &pb.AWSAuth{
			SignedGetCallerIdentity: signedGetCallerIdentity,
			ClaimedIamArn:           "arn:partition:service:region:account-id:arn-thats-not-mine",
		},
		Salt:              "1234",
		Iterations:        0,
		AuthenticationMsg: "hello, world!",
	})
	if err == nil {
		t.Fatal("using a claimed arn that doesn't belong to me should fail")
	}
}

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
	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			verifiedARN, err := getAwsIdentity(testCase.SignedGetCallerIdentity, pb.ClientLanguage_GO)
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
			match, err := arnsMatch(testCase.ExpectedArn, verifiedARN)
			if err != nil {
				t.Fatal(err)
			}
			if !match {
				t.Fatalf("expected %s but received %s", testCase.ExpectedArn, verifiedARN)
			}
		})
	}
}

func TestToDatabaseARN(t *testing.T) {
	// Make sure role session names get stripped for assumed roles because
	// users won't be planting creds in databases with session names, since
	// they change all the time.
	result, err := toDatabaseARN("arn:aws:sts::account-id:assumed-role/role-name/role-session-name")
	if err != nil {
		t.Fatal(err)
	}
	if result != "arn:aws:iam::account-id:role/role-name" {
		t.Fatalf("expected %s but received %s", "arn:aws:iam::account-id:role/role-name", result)
	}

	// Leave other arns alone.
	result, err = toDatabaseARN("arn:aws:sts::123456789012:federated-user/my-federated-user-name")
	if err != nil {
		t.Fatal(err)
	}
	if result != "arn:aws:sts::123456789012:federated-user/my-federated-user-name" {
		t.Fatalf("expected %s but received %s", "arn:aws:sts::123456789012:federated-user/my-federated-user-name", result)
	}
}

func TestXorBytes(t *testing.T) {
	result := xorBytes([]byte{0}, []byte{0})
	expected := []byte{0}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %#v, but received %#v", expected, result)
	}

	result = xorBytes([]byte{1}, []byte{1})
	expected = []byte{0}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %#v, but received %#v", expected, result)
	}

	result = xorBytes([]byte{0, 1, 1}, []byte{0, 1, 1})
	expected = []byte{0, 0, 0}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %#v, but received %#v", expected, result)
	}

	result = xorBytes([]byte{1, 1, 1}, []byte{0, 0, 0})
	expected = []byte{1, 1, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %#v, but received %#v", expected, result)
	}
}

// TestNoRaces is intended to be run via "$ go test -v -race" and will
// fail if a race is detected.
func TestNoRaces(t *testing.T) {
	// These tests rely upon the file back-end, so unset the Vault addr if it exists.
	_ = os.Setenv(vault.EnvVaultAddress, "")

	// Get these in advance so they'll be only read. You can't have
	// races with objects that are only read.
	signedGetCallerIdentity, err := testEnv.SignedGetCallerIdentity(t)
	if err != nil {
		t.Fatal(err)
	}
	claimedARN := testEnv.ClaimedArn()

	// Create and start the authenticator as we normally would.
	authenticator, err := NewAuthenticator()
	if err != nil {
		t.Fatal(err)
	}
	go authenticator.run()

	// Try to create a race.
	start := make(chan interface{})
	done := make(chan interface{})
	for i := 0; i < 1000; i++ {
		go func() {
			<-start
			// We don't care about the response, we just want to hit as much
			// of the authenticator's code as possible.
			authenticator.GetPGSHA256Hash(nil, &pb.PGSHA256HashRequest{
				Authtype:       pb.AuthType_AWS,
				ClientLanguage: pb.ClientLanguage_GO,
				Dbhost:         "dbsha256",
				Dbport:         "5432",
				Dbuser:         "bob",
				Awsauth: &pb.AWSAuth{
					SignedGetCallerIdentity: signedGetCallerIdentity,
					ClaimedIamArn:           claimedARN,
				},
				Salt:              "1234",
				Iterations:        0,
				AuthenticationMsg: "hello, world!",
			})
			done <- true
		}()
		go func() {
			<-start
			// We don't care about the response, we just want to hit as much
			// of the authenticator's code as possible.
			authenticator.GetPGMD5Hash(nil, &pb.PGMD5HashRequest{
				Authtype:       pb.AuthType_AWS,
				ClientLanguage: pb.ClientLanguage_GO,
				Dbhost:         "dbmd5",
				Dbport:         "5432",
				Dbuser:         "bob",
				Awsauth: &pb.AWSAuth{
					SignedGetCallerIdentity: signedGetCallerIdentity,
					ClaimedIamArn:           testEnv.ClaimedArn(),
				},
				Salt: []byte{1, 2, 3, 4},
			})
			done <- true
		}()
	}
	// Close the start chan to fire all these calls at once.
	close(start)

	// Wait for them to finish, but in case they deadlock, time out.
	timer := time.NewTimer(10 * time.Second)
	for i := 0; i < 2000; i++ {
		select {
		case <-done:
		case <-timer.C:
			t.Fatal("over ten seconds elapsed")
		default:
		}
	}
}

// TestFuzzAuthenticator simply fuzzes its two request-receiving
// methods to ensure a panic isn't caused by random values. If
// a panic is produced, the test will fail.
func TestFuzzAuthenticator(t *testing.T) {
	// These tests rely upon the file back-end, so unset the Vault addr if it exists.
	_ = os.Setenv(vault.EnvVaultAddress, "")

	// This test generates a lot of error logs, so quiet them to
	// avoid them drowning out other tests.
	log.SetLevel(log.FatalLevel)

	authenticator, err := NewAuthenticator()
	if err != nil {
		t.Fatal(err)
	}

	fuzzer := fuzz.New()
	for i := 0; i < 1000; i++ {
		req1 := &pb.PGSHA256HashRequest{}
		fuzzer.Fuzz(req1)
		authenticator.GetPGSHA256Hash(nil, req1)

		req2 := &pb.PGMD5HashRequest{}
		fuzzer.Fuzz(req2)
		authenticator.GetPGMD5Hash(nil, req2)
	}
}

// This allows us to only get the signedGetCallerIdentity string once, but
// to reuse it throughout tests through the testEnv variable, reducing load
// on AWS.
type env struct {
	signedGetCallerIdentity string
}

func (e *env) ClaimedArn() string {
	return os.Getenv(envVarTestRole)
}

func (e *env) SignedGetCallerIdentity(t *testing.T) (string, error) {

	if os.Getenv(envVarTestRole) == "" {
		t.Fatalf("skipping because %s is unset", envVarTestRole)
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
