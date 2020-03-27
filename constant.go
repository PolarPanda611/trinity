package trinity

// GRPCMethodKey used in passing grpc value
const GRPCMethodKey string = "method"

// TraceIDKey used in passing grpc value
const TraceIDKey string = "trace_id"

// ReqUserNameKey used in passing grpc value
const ReqUserNameKey string = "req_user_name"

// DefaultGRPCHealthCheck default grpc health check name
const DefaultGRPCHealthCheck string = "/grpc.health.v1.Health/Check"
