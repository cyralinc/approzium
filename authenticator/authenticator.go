package main

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/approzium/approzium/authenticator/credmgrs"
	pb "github.com/approzium/approzium/authenticator/protos"
	"github.com/aws/aws-sdk-go/aws/arn"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// validSTSEndpoints is presented as a variable so it
	// can be edited for testing if we need to mock the AWS
	// test server. This list is based off of
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_enable-regions.html
	validSTSEndpoints = []string{
		"sts.amazonaws.com",
		"sts.us-east-2.amazonaws.com",
		"sts.us-east-1.amazonaws.com",
		"sts.us-west-1.amazonaws.com",
		"sts.us-west-2.amazonaws.com",
		"sts.ap-east-1.amazonaws.com",
		"sts.ap-south-1.amazonaws.com",
		"sts.ap-northeast-2.amazonaws.com",
		"sts.ap-southeast-1.amazonaws.com",
		"sts.ap-southeast-2.amazonaws.com",
		"sts.ap-northeast-1.amazonaws.com",
		"sts.ca-central-1.amazonaws.com",
		"sts.eu-central-1.amazonaws.com",
		"sts.eu-west-1.amazonaws.com",
		"sts.eu-west-2.amazonaws.com",
		"eu-south-1",
		"sts.eu-west-3.amazonaws.com",
		"sts.eu-north-1.amazonaws.com",
		"sts.me-south-1.amazonaws.com",
		"af-south-1",
		"sts.sa-east-1.amazonaws.com",
	}

	// The initial choice for a max is based on
	// https://www.postgresql.org/docs/8.3/pgcrypto.html and
	// https://tools.ietf.org/html/rfc7677.
	// We want to allow enough iterations to be secure, but
	// not so many that the iterations could be used to effectively
	// DOS us by sending us looping for a long amount of time.
	// The RFC recommends at least 15,000 iterations, so we just
	// allow up to 10 times as much in case folks are being extra
	// secure. We are open to making this higher or lower based on
	// community feedback.
	maxIterations = uint32(15000 * 10)
)

type Authenticator struct {
	credMgr credmgrs.CredentialManager
	counter int
}

func NewAuthenticator() (*Authenticator, error) {
	credMgr, err := credmgrs.RetrieveConfigured()
	if err != nil {
		return nil, err
	}
	return &Authenticator{
		credMgr: credMgr,
	}, nil
}

func (a *Authenticator) run() {
	for {
		log.Printf("authenticator running. %d requests received", a.counter)
		time.Sleep(10 * time.Second)
	}
}

func executeGetCallerIdentity(request string) (string, error) {
	resp, err := http.Post(request, "", nil)
	if err != nil {
		return "", err
	}

	// We can ignore errors here because if it's not a valid response body, we'll still error
	// below when we try to unmarshal the XML.
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received unexpected get caller identity response %d: %s", resp.StatusCode, respBody)
	}

	type GetCallerIdentityResponse struct {
		IamArn string `xml:"GetCallerIdentityResult>Arn"`
	}
	response := GetCallerIdentityResponse{}
	err = xml.Unmarshal(respBody, &response)
	if err != nil {
		return "", err
	}
	return response.IamArn, nil
}

// verifyIdentity takes a signed get caller identity string and executes
// the request to the given AWS STS endpoint. It returns the caller's
// full IAM ARN.
func verifyIdentity(signedGetCallerIdentity string) (string, error) {
	u, err := url.Parse(signedGetCallerIdentity)
	if err != nil {
		return "", err
	}

	// Ensure the STS endpoint we'll be using is an AWS endpoint, and it's not
	// just some random server set up to mimic valid AWS STS responses.
	isValidSTSEndpoint := false
	for _, validSTSEndpoint := range validSTSEndpoints {
		if u.Host == validSTSEndpoint {
			isValidSTSEndpoint = true
			break
		}
	}
	if !isValidSTSEndpoint {
		return "", fmt.Errorf("%s is not a valid STS endpoint", u.Host)
	}

	// Ensure the call getting executed is actually the GetCallerIdentity call,
	// and not some other call that happens to return the expected XML fields.
	query := u.Query()
	if query.Get("Action") != "GetCallerIdentity" {
		return "", fmt.Errorf("invalid action for GetCallerIdentity: %s", query.Get("Action"))
	}
	return executeGetCallerIdentity(signedGetCallerIdentity)
}

