package factory

import (
	"github.com/vmlellis/imersao/codepix-go/application/usecase"
	"github.com/vmlellis/imersao/codepix-go/infrastructure/repository"
	"gorm.io/gorm"
)

func TransactionUseCaseFactory(database *gorm.DB) usecase.TransactionUseCase {
	pixRepository := repository.PixKeyRepositoryDb{Db: database}
	transactionRepository := repository.TransactionRepositoryDb{Db: database}

	transactionUseCase := usecase.TransactionUseCase{
		TransactionRepository: &transactionRepository,
		PixRepository:         &pixRepository,
	}

	return transactionUseCase
}
