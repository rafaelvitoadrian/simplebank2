package gapi

import (
	db "github.com/rafaelvitoadrian/simplebank2/db/sqlc"
	"github.com/rafaelvitoadrian/simplebank2/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Password:          user.HashedPassword,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
