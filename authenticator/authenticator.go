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
	"regexp"
	"strings"
	"time"

	"github.com/approzium/approzium/authenticator/credmgrs"
	pb "github.com/approzium/approzium/authenticator/protos"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		log.Printf("authenticator running. %d requests received\n", a.counter)
		time.Sleep(10 * time.Second)
	}
}

func IAMFormatToSTS(iamArn string) (string, error) {
	matches := regexp.MustCompile(`arn:aws:iam::(.*):role/(.*)`).FindStringSubmatch(iamArn)
	if matches == nil {
		return "", fmt.Errorf("provided IAM role ARN is not properly formatted, expected format: arn:aws:iam::accountid:role/rolename")
	}
	accountId := matches[1]
	role := matches[2]
	return fmt.Sprintf("arn:aws:sts::%s:assumed-role/%s", accountId, role), nil
}

func executeGetCallerIdentity(request string) (string, error) {
	resp, err := http.Post(request, "", nil)
	if err != nil {
		return "", err
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	type GetCallerIdentityResponse struct {
		IamArn string `xml:"GetCallerIdentityResult>Arn"`
	}
	response := GetCallerIdentityResponse{}
	err = xml.Unmarshal(responseData, &response)
	if err != nil {
		return "", err
	}
	iamArn := strings.Trim(response.IamArn, "{}")
	return iamArn, nil
}

func verifyService(claimedIamArn, signedGetCallerIdentity string) error {
	log.Printf("verifying service for role: %s\n", claimedIamArn)
	actualIamArn, err := executeGetCallerIdentity(signedGetCallerIdentity)
	if err != nil {
		return fmt.Errorf("could not execute GetCallerIdentity %s", err)
	}
	// have to change formats of arns to be able to do string comparison
	claimedIamArn, err = IAMFormatToSTS(claimedIamArn)
	if err != nil {
		return fmt.Errorf("could not parse claimed IAM ARN %s", err)
	}

	// uses prefix check because user might have added a session tag in their claimed ARN
	// for example, the following two IAMs should match
	// arn:aws:sts::403019568400:assumed-role/dev
	// arn:aws:sts::403019568400:assumed-role/dev/Service1
	if strings.HasPrefix(actualIamArn, claimedIamArn) {
		return nil
	} else {
		return fmt.Errorf("actual IAM ARN %s does not match claimed IAM ARN %s", actualIamArn, claimedIamArn)
	}
}

func (a *Authenticator) GetPGMD5Hash(ctx context.Context, req *pb.PGMD5HashRequest) (*pb.PGMD5Response, error) {
	a.counter++

	claimedIamArn := req.GetClaimedIamArn()
	dbHost := req.GetDbhost()
	dbPort := req.GetDbport()
	dbUser := req.GetDbuser()
	log.Printf("received GetPGMD5Hash request with claimedIamArn: %s\n", claimedIamArn)
	err := verifyService(claimedIamArn, req.GetSignedGetCallerIdentity())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	password, err := a.getCreds(credmgrs.DBKey{
		IAMArn: claimedIamArn,
		DBHost: dbHost,
		DBPort: dbPort,
		DBUser: dbUser,
	})
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	salt := req.GetSalt()

	if len(salt) != 4 {
		msg := fmt.Sprintf("expected salt to be 4 bytes long, but got %d bytes", len(salt))
		log.Error(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}

	return &pb.PGMD5Response{Hash: computePGMD5Hash(dbUser, password, salt)}, nil
}

func (a *Authenticator) GetPGSHA256Hash(ctx context.Context, req *pb.PGSHA256HashRequest) (*pb.PGSHA256Response, error) {
	a.counter++

	claimedIamArn := req.GetClaimedIamArn()
	dbHost := req.GetDbhost()
	dbPort := req.GetDbport()
	dbUser := req.GetDbuser()
	log.Printf("received GetPGSHA256Hash request with claimedIamArn: %s\n", claimedIamArn)
	err := verifyService(claimedIamArn, req.GetSignedGetCallerIdentity())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	password, err := a.getCreds(credmgrs.DBKey{
		IAMArn: claimedIamArn,
		DBHost: dbHost,
		DBPort: dbPort,
		DBUser: dbUser,
	})
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	salt := req.GetSalt()

	if len(salt) == 0 {
		msg := fmt.Sprintf("salt not provided")
		log.Error(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}
	iterations := int(req.GetIterations())
	saltedPass, err := computePGSHA256SaltedPass(password, salt, iterations)
	if err != nil {
		msg := fmt.Sprintf("Could not compute hash %s", err)
		log.Error(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}

	authMsg := req.GetAuthenticationMsg()

	if len(authMsg) == 0 {
		msg := fmt.Sprintf("authentication message not provided")
		log.Error(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}
	cproof := computePGSHA256Cproof(saltedPass, authMsg)
	sproof := computePGSHA256Sproof(saltedPass, authMsg)
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
