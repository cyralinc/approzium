package identity

import (
	"errors"
	"fmt"

	pb "github.com/approzium/approzium/authenticator/server/protos"
	log "github.com/sirupsen/logrus"
)

type Proof struct {
	AwsAuth *pb.AWSIdentity
}

type Verified struct {
	IamArn string
}

func (v *Verified) Matches(claimedIdentity interface{}) (bool, error) {
	if v.IamArn != "" {
		claimedArn, ok := claimedIdentity.(string)
		if !ok {
			return false, fmt.Errorf("expected claimed ARN but received %s, a %T", claimedIdentity, claimedIdentity)
		}
		return arnsMatch(claimedArn, v.IamArn)
	}
	return false, errors.New("no verified identity detected")
}

// Get is for getting a verified identity using its proof of identity.
func Get(logger *log.Entry, proof *Proof, clientLang pb.ClientLanguage) (*Verified, error) {
	// Log here, and pass verification down.
	verified, err := get(proof, clientLang)
	if err != nil {
		logger.Warn(fmt.Sprintf("couldn't verify %s: %s", proof, err))
	} else {
		logger.Info(fmt.Sprintf("verified %s", verified))
	}
	return verified, err
}

// get is the main method that actually verifies the proof it's given for an identity.
func get(proof *Proof, clientLang pb.ClientLanguage) (*Verified, error) {
	if proof.AwsAuth == nil {
		return nil, fmt.Errorf("AWS auth info is required")
	}
	iamArn, err := getAwsIdentity(proof.AwsAuth.SignedGetCallerIdentity, clientLang)
	if err != nil {
		return nil, err
	}
	return &Verified{
		IamArn: iamArn,
	}, nil
}
