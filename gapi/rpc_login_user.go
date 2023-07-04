package gapi

import (
	"context"
	"database/sql"

	db "github.com/rafaelvitoadrian/simplebank2/db/sqlc"
	"github.com/rafaelvitoadrian/simplebank2/pb"
	"github.com/rafaelvitoadrian/simplebank2/utils"
	"github.com/rafaelvitoadrian/simplebank2/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := ValidateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Not Found Account: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "Error : %s", err)
	}

	err = utils.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Password Or User Wrong: %s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessDurationToken)
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "Failed To Create Access Token")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed To Create Refresh Token")
	}

	mtdt := server.extractMetaData(ctx)
	Session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed To Create Session")
	}

	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             Session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiredAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiredAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return rsp, nil
}

func ValidateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValdiateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