func toDatabaseARN(fullIAMArn string) (string, error) {
	parsedArn, err := arn.Parse(fullIAMArn)
	if err != nil {
		return "", err
	}
	log.Debugf("received login attempt from %+v", parsedArn)
	if !strings.HasPrefix(parsedArn.Resource, "assumed-role") {
		// This is a regular ARN, so we should return it as-is for use in accessing
		// database credentials.
		return fullIAMArn, nil
	}
	// For assumed-role arns, they may have a session tag that we want to strip off
	// for accessing database credentials.
	fields := strings.Split(parsedArn.Resource, "/")
	if len(fields) != 3 {
		return "", fmt.Errorf("unexpected resource format for %s", parsedArn.Resource)
	}
	return fmt.Sprintf("arn:%s:%s::%s:%s/%s", parsedArn.Partition, parsedArn.Service, parsedArn.AccountID, fields[0], fields[1]), nil
}

func (a *Authenticator) GetPGMD5Hash(ctx context.Context, req *pb.PGMD5HashRequest) (*pb.PGMD5Response, error) {
	a.counter++

	// Return early if we didn't get a valid salt.
	salt := req.GetSalt()
	if len(salt) != 4 {
		msg := fmt.Sprintf("expected salt to be 4 bytes long, but got %d bytes", len(salt))
		log.Error(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}

	// To expedite handling the request, let's verify the caller's identity at the same
	// time as getting the password.
	verifiedIAMArnChan := make(chan string, 1)
	verificationErrChan := make(chan error, 1)
	go func() {
		verifiedIAMArn, err := verifyIdentity(req.GetSignedGetCallerIdentity())
		if err != nil {
			verificationErrChan <- status.Errorf(codes.Unauthenticated, err.Error())
			return
		}
		verifiedIAMArnChan <- verifiedIAMArn
	}()

	// Get the credentials.
	claimedIamArn := req.GetClaimedIamArn()
	dbHost := req.GetDbhost()
	dbPort := req.GetDbport()
	dbUser := req.GetDbuser()

	databaseArn, err := toDatabaseARN(claimedIamArn)
	if err != nil {
		return nil, err
	}
	password, err := a.getCreds(credmgrs.DBKey{
		IAMArn: databaseArn,
		DBHost: dbHost,
		DBPort: dbPort,
		DBUser: dbUser,
	})
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Make sure the ARN they claimed they had to get the creds was their actual ARN.
	select {
	case verifiedIAMArn := <-verifiedIAMArnChan:
		if verifiedIAMArn != claimedIamArn {
			return nil, status.Errorf(codes.Unauthenticated, fmt.Sprintf("claimed IAM ARN %s did not match actual IAM arn of %s", claimedIamArn, verifiedIAMArn))
		}
	case err = <-verificationErrChan:
		return nil, err
	}

	// Everything checked out.
	return &pb.PGMD5Response{Hash: computePGMD5Hash(dbUser, password, salt)}, nil
}

func (a *Authenticator) GetPGSHA256Hash(ctx context.Context, req *pb.PGSHA256HashRequest) (*pb.PGSHA256Response, error) {
	a.counter++

	// Return early if we didn't get a valid auth message or salt.
	authMsg := req.GetAuthenticationMsg()
	if len(authMsg) == 0 {
		msg := fmt.Sprintf("authentication message not provided")
		log.Error(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}

	salt := req.GetSalt()
	if len(salt) == 0 {
		msg := fmt.Sprintf("salt not provided")
		log.Error(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}

	iterations := req.GetIterations()
	if iterations > maxIterations {
		// Using a very high number of iterations could cause us to loop and a lot of
		// those requests could quickly take us down, like a DOS attack.
		msg := fmt.Sprintf("iterations too high, received %d but maximum is %d", iterations, maxIterations)
		log.Error(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}

	// To expedite handling the request, let's verify the caller's identity at the same
	// time as getting the password.
	verifiedIAMArnChan := make(chan string, 1)
	verificationErrChan := make(chan error, 1)
	go func() {
		verifiedIAMArn, err := verifyIdentity(req.GetSignedGetCallerIdentity())
		if err != nil {
			verificationErrChan <- status.Errorf(codes.Unauthenticated, err.Error())
			return
		}
		verifiedIAMArnChan <- verifiedIAMArn
	}()

	// Get the credentials.
	claimedIamArn := req.GetClaimedIamArn()
	dbHost := req.GetDbhost()
	dbPort := req.GetDbport()
	dbUser := req.GetDbuser()

	databaseArn, err := toDatabaseARN(claimedIamArn)
	if err != nil {
		return nil, err
	}

	password, err := a.getCreds(credmgrs.DBKey{
		IAMArn: databaseArn,
		DBHost: dbHost,
		DBPort: dbPort,
		DBUser: dbUser,
	})
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	saltedPass, err := computePGSHA256SaltedPass(password, salt, int(iterations))
	if err != nil {
		msg := fmt.Sprintf("Could not compute hash %s", err)
		log.Error(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}

	cproof := computePGSHA256Cproof(saltedPass, authMsg)
	sproof := computePGSHA256Sproof(saltedPass, authMsg)

	// Make sure the ARN they claimed they had to get the creds was their actual ARN.
	select {
	case verifiedIAMArn := <-verifiedIAMArnChan:
		if verifiedIAMArn != claimedIamArn {
			return nil, status.Errorf(codes.Unauthenticated, fmt.Sprintf("claimed IAM ARN %s did not match actual IAM arn of %s", claimedIamArn, verifiedIAMArn))
		}
	case err = <-verificationErrChan:
		return nil, err
	}
	return &pb.PGSHA256Response{Cproof: cproof, Sproof: sproof}, nil
}

func (a *Authenticator) getCreds(identity credmgrs.DBKey) (string, error) {
	creds, err := a.credMgr.Password(identity)
	if err != nil {
		msg := fmt.Errorf("password not found for identity %s", identity)
		log.Error(msg)
		return "", msg
	}
	return creds, nil
}

func computeMD5(s string, salt []byte) string {
	hasher := md5.New()
	io.WriteString(hasher, s)
	hasher.Write(salt)
	hashedBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}

func computePGMD5Hash(user, password string, salt []byte) string {
	firstHash := computeMD5(password, []byte(user))
	secondHash := computeMD5(firstHash, salt)
	return secondHash
}

func computePGSHA256SaltedPass(password string, salt string, iterations int) ([]byte, error) {
	s, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, fmt.Errorf("Bad salt %s", err)
	}
	dk := pbkdf2.Key([]byte(password), s, iterations, 32, sha256.New)
	return dk, nil
}

// assumes a and b are of the same length
func xorBytes(a, b []byte) []byte {
	buf := make([]byte, len(a))

	for i := range a {
		buf[i] = a[i] ^ b[i]
	}

	return buf
}

// SCRAM reference: https://en.wikipedia.org/wiki/Salted_Challenge_Response_Authentication_Mechanism
func computePGSHA256Cproof(spassword []byte, authMsg string) string {
	mac := hmac.New(sha256.New, spassword)
	mac.Write([]byte("Client Key"))
	ckey := mac.Sum(nil)
	ckeyHash := sha256.Sum256(ckey)
	cproofHmac := hmac.New(sha256.New, ckeyHash[:])
	cproofHmac.Write([]byte(authMsg))
	cproof := xorBytes(cproofHmac.Sum(nil), ckey)
	cproof64 := base64.StdEncoding.EncodeToString(cproof)
	return cproof64
}

func computePGSHA256Sproof(spassword []byte, authMsg string) string {
	mac := hmac.New(sha256.New, spassword)
	mac.Write([]byte("Server Key"))
	skey := mac.Sum(nil)
	sproofHmac := hmac.New(sha256.New, skey)
	sproofHmac.Write([]byte(authMsg))
	sproof := sproofHmac.Sum(nil)
	sproof64 := base64.StdEncoding.EncodeToString(sproof)
	return sproof64
}
