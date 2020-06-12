package main

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	pb "approzium/authenticator/protos"
	log "github.com/sirupsen/logrus"
)

type vaultKey struct {
	iam_arn string
	dbhost  string
	dbuser  string
}

type Authenticator struct {
	vault   map[vaultKey]string
	counter int
}

func NewAuthenticator() *Authenticator {
	return &Authenticator{
		vault: newVault(),
	}
}

func (a *Authenticator) run() {
	for {
		log.Printf("authenticator running. %d requests received\n", a.counter)
		time.Sleep(10 * time.Second)
	}
}

func IAMFormatToSTS(iam_arn string) (string, error) {
	matches := regexp.MustCompile(`arn:aws:iam::(.*):role/(.*)`).FindStringSubmatch(iam_arn)
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

func verifyService(claimed_iam_arn, signed_get_caller_identity string) error {
	log.Printf("verifying service for role: %s\n", claimed_iam_arn)
	actual_iam_arn, err := executeGetCallerIdentity(signed_get_caller_identity)
	if err != nil {
		return fmt.Errorf("could not execute GetCallerIdentity %s", err)
	}
	// have to change formats of arns to be able to do string comparison
	claimed_iam_arn, err = IAMFormatToSTS(claimed_iam_arn)
	if err != nil {
		return fmt.Errorf("could not parse claimed IAM ARN %s", err)
	}

	// uses prefix check because user might have added a session tag in their claimed ARN
	// for example, the following two IAMs should match
	// arn:aws:sts::403019568400:assumed-role/dev
	// arn:aws:sts::403019568400:assumed-role/dev/Service1
	if strings.HasPrefix(actual_iam_arn, claimed_iam_arn) {
		return nil
	} else {
		return fmt.Errorf("actual IAM ARN %s does not match claimed IAM ARN %s", actual_iam_arn, claimed_iam_arn)
	}
}

func (a *Authenticator) GetPGMD5Hash(ctx context.Context, req *pb.PGMD5HashRequest) (*pb.PGMD5Response, error) {
	a.counter++

	claimed_iam_arn := req.GetClaimedIamArn()
	dbhost := req.GetDbhost()
	dbuser := req.GetDbuser()
	log.Printf("received GetPGMD5Hash request with claimed_iam_arn: %s\n", claimed_iam_arn)
	err := verifyService(claimed_iam_arn, req.GetSignedGetCallerIdentity())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	password, err := a.getCreds(vaultKey{claimed_iam_arn, dbhost, dbuser})
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

	return &pb.PGMD5Response{Hash: computePGMD5(dbuser, password, salt)}, nil
}

func (a *Authenticator) GetPGSHA256Hash(ctx context.Context, req *pb.PGSHA256HashRequest) (*pb.PGSHA256Response, error) {
	a.counter++

	claimed_iam_arn := req.GetClaimedIamArn()
	dbhost := req.GetDbhost()
	dbuser := req.GetDbuser()
	log.Printf("received GetPGSHA256Hash request with claimed_iam_arn: %s\n", claimed_iam_arn)
	err := verifyService(claimed_iam_arn, req.GetSignedGetCallerIdentity())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	password, err := a.getCreds(vaultKey{claimed_iam_arn, dbhost, dbuser})
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
	spassword, err := computePGSHA256(password, salt, iterations)
	if err != nil {
		msg := fmt.Sprintf("Could not compute hash %s", err)
		log.Error(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}
	return &pb.PGSHA256Response{Spassword: spassword}, nil
}

func (a *Authenticator) getCreds(identity vaultKey) (string, error) {
	if creds, ok := a.vault[identity]; ok {
		return creds, nil
	}
	msg := fmt.Errorf("password not found for identity %s", identity)
	log.Error(msg)
	return "", msg
}

func newVault() map[vaultKey]string {
	creds := make(map[vaultKey]string)
	// for dev purposes: read credentials from a local file
	type secrets []struct {
		Dbhost   string `yaml:"dbhost"`
		Dbuser   string `yaml:"dbuser"`
		Password string `yaml:"password"`
	}
	var devCreds secrets
	yamlFile, err := ioutil.ReadFile("secrets.yaml")
	if err != nil {
		log.Errorf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &devCreds)
	if err != nil {
		log.Errorf("Unmarshal: #%v ", err)
	}
	for _, cred := range devCreds {
		key := vaultKey{"arn:aws:iam::403019568400:role/dev", cred.Dbhost, cred.Dbuser}
		creds[key] = cred.Password
		log.Debugf("added dev credential for host %s", cred.Dbhost)
	}
	return creds
}

func computeMD5(s string, salt []byte) string {
	hasher := md5.New()
	io.WriteString(hasher, s)
	hasher.Write(salt)
	hashedBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}

func computePGMD5(user, password string, salt []byte) string {
	first_hash := computeMD5(password, []byte(user))
	second_hash := computeMD5(first_hash, salt)
	return second_hash
}

func computePGSHA256(password string, salt []byte, iterations int) ([]byte, error) {
	s, err := base64.StdEncoding.DecodeString(string(salt))
	if err != nil {
		return nil, fmt.Errorf("Bad salt %s", err)
	}
	dk := pbkdf2.Key([]byte(password), s, iterations, 32, sha256.New)
	return dk, nil
}
