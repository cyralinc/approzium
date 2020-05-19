package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"

	pb "dbauth/authenticator/messages"
)

type Authenticator struct {
	creds map[string]credentials
}

type credentials struct {
	user     string
	password string
}

func NewAuthenticator() *Authenticator {
	return &Authenticator{
		creds: newCreds(),
	}
}

func (a *Authenticator) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	identity := req.GetIdentity()
	salt := req.GetSalt()
	fmt.Printf("received request to return hashed credentials for identity %s given salt %s",
		identity, salt)

	if salt == "" {
		msg := fmt.Errorf("salt not received")
		return &pb.AuthenticateResponse{
			Status:  pb.AuthenticateResponse_ERROR,
			Message: fmt.Sprintf("%s", msg),
		}, msg
	}

	creds, err := a.getCreds(identity)
	if err != nil {
		return &pb.AuthenticateResponse{
			Status:  pb.AuthenticateResponse_ERROR,
			Message: fmt.Sprintf("%s", err),
		}, err
	}

	hashedCreds := pb.Credentials{
		User:           creds.user,
		HashedPassword: computeMD5(creds.password, salt),
	}

	return &pb.AuthenticateResponse{
		Status:      pb.AuthenticateResponse_SUCCESS,
		Credentials: &hashedCreds,
	}, nil
}

func (a *Authenticator) getCreds(identity string) (credentials, error) {
	if creds, ok := a.creds[identity]; ok {
		return creds, nil
	}
	msg := fmt.Errorf("credentials not found for identity %s", identity)
	return credentials{}, msg
}

func newCreds() map[string]credentials {
	creds := make(map[string]credentials)
	creds["diotim"] = credentials{
		user:     "bob",
		password: "password",
	}
	return creds
}

func computeMD5(s, salt string) string {
	hasher := md5.New()
	io.WriteString(hasher, s)
	io.WriteString(hasher, salt)
	hashedBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}
