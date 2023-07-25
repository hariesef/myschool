package account

import (
	"context"

	acct "myschool/internal/services/account"
	pb "myschool/pkg/controller/rpc/account"
)

type AccountRPCServer struct {
	AccountSvc *acct.AccountService
}

func (s *AccountRPCServer) Create(ctx context.Context, args *pb.NewAccountParam) (*pb.Empty, error) {
	return &pb.Empty{}, s.AccountSvc.Create(ctx, args.Email, args.Password)
}

func (s *AccountRPCServer) Login(ctx context.Context, args *pb.LoginInfo) (*pb.TokenInfo, error) {
	result, err := s.AccountSvc.Login(ctx, args.Email, args.Password)
	if err != nil {
		return &pb.TokenInfo{}, err
	}
	return &pb.TokenInfo{Token: result.Token, Expiry: int32(result.Expiry)}, nil
}

func (s *AccountRPCServer) Logout(ctx context.Context, args *pb.Token) (*pb.Empty, error) {
	return &pb.Empty{}, s.AccountSvc.Logout(ctx, args.Token)
}
