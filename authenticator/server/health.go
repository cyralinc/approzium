package server

import (
	"context"

	pb "github.com/cyralinc/approzium/authenticator/server/protos"
)

func newHealthServer() pb.HealthServer {
	return &healthServer{}
}

type healthServer struct{}

func (h *healthServer) Check(context.Context, *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	// Someday we may wish to grow more sophisticated in our checks, but for now we're simply confirming
	// that the gRPC server is up.
	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}
