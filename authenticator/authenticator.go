package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
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
		time.Sleep(2 * time.Second)
	}
}

func verifyService(claimed_iam_arn, signed_get_caller_identity string) error {
	log.Printf("verifying service for role: %s\n", claimed_iam_arn)
	return nil
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
