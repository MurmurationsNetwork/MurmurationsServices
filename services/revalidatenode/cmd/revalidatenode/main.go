package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/revalidatenode/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/revalidatenode/internal/adapter/repository"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/revalidatenode/internal/usecase"
)

func init() {
	global.Init()
}

func main() {
	nodeUsecase := usecase.NewNodeUsecase((repository.NewNodeRepository(mongo.Client.GetClient())))
	nodeUsecase.RevalidateNodes()
	mongo.Client.Disconnect()
}
