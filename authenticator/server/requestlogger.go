package server

import (
	"context"
	"encoding/json"
	"time"

	pb "github.com/approzium/approzium/authenticator/server/protos"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ctxLogger     = "logger"
	redactedValue = "********"
)

// getRequestLogger gets a logger that includes the request ID
// with every line, and a trace ID if it was sent.
func getRequestLogger(ctx context.Context) *log.Entry {
	rawLogger := ctx.Value(ctxLogger)
	if rawLogger == nil {
		return &log.Entry{}
	}
	logger, ok := rawLogger.(*log.Entry)
	if !ok {
		return &log.Entry{}
	}
	return logger
}

func newRequestLogger(logger *log.Logger, logRaw bool, wrapped pb.AuthenticatorServer) pb.AuthenticatorServer {
	return &requestLogger{
		logger:  logger,
		logRaw:  logRaw,
		wrapped: wrapped,
	}
}

type requestLogger struct {
	// logger is the application-level logger. It holds settings like the log format and level,
	// but lacks request-related context.
	logger *log.Logger

	// If true, will not redact fields that appear to hold sensitive information.
	// Defaults to false.
	logRaw bool

	// wrapped is the AuthenticatorServer we're logging requests and responses for.
	wrapped pb.AuthenticatorServer
}

func (l *requestLogger) GetPGMD5Hash(ctx context.Context, req *pb.PGMD5HashRequest) (*pb.PGMD5Response, error) {
	requestId, requestLogger, err := l.buildContextualLogger()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	l.logRequest(requestLogger, req, req.Awsauth)

	resp, respErr := l.wrapped.GetPGMD5Hash(context.WithValue(ctx, ctxLogger, requestLogger), req)

	if resp == nil {
		resp = &pb.PGMD5Response{}
	}
	resp.Requestid = requestId
	l.logResponse(requestLogger, resp, respErr)
	return resp, respErr
}

func (l *requestLogger) GetPGSHA256Hash(ctx context.Context, req *pb.PGSHA256HashRequest) (*pb.PGSHA256Response, error) {
	requestId, requestLogger, err := l.buildContextualLogger()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	l.logRequest(requestLogger, req, req.Awsauth)

	resp, respErr := l.wrapped.GetPGSHA256Hash(context.WithValue(ctx, ctxLogger, requestLogger), req)

	if resp == nil {
		resp = &pb.PGSHA256Response{}
	}
	resp.Requestid = requestId
	l.logResponse(requestLogger, resp, respErr)
	return resp, respErr
}

func (l *requestLogger) logRequest(requestLogger *log.Entry, req interface{}, awsAuth *pb.AWSAuth) {
	// Log asynchronously to avoid blocking while lots of JSON conversion takes place.
	preciseTime := time.Now().UTC()
	go func() {
		reqMap, err := toMap(req)
		if err != nil {
			requestLogger.Warnf("couldn't log request due to %s", err)
			return
		}
		if !l.logRaw && awsAuth != nil {
			if _, ok := reqMap["awsauth"]; !ok {
				requestLogger.Warn("couldn't log request because aws auth not found")
				return
			}
			reqMap["awsauth"] = map[string]interface{}{
				"claimed_iam_arn":            awsAuth.ClaimedIamArn,
				"signed_get_caller_identity": redactedValue,
			}
		}
		fields := log.Fields{
			"precise_time": preciseTime.Format(time.RFC3339Nano),
		}
		for k, v := range reqMap {
			fields[k] = v
		}
		requestLogger.WithFields(fields).Info("request")
	}()
}

func (l *requestLogger) logResponse(requestLogger *log.Entry, resp interface{}, respErr error) {
	// Log asynchronously to avoid blocking while lots of JSON conversion takes place.
	preciseTime := time.Now().UTC()
	go func() {
		respMap, err := toMap(resp)
		if err != nil {
			requestLogger.Warnf("couldn't log response due to %s", err)
			return
		}
		fields := log.Fields{
			"precise_time": preciseTime.Format(time.RFC3339Nano),
		}
		for k, v := range respMap {
			fields[k] = v
		}
		if respErr != nil {
			fields["error"] = respErr.Error()
			requestLogger.WithFields(fields).Error("response")
			return
		}
		requestLogger.WithFields(fields).Info("response")
	}()
}

func (l *requestLogger) buildContextualLogger() (string, *log.Entry, error) {
	randomUuid, err := uuid.NewRandom()
	if err != nil {
		return "", nil, err
	}
	requestId := randomUuid.String()

	// Create the logger from the application-level one to retain its settings.
	requestLogger := l.logger.WithFields(log.Fields{
		"request_id": requestId,
	})
	return requestId, requestLogger, nil
}

func toMap(obj interface{}) (map[string]interface{}, error) {
	objJson, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	objMap := make(map[string]interface{})
	if err := json.Unmarshal(objJson, &objMap); err != nil {
		return nil, err
	}
	return objMap, nil
}
