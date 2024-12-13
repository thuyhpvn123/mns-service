package app

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/meta-node-blockchain/meta-node-mns/internal/config"
	"github.com/meta-node-blockchain/meta-node-mns/internal/controller"
	"github.com/meta-node-blockchain/meta-node-mns/internal/cronjob"
	"github.com/meta-node-blockchain/meta-node-mns/internal/database"
	"github.com/meta-node-blockchain/meta-node-mns/internal/handlers"
	"github.com/meta-node-blockchain/meta-node-mns/internal/repository"
	"github.com/meta-node-blockchain/meta-node-mns/internal/service"
	"github.com/meta-node-blockchain/meta-node-mns/internal/usecase"
	"github.com/meta-node-blockchain/meta-node-mns/route"
	"github.com/meta-node-blockchain/meta-node/cmd/client"
	c_config "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/types"
)

type App struct {
	Config      *config.AppConfig
	ApiApp      *gin.Engine
	ChainClient *client.Client
	StopChan    chan bool
	EventChan   chan types.EventLogs
	Handler     *handlers.Handler
	StorageClient *client.Client
}

func NewApp(
	configPath string,
	loglevel int,
) (*App, error) {
	var loggerConfig = &logger.LoggerConfig{
		Flag:    loglevel,
		Outputs: []*os.File{os.Stdout},
	}
	logger.SetConfig(loggerConfig)
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("can not load config", err)
	}
	app := &App{}
	engine := gin.Default()
	app.ChainClient, err = client.NewClient(
		&c_config.ClientConfig{
			Version_:                config.MetaNodeVersion,
			PrivateKey_:             config.PrivateKey_,
			ParentAddress:           config.NodeAddress,
			ParentConnectionAddress: config.NodeConnectionAddress,
			DnsLink_:                config.DnsLink(),
		},
	)

	if err != nil {
		logger.Error(fmt.Sprintf("error when create chain client %v", err))
		return nil, err
	}
	database.StartMySQL(config)
	db := database.GetMySqlConn()
	nameRepo := repository.NewNameRepository(db)

	nameUsecase := usecase.NewNameUsecase(nameRepo)

	app.StorageClient, err = client.NewStorageClient(
		&c_config.ClientConfig{
		  Version_:                config.MetaNodeVersion,
		  PrivateKey_:             config.PrivateKey_,
		  ParentAddress:           config.StorageAddress,
		  ParentConnectionAddress: config.NodeConnectionAddress,
		  DnsLink_:                config.DnsLink(),
		},
		[]common.Address{
			common.HexToAddress(config.NamewrapperAddress),
			common.HexToAddress(config.RegistrarControllerAddress),
		},
	)
	if err != nil {
		logger.Error(fmt.Sprintf("error when create storage client %v", err))
		return nil, err
	}
	//
	app.EventChan = app.StorageClient.GetEventLogsChan()
	// create customdomain abi
	reader, err := os.Open(config.CustomDomainABIPath)
	if err != nil {
		logger.Error("Error occured while read resolver abi")
		return nil, err
	}
	defer reader.Close()

	customDomainAbi, err := abi.JSON(reader)
	if err != nil {
		logger.Error("Error occured while parse resolver smart contract abi")
		return nil, err
	}
	// create namewrapper abi
	namewrapperReader, err := os.Open(config.NamewrapperABIPath)
	if err != nil {
		logger.Error("Error occured while read namewrapper abi")
		return nil, err
	}
	defer namewrapperReader.Close()

	nameWrapperAbi, err := abi.JSON(namewrapperReader)
	if err != nil {
		logger.Error("Error occured while parse namewrapper smart contract abi")
		return nil, err
	}
	// create registrarController abi
	registrarControllerReader, err := os.Open(config.RegistrarControllerABIPath)
	if err != nil {
		logger.Error("Error occured while read registrarController abi")
		return nil, err
	}
	defer registrarControllerReader.Close()

	registrarControllerAbi, err := abi.JSON(registrarControllerReader)
	if err != nil {
		logger.Error("Error occured while parse registrarController smart contract abi")
		return nil, err
	}
	
	// Initialize services
	servs := service.NewSendTransactionService(
		app.ChainClient,
		&customDomainAbi,
		common.HexToAddress(config.CustomDomainAddress),
		common.HexToAddress(config.ResolverAddress),
	)
	controller := controller.NewController(nameUsecase, servs)
	//Initialize Cronjob
	cronjob.Start(controller)
	defer cronjob.Stop()

	app.Config = config
	app.Handler = handlers.NewMNSHandler(
		nameUsecase,
		&nameWrapperAbi,
		&registrarControllerAbi,
	)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	//
	engine.Use(cors.New(corsConfig))

	route.InitialRoutes(
		engine,
		controller,
	)
	app.ApiApp = engine
	return app, nil
}

func (app *App) Run() {
	app.StopChan = make(chan bool)
	go func() {
		app.ApiApp.Run(app.Config.API_PORT)
	}()
	for {
		select {
		case <-app.StopChan:
			return
		case eventLogs := <-app.EventChan:
			app.Handler.HandleEvent(eventLogs)
		}
	}
}
  