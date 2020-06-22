package main

import (
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	pb "github.com/approzium/approzium/authenticator/protos"
)

const exampleSTSReq = "https://sts.amazonaws.com/?" +
	"Action=GetCallerIdentity&" +
	"Version=2011-06-15&" +
	"X-Amz-Algorithm=AWS4-HMAC-SHA256&" +
	"X-Amz-Credential=ASIA2VNACNFCQTKAMK7P%2F20200619%2Fus-east-1%2Fsts%2Faws4_request&" +
	"X-Amz-Date=20200619T213034Z&" +
	"X-Amz-Expires=3600&" +
	"X-Amz-SignedHeaders=host&" +
	"X-Amz-Security-Token=FwoGZXIvYXdzEKf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaDI3Fc5sH4PVjCltv%2BCKsAZJbeqGXFK4iRoGzVACyFb1lirj1pCg278WgTVOEwA9cbSaSz%2FbLSXjsWAwbQGTo8KzcqLuvHV9IALku5ncJz4XXJ2WQtN7S5qpv%2BJ%2BK0U7hA%2Bk0ktRlhpoUWbJJeV7RCkMF5xSkOSQ3T4RB0PVH3kALjAcEhrEXwMAHD%2FU7RUXQaQkJYNA%2Ba7InAU8%2BorE4Ksuw4YKZGPYaHaEEnDIq7x0phmY0PFcqqor9mSMo%2Bty09wUyLX4CJ2MtA4fowq5hUuIluUPoVNM1Tk9YWY31VcWeakteFoSOHNqoo4kvwcNMVQ%3D%3D&" +
	"X-Amz-Signature=26cccaf1eb751690921676ce3bb8272f3dd3119d7c93995aa468da356814a17"

var returnUnauthorizedArn = false

