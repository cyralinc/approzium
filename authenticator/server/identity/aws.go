package identity

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	pb "github.com/cyralinc/approzium/authenticator/server/protos"
	log "github.com/sirupsen/logrus"
)

// validSTSEndpoints is presented as a variable so it
// can be edited for testing if we need to mock the AWS
// test server. This list is based off of
// https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_enable-regions.html
var validSTSEndpoints = []string{
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

type aws struct{}

func (a *aws) Get(_ *log.Entry, proof *Proof) (*Verified, error) {
	if proof.AwsAuth == nil {
		return nil, fmt.Errorf("AWS auth info is required")
	}
	iamArn, err := a.getAwsIdentity(proof.AwsAuth.SignedGetCallerIdentity, proof.ClientLang)
	if err != nil {
		return nil, err
	}
	return &Verified{
		authType: authTypeAws,
		iamArn:   iamArn,
	}, nil
}

func (a *aws) Matches(_ *log.Entry, claimedIdentity string, verifiedIdentity *Verified) (bool, error) {
	return a.arnsMatch(claimedIdentity, verifiedIdentity.iamArn)
}

// getAwsIdentity takes a signed get caller identity string and executes
// the request to the given AWS STS endpoint. It returns the caller's
// full IAM arn.
func (a *aws) getAwsIdentity(signedGetCallerIdentity string, clientLanguage pb.ClientLanguage) (string, error) {
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
	return a.executeGetCallerIdentity(signedGetCallerIdentity, clientLanguage)
}

func (a *aws) executeGetCallerIdentity(signedGetCallerIdentity string, clientLanguage pb.ClientLanguage) (string, error) {
	var resp *http.Response
	var err error
	switch clientLanguage {
	case pb.ClientLanguage_GO:
		resp, err = http.Get(signedGetCallerIdentity) // #nosec: URL is verified before execution
	case pb.ClientLanguage_PYTHON:
        resp, err = http.Post(signedGetCallerIdentity, "", nil) // #nosec: URL is verified before execution
	case pb.ClientLanguage_LANGUAGE_NOT_PROVIDED:
		return "", fmt.Errorf("client language must be provided for AWS authentication")
	default:
		return "", fmt.Errorf("unsupported SDK type of %d", clientLanguage)
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received unexpected get caller identity response %d: %s", resp.StatusCode, respBody)
	}

	type GetCallerIdentityResponse struct {
		IamArn string `xml:"GetCallerIdentityResult>Arn"`
	}
	response := GetCallerIdentityResponse{}
	if err = xml.Unmarshal(respBody, &response); err != nil {
		return "", err
	}
	return response.IamArn, nil
}

// arnsMatch compares a claimed arn that the caller states they'll
// have, and an actual arn returned by the AWS get caller identity call.
func (a *aws) arnsMatch(claimedArn, actualArn string) (bool, error) {
	if claimedArn == "" || actualArn == "" {
		return false, nil
	}
	if claimedArn == actualArn {
		return true, nil
	}

	// If they're not immediately equal, check for special situations
	// where we would still allow a match.
	// We allow role arn to match the assumed role arn folks would have
	// for that role.
	type WrappedARN struct {
		arn.ARN
		RoleName string
	}

	var assumedRole *WrappedARN
	var role *WrappedARN
	for _, rawArn := range []string{claimedArn, actualArn} {
		parsed, err := arn.Parse(rawArn)
		if err != nil {
			return false, err
		}
		wrappedARN := &WrappedARN{
			ARN: parsed,
		}
		if strings.HasPrefix(wrappedARN.Resource, "assumed-role/") {
			fields := strings.Split(wrappedARN.Resource, "/")
			// Assumed role resource strings look like "assumed-role/rolename/session",
			// but they may not have session on the end.
			if len(fields) < 2 || len(fields) > 3 {
				return false, fmt.Errorf("received assumed role arn that doesn't match the expected format: %s", rawArn)
			}
			wrappedARN.RoleName = fields[1]
			assumedRole = wrappedARN
			continue
		}
		if strings.HasPrefix(wrappedARN.Resource, "role/") {
			fields := strings.Split(wrappedARN.Resource, "/")
			// Role resource strings look like "role/rolename".
			if len(fields) != 2 {
				return false, fmt.Errorf("received role arn that doesn't match the expected format: %s", rawArn)
			}
			wrappedARN.RoleName = fields[1]
			role = wrappedARN
			continue
		}
	}
	if assumedRole == nil || role == nil {
		// Since we only special case matching role arns with assumed role arns,
		// we can conclude that these don't match.
		return false, nil
	}

	// Compare the role arn and the assumed role arn to ensure they match.
	if role.Partition != assumedRole.Partition {
		return false, fmt.Errorf("partitions don't match, claimed arn %s, actual arn %s", claimedArn, actualArn)
	}
	if role.Service != "iam" {
		return false, fmt.Errorf("received unexpected service for role: %s", role.String())
	}
	if assumedRole.Service != "sts" {
		return false, fmt.Errorf("received unexpected service for assumed role: %s", assumedRole.String())
	}
	if role.Region != assumedRole.Region {
		return false, fmt.Errorf("regions don't match, claimed arn %s, actual arn %s", claimedArn, actualArn)
	}
	if role.AccountID != assumedRole.AccountID {
		return false, fmt.Errorf("account IDs don't match, claimed arn %s, actual arn %s", claimedArn, actualArn)
	}
	if role.RoleName != assumedRole.RoleName {
		return false, fmt.Errorf("role names don't match, claimed arn %s, actual arn %s", claimedArn, actualArn)
	}
	return true, nil
}
