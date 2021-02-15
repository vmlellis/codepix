package usecase

import (
	"github.com/vmlellis/imersao/codepix-go/domain/contract"
	"github.com/vmlellis/imersao/codepix-go/domain/model"
)

type PixUseCase struct {
	PixKeyRepository contract.PixKeyRepositoryInterface
}

func (p *PixUseCase) RegisterKey(key, kind, accountId string) (*model.PixKey, error) {
	account, err := p.PixKeyRepository.FindAccount(accountId)
	if err != nil {
		return nil, err
	}

	pixKey, err := model.NewPixKey(kind, account, key)
	if err != nil {
		return nil, err
	}

	_, err = p.PixKeyRepository.RegisterKey(pixKey)
	if err != nil {
		return nil, err
	}

	return pixKey, nil
}

func (p *PixUseCase) FindKey(key string, kind string) (*model.PixKey, error) {
	pixKey, err := p.PixKeyRepository.FindKeyByKind(key, kind)
	if err != nil {
		return nil, err
	}

	return pixKey, nil
}
