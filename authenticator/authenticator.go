package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	pb "dbauth/authenticator/messages"
	log "github.com/sirupsen/logrus"
)

type Authenticator struct {
	creds   map[string]credentials
	counter int
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

func (a *Authenticator) run() {
	for {
		log.Printf("authenticator running. %d requests received\n", a.counter)
		time.Sleep(2 * time.Second)
	}
}

func (a *Authenticator) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	a.counter++
	identity := req.GetIdentity()
	salt := req.GetSalt()
	log.Printf("received request to return hashed credentials for identity %s given salt %s\n",
		identity, salt)

	if len(salt) == 0 {
		msg := fmt.Errorf("salt not received")
		log.Error(msg)
		return &pb.AuthenticateResponse{
			Status:  pb.AuthenticateResponse_ERROR,
			Message: fmt.Sprintf("%s", msg),
		}, msg
	}

	creds, err := a.getCreds(identity)
	if err != nil {
		log.Error(err)
		return &pb.AuthenticateResponse{
			Status:  pb.AuthenticateResponse_ERROR,
			Message: fmt.Sprintf("%s", err),
		}, err
	}

	hashedCreds := pb.Credentials{
		User:           creds.user,
		HashedPassword: computePGMD5(identity, creds.password, salt),
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
	log.Error(msg)
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

func computeMD5(s string, salt []byte) string {
	hasher := md5.New()
	io.WriteString(hasher, s)
	hasher.Write(salt[:4])
	hashedBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}

func computePGMD5(user, password string, salt []byte) string {
	first_hash := computeMD5(password, []byte(user))
	second_hash := computeMD5(first_hash, salt)
	return "md5" + second_hash
}
