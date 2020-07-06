package server

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/approzium/approzium/authenticator/server/api"
	"github.com/approzium/approzium/authenticator/server/config"
	pb "github.com/approzium/approzium/authenticator/server/protos"
	testtools "github.com/approzium/approzium/authenticator/server/testing"
	"github.com/google/gofuzz"
	vault "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

var (
	testEnv      = &testtools.AwsEnv{}
	testLogEntry = func() *log.Entry {
		logEntry := log.WithFields(log.Fields{"test": "logger"})
		logEntry.Level = log.FatalLevel
		return logEntry
	}()
	testCtx = context.WithValue(context.Background(), ctxLogger, testLogEntry)
)

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

	authenticator, err := buildServer(testtools.TestLogger(), config.Config{})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := authenticator.GetPGMD5Hash(testCtx, &pb.PGMD5HashRequest{
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
	resp, err = authenticator.GetPGMD5Hash(testCtx, &pb.PGMD5HashRequest{
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

	authenticator, err := buildServer(testtools.TestLogger(), config.Config{})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := authenticator.GetPGSHA256Hash(testCtx, &pb.PGSHA256HashRequest{
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
	resp, err = authenticator.GetPGSHA256Hash(testCtx, &pb.PGSHA256HashRequest{
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

func TestToDatabaseARN(t *testing.T) {
	// Make sure role session names get stripped for assumed roles because
	// users won't be planting creds in databases with session names, since
	// they change all the time.
	result, err := toDatabaseARN(testLogEntry, "arn:aws:sts::account-id:assumed-role/role-name/role-session-name")
	if err != nil {
		t.Fatal(err)
	}
	if result != "arn:aws:iam::account-id:role/role-name" {
		t.Fatalf("expected %s but received %s", "arn:aws:iam::account-id:role/role-name", result)
	}

	// Leave other arns alone.
	result, err = toDatabaseARN(testLogEntry, "arn:aws:sts::123456789012:federated-user/my-federated-user-name")
	if err != nil {
		t.Fatal(err)
	}
	if result != "arn:aws:sts::123456789012:federated-user/my-federated-user-name" {
		t.Fatalf("expected %s but received %s", "arn:aws:sts::123456789012:federated-user/my-federated-user-name", result)
	}
}

func TestXorBytes(t *testing.T) {
	result, err := xorBytes([]byte{0}, []byte{0})
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte{0}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %#v, but received %#v", expected, result)
	}

	result, err = xorBytes([]byte{1}, []byte{1})
	if err != nil {
		t.Fatal(err)
	}
	expected = []byte{0}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %#v, but received %#v", expected, result)
	}

	result, err = xorBytes([]byte{0, 1, 1}, []byte{0, 1, 1})
	if err != nil {
		t.Fatal(err)
	}
	expected = []byte{0, 0, 0}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %#v, but received %#v", expected, result)
	}

	result, err = xorBytes([]byte{1, 1, 1}, []byte{0, 0, 0})
	if err != nil {
		t.Fatal(err)
	}
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
	authenticator, err := buildServer(testtools.TestLogger(), config.Config{})
	if err != nil {
		t.Fatal(err)
	}

	// Try to create a race.
	start := make(chan interface{})
	done := make(chan interface{})
	for i := 0; i < 1000; i++ {
		go func() {
			<-start
			// We don't care about the response, we just want to hit as much
			// of the authenticator's code as possible.
			authenticator.GetPGSHA256Hash(testCtx, &pb.PGSHA256HashRequest{
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
			authenticator.GetPGMD5Hash(testCtx, &pb.PGMD5HashRequest{
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

	authenticator, err := buildServer(testtools.TestLogger(), config.Config{})
	if err != nil {
		t.Fatal(err)
	}

	fuzzer := fuzz.New()
	for i := 0; i < 1000; i++ {
		req1 := &pb.PGSHA256HashRequest{}
		fuzzer.Fuzz(req1)
		authenticator.GetPGSHA256Hash(context.Background(), req1)

		req2 := &pb.PGMD5HashRequest{}
		fuzzer.Fuzz(req2)
		authenticator.GetPGMD5Hash(context.Background(), req2)
	}
}

func TestFuzzXorBytes(t *testing.T) {
	fuzzer := fuzz.New()
	for i := 0; i < 1000; i++ {

		var a []byte
		fuzzer.Fuzz(&a)

		var b []byte
		fuzzer.Fuzz(&b)

		xorBytes(a, b)
	}
}

func TestMetrics(t *testing.T) {
	// These tests rely upon the file back-end, so unset the Vault addr if it exists.
	_ = os.Setenv(vault.EnvVaultAddress, "")

	// Start the API, which includes an endpoint for Prometheus to mine metrics.
	config := config.Config{
		Host:     "127.0.0.1",
		HTTPPort: 6000,
	}
	_ = api.Start(testtools.TestLogger(), config)

	// Make some calls to increment the metrics.
	svr, err := buildServer(testtools.TestLogger(), config)
	if err != nil {
		t.Fatal(err)
	}
	svr.GetPGSHA256Hash(testCtx, &pb.PGSHA256HashRequest{})

	// See what we get for metrics.
	resp, err := http.Get("http://localhost:6000/v1/metrics/prometheus")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 but received %d", resp.StatusCode)
	}

	actualResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse, err := ioutil.ReadFile("testing/prometheus.txt")
	if err != nil {
		t.Fatal(err)
	}

	if clean(expectedResponse) != clean(actualResponse) {
		t.Fatalf("expected %s but received %s", expectedResponse, actualResponse)
	}
}

func clean(b []byte) string {
	cleaned := strings.ReplaceAll(string(b), " ", "")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	return cleaned
}
