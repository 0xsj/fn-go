package middleware

import (
	"context"
	"time"

	"github.com/0xsj/fn-go/pkg/common/logging"
	"github.com/0xsj/fn-go/pkg/security"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(logger logging.Logger) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        start := time.Now()
        
        requestID := generateRequestID()
        ctx = context.WithValue(ctx, "request_id", requestID)
        
        logger.Info("Received request", 
            logging.F("method", info.FullMethod),
            logging.F("request_id", requestID))
        
        resp, err := handler(ctx, req)
        
        duration := time.Since(start)
        if err != nil {
            logger.Error("Request failed",
                logging.F("method", info.FullMethod),
                logging.F("request_id", requestID),
                logging.F("duration", duration.String()),
                logging.F("error", err.Error()))
        } else {
            logger.Info("Request succeeded",
                logging.F("method", info.FullMethod),
                logging.F("request_id", requestID),
                logging.F("duration", duration.String()))
        }
        
        return resp, err
    }
}


func AuthUnaryServerInterceptor(jwtConfig security.JWTConfig) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        if isPublicMethod(info.FullMethod) {
            return handler(ctx, req)
        }
        
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, status.Error(codes.Unauthenticated, "missing metadata")
        }
        
        authHeader, ok := md["authorization"]
        if !ok || len(authHeader) == 0 {
            return nil, status.Error(codes.Unauthenticated, "missing authorization header")
        }
        
        token := authHeader[0]
        claims, err := security.ValidateToken(jwtConfig, token)
        if err != nil {
            return nil, status.Error(codes.Unauthenticated, "invalid token")
        }
        
        ctx = context.WithValue(ctx, "user_id", claims.UserID)
        ctx = context.WithValue(ctx, "role", claims.Role)
        
        return handler(ctx, req)
    }
}



func generateRequestID() string {
    return time.Now().Format("20060102150405") + "-" + randomString(6)
}

func randomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
    }
    return string(b)
}

func isPublicMethod(method string) bool {
    publicMethods := map[string]bool{
        "/auth.AuthService/Login":  true,
        "/health.HealthService/Check": true,
    }
    return publicMethods[method]
}