package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/rafaelvitoadrian/simplebank2/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorization(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("Missing Metadata")
	}

	values := md.Get(authorizationHeader)

	if len(values) == 0 {
		return nil, fmt.Errorf("Missing Token")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("Invalid Authorization Format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsuport authorization type : %s", authType)
	}

	accesToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accesToken)
	if err != nil {
		return nil, fmt.Errorf("Invalid Access Token : %s", err)
	}

	return payload, nil
}
