package identity

import (
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	pb "github.com/cyralinc/approzium/authenticator/server/protos"
	log "github.com/sirupsen/logrus"
)

// GetCallerIdentity strings expire every 15 minutes to prevent
// replay attacks, so let's refresh ours every 5.
var refreshEvery = 5 * time.Minute

// To not assume a role, simply provide "".
func newAwsIdentityHandler(logger *log.Logger, roleArnToAssume string) (*awsIdentityHandler, error) {
	h := &awsIdentityHandler{
		logger:          logger,
		roleArnToAssume: roleArnToAssume,
	}
	// Initially ensure we can get a caller identity on the main thread so
	// we can provide immediate feedback if it's misconfigured.
	if err := h.refreshProof(); err != nil {
		return nil, err
	}
	// Now refresh identity in the background so it's super fast to retrieve
	// the latest identity.
	h.startBackgroundRefresh()
	return h, nil
}

type awsIdentityHandler struct {
	logger          *log.Logger
	roleArnToAssume string

	identityLock   sync.RWMutex
	latestIdentity *pb.AWSIdentity
}

func (h *awsIdentityHandler) RetrieveAWSIdentity() *pb.AWSIdentity {
	h.identityLock.RLock()
	defer h.identityLock.RUnlock()
	return h.latestIdentity
}

func (h *awsIdentityHandler) startBackgroundRefresh() {
	go func() {
		ticker := time.NewTicker(refreshEvery)
		for true {
			select {
			case <-ticker.C:
				if err := h.refreshProof(); err != nil {
					// In case the issue was transient, just warn and continue. The application
					// will continue using the latest creds and attempting to refresh them. If we're
					// never able to refresh the credentials again, the Approzium authenticator will
					// eventually return a failure response due to an expired GetCallerIdentity string.
					h.logger.Warnf("unable to refresh AWS get caller identity creds due to %s", err)
					continue
				}
				if h.logger.IsLevelEnabled(log.DebugLevel) {
					h.identityLock.RLock()
					h.logger.Debugf("identity refreshed to %+v", h.latestIdentity)
					h.identityLock.RUnlock()
				}
			}
		}
	}()
}

// This has been tested and confirmed to work with:
//		- A given role to assume (uses the role arn and succeeds).
//		- AWS EC2 instances that have been launched into a role (uses the role arn and succeeds).
func (h *awsIdentityHandler) refreshProof() error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	var svc *sts.STS
	if h.roleArnToAssume == "" {
		svc = sts.New(sess)
	} else {
		creds := stscreds.NewCredentials(sess, h.roleArnToAssume)
		svc = sts.New(sess, &aws.Config{Credentials: creds})
	}

	req, _ := svc.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{})
	signedGetCallerIdentity, err := req.Presign(time.Minute * 15)
	if err != nil {
		return err
	}

	callerIdentity, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}
	if callerIdentity.Arn == nil {
		return fmt.Errorf("no ARN received in %+v", callerIdentity)
	}

	id := &pb.AWSIdentity{
		SignedGetCallerIdentity: signedGetCallerIdentity,
		ClaimedIamArn:           *callerIdentity.Arn,
	}

	h.identityLock.Lock()
	defer h.identityLock.Unlock()
	h.latestIdentity = id
	return nil
}
