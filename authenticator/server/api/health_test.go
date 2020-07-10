package api

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/cyralinc/approzium/authenticator/server/config"
	pb "github.com/cyralinc/approzium/authenticator/server/protos"
	testtools "github.com/cyralinc/approzium/authenticator/server/testing"
	"google.golang.org/grpc"
)

func TestHealthChecker(t *testing.T) {
	c := config.Config{
		Host:     "127.0.0.1",
		GRPCPort: 6001,
	}
	serviceAddress := fmt.Sprintf("%s:%d", c.Host, c.GRPCPort)
	lis, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		t.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterHealthServer(s, &testHealthServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Fatalf("Server exited with error: %v", err)
		}
	}()

	checker := newHealthChecker(testtools.TestLogger(), c)
	testWriter := &testtools.TestResponseWriter{}
	checker.ServeHTTP(testWriter, nil)

	if testWriter.LastStatusCodeReceived != 200 {
		t.Fatalf("expected 200 but received %d", testWriter.LastStatusCodeReceived)
	}
}

type testHealthServer struct{}

func (h *testHealthServer) Check(context.Context, *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}
