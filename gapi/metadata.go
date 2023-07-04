package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	userAgentGetCTX     = "grpcgateway-user-agent"
	userAgentHeaderGRPC = "user-agent"
	ipClientGetCTX      = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetaData(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// log.Print("md: %+v\n", md)
		if userAgents := md.Get(userAgentGetCTX); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeaderGRPC); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if ClientIPs := md.Get(ipClientGetCTX); len(ClientIPs) > 0 {
			mtdt.ClientIP = ClientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
