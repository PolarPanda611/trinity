package trinity

import (
	"context"
	"fmt"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor record log
// func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 	resp, err := handler(ctx, req)
// 	// md, _ := metadata.FromIncomingContext(ctx)

// 	// logger.Logger.Log("gRPC method: ", info.FullMethod, "trace_id :", md["trace_id"][0], "currentUser :", md["current_user"][0], "request :", fmt.Sprintf("%v", req), "response :", fmt.Sprintf("%v", resp), "error: ", fmt.Sprintf("%v", err))
// 	return resp, err
// }

// // RecoveryInterceptor recovery from panic
// func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
// 	defer func() {
// 		if e := recover(); e != nil {
// 			debug.PrintStack()
// 			err = status.Errorf(codes.Internal, "Panic err: %v", e)
// 			// logger.Logger.Log("gRPC method: ", info.FullMethod, "request :", fmt.Sprintf("%v", req), "Panic err :", fmt.Sprintf("%v", err))
// 		}
// 	}()
// 	return handler(ctx, req)
// }

// LoggingInterceptor record log
func LoggingInterceptor(log Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}
		md, _ := metadata.FromIncomingContext(ctx)
		md.Append("method", info.FullMethod)
		if _, ok := md["trace_id"]; !ok {
			md.Append("trace_id", "")
		}
		if _, ok := md["req_user_name"]; !ok {
			md.Append("req_user_name", "")
		}
		resp, err := handler(ctx, req)
		log.FormatLogger(
			GRPCMethod(info.FullMethod),
			TraceID(md["trace_id"][0]),
			ReqUserName(md["req_user_name"][0]),
		).Print("Req", fmt.Sprintf("%v", req), "Res", fmt.Sprintf("%v", resp), "Error", fmt.Sprintf("%v", err))
		return resp, err
	}
}

// RecoveryInterceptor recovery from panic
func RecoveryInterceptor(log Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if e := recover(); e != nil {
				debug.PrintStack()
				err = status.Errorf(codes.Internal, "Panic err: %v", e)
				log.Print("gRPC method: ", info.FullMethod, "request :", fmt.Sprintf("%v", req), "Panic err :", fmt.Sprintf("%v", err))
			}
		}()
		return handler(ctx, req)
	}
}

// UserAuthInterceptor record log
func UserAuthInterceptor(log Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		md, _ := metadata.FromIncomingContext(ctx)
		currentUser, ok := md["current_user"]
		if !ok || currentUser[0] == "" {
			return nil, status.Error(codes.Unauthenticated, "current user not found")
		}
		return handler(ctx, req)
	}
}
