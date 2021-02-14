package contract

import "github.com/vmlellis/imersao/codepix-go/domain/model"

type PixKeyRepositoryInterface interface {
	Register(pixKey *model.PixKey) (*model.PixKey, error)
	FindKeyByKind(key, kind string) (*model.PixKey, error)
	AddBank(bank *model.Bank) error
	AddAccount(account *model.Account)
	FindAccount(id string) (*model.Account, error)
}
