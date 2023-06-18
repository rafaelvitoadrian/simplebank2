package gapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/rafaelvitoadrian/simplebank2/db/sqlc"
	"github.com/rafaelvitoadrian/simplebank2/pb"
	"github.com/rafaelvitoadrian/simplebank2/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassowrd, err := utils.HashedPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed Hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
		HashedPassword: hashedPassowrd,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// log.Println(pqErr)
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "account already exist: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "error: %s", err)
	}
	// return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}
