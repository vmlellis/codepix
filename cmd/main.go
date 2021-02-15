package main

import (
	"os"

	"github.com/vmlellis/imersao/codepix-go/application/grpc"
	"github.com/vmlellis/imersao/codepix-go/infrastructure/db"
	"gorm.io/gorm"
)

var database *gorm.DB

func main() {
	database := db.ConnectDB(os.Getenv("env"))
	grpc.StartGrpcServer(database, 50051)
}
