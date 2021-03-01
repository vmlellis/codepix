package kafka

import (
	"fmt"
	"os"

	"gorm.io/gorm"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/vmlellis/imersao/codepix-go/application/factory"
	appModel "github.com/vmlellis/imersao/codepix-go/application/model"
	"github.com/vmlellis/imersao/codepix-go/application/usecase"
	"github.com/vmlellis/imersao/codepix-go/domain/model"
)

type KafkaProcessor struct {
	Database      *gorm.DB
	KafkaProducer *KafkaProducer
}

func NewKafkaProcessor(database *gorm.DB, producer *KafkaProducer) *KafkaProcessor {
	return &KafkaProcessor{
		Database:      database,
		KafkaProducer: producer,
	}
}

func (k *KafkaProcessor) Consume() error {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
		"group.id":          os.Getenv("kafkaConsumerGroupId"),
		"auto.offset.reset": "earliest",
	}
	c, err := ckafka.NewConsumer(configMap)
	if err != nil {
		return err
	}

	topics := []string{
		os.Getenv("kafkaTransactionTopic"),
		os.Getenv("kafkaTransactionConfirmationTopic"),
	}
	c.SubscribeTopics(topics, nil)

	fmt.Println("kafka consumer has been started")
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			k.processMessage(msg)
		}
	}
}

func (k *KafkaProcessor) processMessage(msg *ckafka.Message) {
	fmt.Println(string(msg.Value))

	transactionsTopic := os.Getenv("kafkaTransactionTopic")
	transactionConfirmationTopic := os.Getenv("kafkaTransactionConfirmationTopic")

	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionsTopic:
		k.processTransaction(msg)
	case transactionConfirmationTopic:
		k.processTransactionConfirmation(msg)
	default:
		fmt.Println("not a valid topic", string(msg.Value))
	}
}

func (k *KafkaProcessor) processTransaction(msg *ckafka.Message) error {
	fmt.Println("===> processTransaction")
	fmt.Println(string(msg.Value))

	transaction := appModel.NewTransaction()
	err := transaction.ParseJson(msg.Value)
	if err != nil {
		fmt.Println("error parse transaction", err)
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(k.Database)

	createdTransaction, err := transactionUseCase.Register(
		transaction.AccountID,
		transaction.Amount,
		transaction.PixKeyTo,
		transaction.PixKeyKindTo,
		transaction.Description,
	)
	if err != nil {
		fmt.Println("error registering transaction", err)
		return err
	}

	topic := "bank" + createdTransaction.PixKeyTo.Account.Bank.Code
	transaction.ID = createdTransaction.ID
	transaction.Status = model.TransactionPending
	transactionJson, err := transaction.ToJson()
	if err != nil {
		return err
	}

	err = k.KafkaProducer.Publish(string(transactionJson), topic)
	if err != nil {
		return err
	}

	return nil
}

func (k *KafkaProcessor) processTransactionConfirmation(msg *ckafka.Message) error {
	fmt.Println("===> processTransactionConfirmation")
	fmt.Println(string(msg.Value))

	transaction := appModel.NewTransaction()
	err := transaction.ParseJson(msg.Value)
	if err != nil {
		fmt.Println("error parse transaction confirmation", err)
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(k.Database)

	if transaction.Status == model.TransactionConfirmed {
		fmt.Println("===> confirm transaction")
		err = k.confirmTransaction(transaction, transactionUseCase)
		if err != nil {
			return err
		}
	} else if transaction.Status == model.TransactionCompleted {
		fmt.Println("===> complete transaction")
		_, err := transactionUseCase.Complete(transaction.ID)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (k *KafkaProcessor) confirmTransaction(transaction *appModel.Transaction, transactionUseCase usecase.TransactionUseCase) error {
	confirmedTransaction, err := transactionUseCase.Confirm(transaction.ID)
	if err != nil {
		fmt.Println("error to confirm transaction", err)
		return err
	}

	topic := "bank" + confirmedTransaction.AccountFrom.Bank.Code
	transactionJson, err := transaction.ToJson()
	if err != nil {
		fmt.Println("error transaction to json", err)
		return err
	}

	fmt.Println("publishing:", string(transactionJson))
	err = k.KafkaProducer.Publish(string(transactionJson), topic)
	if err != nil {
		fmt.Println("error publish in kafka", err)
		return err
	}

	return nil
}
