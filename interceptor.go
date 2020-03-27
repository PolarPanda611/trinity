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
func LoggingInterceptor(log Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == DefaultGRPCHealthCheck {
			return handler(ctx, req)
		}
		md, _ := metadata.FromIncomingContext(ctx)
		md.Append(GRPCMethodKey, info.FullMethod)
		if _, ok := md[TraceIDKey]; !ok {
			md.Append(TraceIDKey, "")
		}
		if _, ok := md[ReqUserNameKey]; !ok {
			md.Append(ReqUserNameKey, "")
		}
		resp, err := handler(ctx, req)
		log.FormatLogger(
			GRPCMethod(info.FullMethod),
			TraceID(md[TraceIDKey][0]),
			ReqUserName(md[ReqUserNameKey][0]),
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
		if info.FullMethod == DefaultGRPCHealthCheck {
			return handler(ctx, req)
		}

		md, _ := metadata.FromIncomingContext(ctx)
		currentUser, ok := md[ReqUserNameKey]
		if !ok || currentUser[0] == "" {
			return nil, status.Error(codes.Unauthenticated, "request user name not found")
		}
		return handler(ctx, req)
	}
}
