package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	pb "dbauth/authenticator/protos"
	log "github.com/sirupsen/logrus"
)

type Authenticator struct {
	vault   map[string]credentials
	counter int
}

type credentials struct {
	user     string
	password string
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

func verifyService(claimed_iam_arn, signed_get_caller_identity string) error {
	resp, err := http.Post(signed_get_caller_identity, "", nil)
	log.Printf("verifying service for role: %s\n", claimed_iam_arn)
	if err != nil {
		return err
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	responseString := string(responseData)
	re := regexp.MustCompile(`<Arn>(.*)</Arn>`)
	match := re.FindStringSubmatch(responseString)
	if match == nil {
		return fmt.Errorf("no ARN found in response")
	}
	returned_arn := match[1]
	// check if returned and provided ARNs match
	matches := regexp.MustCompile(`arn:aws:iam::(.*):role/(.*)`).FindStringSubmatch(claimed_iam_arn)
	if matches == nil {
		return fmt.Errorf("provided IAM role ARN not properly formatted")
	}
	accountId := matches[1]
	role := matches[2]
	expected := fmt.Sprintf("arn:aws:sts::%s:assumed-role/%s", accountId, role)
	// uses prefix check because user might have added a session tag in their claimed ARN
	// example:
	// arn:aws:sts::403019568400:assumed-role/dev
	// arn:aws:sts::403019568400:assumed-role/dev/Service1
	if strings.HasPrefix(returned_arn, expected) {
		return nil
	} else {
		return fmt.Errorf("received ARN %s does not match claimed ARN %s", match[1], claimed_iam_arn)
	}
}

func (a *Authenticator) GetDBUser(ctx context.Context, req *pb.DBUserRequest) (*pb.DBUserResponse, error) {
	a.counter++

	claimed_iam_arn := req.GetIamRoleArn()
	log.Printf("received GetDBUser request\n")
	err := verifyService(claimed_iam_arn, req.GetSignedGetCallerIdentity())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	creds, err := a.getCreds(claimed_iam_arn + "/" + req.GetDbname())
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &pb.DBUserResponse{Dbuser: creds.user}, nil
}

func (a *Authenticator) GetDBHash(ctx context.Context, req *pb.DBHashRequest) (*pb.DBHashResponse, error) {
	a.counter++

	claimed_iam_arn := req.GetIamRoleArn()
	log.Printf("received GetDBHash request\n")
	err := verifyService(claimed_iam_arn, req.GetSignedGetCallerIdentity())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	creds, err := a.getCreds(claimed_iam_arn + "/" + req.GetDbname())
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

	return &pb.DBHashResponse{Hash: computePGMD5(creds.user, creds.password, salt)}, nil
}

func (a *Authenticator) getCreds(identity string) (credentials, error) {
	if creds, ok := a.vault[identity]; ok {
		return creds, nil
	}
	msg := fmt.Errorf("credentials not found for identity %s", identity)
	log.Error(msg)
	return credentials{}, msg
}

func newVault() map[string]credentials {
	creds := make(map[string]credentials)
	creds["diotim"] = credentials{
		user:     "bob",
		password: "password",
	}
	creds["arn:aws:iam::403019568400:role/dev/db"] = credentials{
		user:     "bob",
		password: "password",
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
	return "md5" + second_hash
}
