package main

import (
	"encoding/json"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/antiphy/mcs/consts"
	"github.com/antiphy/mcs/dal"
	"github.com/antiphy/mcs/messagesender"
	"github.com/antiphy/mcs/models"
	"github.com/antiphy/mcs/server/service"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type config struct {
	DNS     string   `json:"dsn"`
	Borkers []string `json:"borkers"`
}

func main() {
	// config file or env or etcd
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	var cfg config
	err = json.NewDecoder(configFile).Decode(&cfg)
	if err != nil {
		panic("invalid cfg")
	}
	db, err := gorm.Open(mysql.Open(cfg.DNS))
	if err != nil {
		panic(err)
	}
	if len(os.Args) > 1 && os.Args[1] == "init" {
		initDB(db)
		return
	}

	logConfig := zap.NewProductionConfig()
	logConfig.OutputPaths = []string{"./consumer.log"}
	logConfig.ErrorOutputPaths = []string{"./consumer_err.log"}
	logConfig.Level.SetLevel(zap.ErrorLevel)
	logger, err := logConfig.Build(zap.WithCaller(true))
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}

	borkers := []string{}
	config := sarama.NewConfig()
	client, _ := sarama.NewClient(borkers, config)
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		logger.Fatal("new consumer failed: " + err.Error())
		return
	}
	offsetManager, err := sarama.NewOffsetManagerFromClient(consts.GroupIDMCSServer, client)

	datasource := dal.NewDatasource(db)
	messageSender := messagesender.NewMsgSender()

	// used to check if consume/sendmsg stopped
	consumeStopped := make(chan struct{}, 1)
	sendMessageStopped := make(chan struct{}, 1)
	campaignChannel := make(chan models.Campaign, 10)

	go service.CampaignMonitor(datasource, campaignChannel, *logger)
	go service.ConsumeMessage(datasource, consumer, offsetManager, *logger, consumeStopped)
	go service.SendMessage(datasource, messageSender, *logger, campaignChannel, sendMessageStopped)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	logger.Info("interrupt signal acceptted, stopping.....")

	service.StopConsume()
	service.StopSendMessage()

	<-consumeStopped
	<-sendMessageStopped

	os.Exit(0)
}

func initDB(db *gorm.DB) {
	// check if db exists
	err := db.Exec(`create database 'mcs' if not exists`).Error
	if err != nil {
		panic(err)
	}
	// two ways to init table:
	// 1. use gorm.DB.Create
	err = db.Migrator().CreateTable(&models.Campaign{})
	if err != nil {
		panic(err)
	}
	err = db.Migrator().CreateTable(&models.Message{})
	if err != nil {
		panic(err)
	}
	err = db.Migrator().CreateTable(&models.Recipient{})
	if err != nil {
		panic(err)
	}
	// 2. use db.Exec execute sql script
}
