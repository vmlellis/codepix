package grpc

import (
	"context"

	"github.com/vmlellis/imersao/codepix-go/application/grpc/pb"
	"github.com/vmlellis/imersao/codepix-go/application/usecase"
)

type PixGrpcService struct {
	PixUseCase usecase.PixUseCase
	pb.UnimplementedPixServiceServer
}

func (p *PixGrpcService) RegisterPixKey(ctx context.Context, in *pb.PixKeyRegistration) (*pb.PixKeyCreatedResult, error) {
	key, err := p.PixUseCase.RegisterKey(in.Key, in.Kind, in.AccountId)
	if err != nil {
		return &pb.PixKeyCreatedResult{
			Status: "not created",
			Error:  err.Error(),
		}, err
	}

	return &pb.PixKeyCreatedResult{
		Id:     key.ID,
		Status: "created",
	}, nil
}

func (p *PixGrpcService) Find(ctx context.Context, in *pb.PixKey) (*pb.PixKeyInfo, error) {
	pixKey, err := p.PixUseCase.FindKey(in.Key, in.Kind)
	if err != nil {
		return &pb.PixKeyInfo{}, err
	}

	var account *pb.Account
	if pixKey.Account != nil {
		acc := pixKey.Account
		bankName := ""
		if acc.Bank != nil {
			bankName = acc.Bank.Name
		}

		account = &pb.Account{
			AccountId:     pixKey.AccountID,
			AccountNumber: acc.Number,
			BankId:        acc.BankID,
			BankName:      bankName,
			OwnerName:     acc.OwnerName,
			CreatedAt:     acc.CreatedAt.String(),
		}
	}

	return &pb.PixKeyInfo{
		Id:        pixKey.ID,
		Kind:      pixKey.Kind,
		Key:       pixKey.Key,
		Account:   account,
		CreatedAt: pixKey.CreatedAt.String(),
	}, nil
}

func NewPixGrpcService(usecase usecase.PixUseCase) *PixGrpcService {
	return &PixGrpcService{
		PixUseCase: usecase,
	}
}
