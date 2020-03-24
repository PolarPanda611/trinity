package trinity

import (
	"context"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor record log
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	// md, _ := metadata.FromIncomingContext(ctx)

	// logger.Logger.Log("gRPC method: ", info.FullMethod, "trace_id :", md["trace_id"][0], "currentUser :", md["current_user"][0], "request :", fmt.Sprintf("%v", req), "response :", fmt.Sprintf("%v", resp), "error: ", fmt.Sprintf("%v", err))
	return resp, err
}

// RecoveryInterceptor recovery from panic
func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
			// logger.Logger.Log("gRPC method: ", info.FullMethod, "request :", fmt.Sprintf("%v", req), "Panic err :", fmt.Sprintf("%v", err))
		}
	}()

	return handler(ctx, req)
}