func TestAuthenticator_GetPGMD5Hash(t *testing.T) {
	mockAwsServer := setup(t)
	defer mockAwsServer.Close()

	validTestRequest := strings.ReplaceAll(exampleSTSReq, "https://sts.amazonaws.com", mockAwsServer.URL)

	authenticator, err := NewAuthenticator()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := authenticator.GetPGMD5Hash(nil, &pb.PGMD5HashRequest{
		SignedGetCallerIdentity: validTestRequest,
		ClaimedIamArn:           "arn:aws:iam::403019568400:assumed-role/dev",
		Dbhost:                  "dbmd5",
		Dbport:                  "5432",
		Dbuser:                  "bob",
		Salt:                    []byte{1, 2, 3, 4},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Hash != "d576ce99165615bb3f4331154ca6660c" {
		t.Fatalf("expected %s but received %s", "d576ce99165615bb3f4331154ca6660c", resp.Hash)
	}

	// Now use a bad claimed ARN and make sure we fail.
	returnUnauthorizedArn = true
	defer func() {
		returnUnauthorizedArn = false
	}()
	resp, err = authenticator.GetPGMD5Hash(nil, &pb.PGMD5HashRequest{
		SignedGetCallerIdentity: validTestRequest,
		ClaimedIamArn:           "arn:aws:sts::123456789012:federated-user/my-federated-user-name",
		Dbhost:                  "dbmd5",
		Dbport:                  "5432",
		Dbuser:                  "bob",
		Salt:                    []byte{1, 2, 3, 4},
	})
	if err == nil {
		t.Fatal("using a claimed ARN that doesn't belong to me should fail")
	}
}

func TestAuthenticator_GetPGSHA256Hash(t *testing.T) {
	mockAwsServer := setup(t)
	defer mockAwsServer.Close()

	validTestRequest := strings.ReplaceAll(exampleSTSReq, "https://sts.amazonaws.com", mockAwsServer.URL)

	authenticator, err := NewAuthenticator()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := authenticator.GetPGSHA256Hash(nil, &pb.PGSHA256HashRequest{
		SignedGetCallerIdentity: validTestRequest,
		ClaimedIamArn:           "arn:aws:iam::403019568400:assumed-role/dev",
		Dbhost:                  "dbsha256",
		Dbport:                  "5432",
		Dbuser:                  "bob",
		Salt:                    "1234",
		Iterations:              0,
		AuthenticationMsg:       "hello, world!",
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

	// Now use a bad claimed ARN and make sure we fail.
	returnUnauthorizedArn = true
	defer func() {
		returnUnauthorizedArn = false
	}()
	resp, err = authenticator.GetPGSHA256Hash(nil, &pb.PGSHA256HashRequest{
		SignedGetCallerIdentity: validTestRequest,
		ClaimedIamArn:           "arn:aws:sts::123456789012:federated-user/my-federated-user-name",
		Dbhost:                  "dbsha256",
		Dbport:                  "5432",
		Dbuser:                  "bob",
		Salt:                    "1234",
		Iterations:              0,
		AuthenticationMsg:       "hello, world!",
	})
	if err == nil {
		t.Fatal("using a claimed ARN that doesn't belong to me should fail")
	}
}

func TestVerifyService(t *testing.T) {
	// Create a mock test server where we'll receive AWS calls.
	mockAwsServer := setup(t)
	defer mockAwsServer.Close()

	validTestRequest := strings.ReplaceAll(exampleSTSReq, "https://sts.amazonaws.com", mockAwsServer.URL)
	testCases := []struct {
		TestName                string
		SignedGetCallerIdentity string
		ExpectedArn             string
		ExpectErr               bool
	}{
		{
			TestName:                "Sunny path, regular ARN",
			SignedGetCallerIdentity: validTestRequest,
			ExpectedArn:             "arn:aws:iam::403019568400:assumed-role/dev",
			ExpectErr:               false,
		},
		{
			TestName:                "Empty values",
			SignedGetCallerIdentity: "",
			ExpectErr:               true,
		},
		{
			TestName:                "Malicious base URL injected",
			SignedGetCallerIdentity: strings.ReplaceAll(validTestRequest, mockAwsServer.URL, "127.0.0.1"),
			ExpectErr:               true,
		},
		{
			TestName:                "Different call than GetCallerIdentity",
			SignedGetCallerIdentity: strings.ReplaceAll(validTestRequest, "GetCallerIdentity", "GetSessionToken"),
			ExpectErr:               true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			verifiedARN, err := verifyIdentity(testCase.SignedGetCallerIdentity)
			if testCase.ExpectErr {
				if err == nil {
					t.Fatal("expected err")
				} else {
					// We expected an error and received it, so we've succeeded
					// and there's nothing else to do here.
					return
				}
			}
			// We don't expect an error. Let's make sure we got the expected response.
			if verifiedARN != testCase.ExpectedArn {
				t.Fatalf("expected %s but received %s", testCase.ExpectedArn, verifiedARN)
			}
			if err != nil {
				t.Fatal(err)
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
	if result != "arn:aws:sts::account-id:assumed-role/role-name" {
		t.Fatalf("expected %s but received %s", "arn:aws:sts::account-id:assumed-role/role-name/role-session-name", result)
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

func setup(t *testing.T) *httptest.Server {
	// These tests rely upon the file back-end, so unset the Vault addr if it exists.
	os.Setenv(vault.EnvVaultAddress, "")

	// Now mock the AWS server.
	defaultArn := "arn:aws:iam::403019568400:assumed-role/dev"
	unauthorizedArn := "arn:aws:sts::123456789012:federated-user/malicious-user"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		responseTemplate := `<GetCallerIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetCallerIdentityResult>
    <Arn>%s</Arn>
    <UserId>123456789012:my-federated-user-name</UserId>
    <Account>123456789012</Account>
  </GetCallerIdentityResult>
  <ResponseMetadata>
    <RequestId>01234567-89ab-cdef-0123-456789abcdef</RequestId>
  </ResponseMetadata>
</GetCallerIdentityResponse>`
		if returnUnauthorizedArn {
			w.Write([]byte(fmt.Sprintf(responseTemplate, unauthorizedArn)))
			return
		} else {
			w.Write([]byte(fmt.Sprintf(responseTemplate, defaultArn)))
			return
		}
	}))

	tsUrl, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Overwrite the valid baseURLs to fire requests at our test server.
	validSTSEndpoints = []string{
		tsUrl.Host,
	}
	return ts
}
