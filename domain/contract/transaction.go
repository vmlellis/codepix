package contract

import "github.com/vmlellis/imersao/codepix-go/domain/model"

type TransactionRepositoryInterface interface {
	Register(transaction *model.Transaction) error
	Save(transaction *model.Transaction) error
	Find(id string) (*model.Transaction, error)
}
