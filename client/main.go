package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/antiphy/mcs/consts"
	"github.com/antiphy/mcs/dal"
	"github.com/antiphy/mcs/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("usage: ./client [filepath] [Campaign_ID] [Message_template] [Scheduled_send_time]")
		fmt.Println("filepath can be a url started with [http://] or [https://]")
		return
	}

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

	kConfig := sarama.NewConfig()
	kConfig.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(cfg.Borkers, kConfig)
	if err != nil {
		panic(err)
	}

	var scanner *bufio.Scanner
	if isUrl(os.Args[1]) {
		response, err := http.Get(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		scanner = bufio.NewScanner(response.Body)
	} else {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Println("open csv file failed: " + err.Error())
			return
		}
		scanner = bufio.NewScanner(file)
	}
	db, err := gorm.Open(mysql.Open(cfg.DNS))
	if err != nil {
		panic(err)
	}
	datasource := dal.NewDatasource(db)

	fmt.Println("processing ...")

	// TODO: add some tips  while processing
	producerInput := producer.Input()
	go handleKafkaErrors(producer.Errors())
	var count int
	for scanner.Scan() {
		count++
		line := scanner.Text()
		if isHeader(line) {
			continue
		}
		var message models.Message
		//TODO: set message fields
		bts, _ := json.Marshal(message)
		kMessage := sarama.ProducerMessage{
			Topic: consts.TopicMessage,
			Value: sarama.ByteEncoder(bts),
		}

		producerInput <- &kMessage
	}
	producer.Close()

	// TODO: consider to create db record here in case that campaign id duplicated
	ScheduledSendTime, err := time.Parse(time.RFC3339, os.Args[4])
	if err != nil {
		panic("invalid Scheduled_send_time :" + os.Args[4])
	}
	campaign := models.Campaign{
		CampaignID:         os.Args[2],
		Path:               os.Args[1],
		Template:           os.Args[3],
		ScheduledTimestamp: ScheduledSendTime.Unix(),
		CreateTimestamp:    time.Now().Unix(),
	}
	err = datasource.CreateCampaign(&campaign)
	// hanle err

	fmt.Printf("process %d line completed", count-1)
}

func handleKafkaErrors(asyncErrors <-chan *sarama.ProducerError) {
	// TODO
}

func isUrl(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func isHeader(line string) bool {
	// TODO
	return false
}

type config struct {
	DNS     string   `json:"dsn"`
	Borkers []string `json:"borkers"`
}
