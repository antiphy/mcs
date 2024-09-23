package service

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/antiphy/mcs/consts"
	"github.com/antiphy/mcs/models"
	"go.uber.org/zap"
)

var stopConsumeSignal = make(chan struct{}, 1)

func StopConsume() {
	stopConsumeSignal <- struct{}{}
}

// TODO: use consumer group
func ConsumeMessage(datasource models.Datasource, consumer sarama.Consumer, offsetManager sarama.OffsetManager, logger zap.Logger, stopped chan struct{}) error {
	partionOffsetManager, _ := offsetManager.ManagePartition(consts.TopicMessage, 0)

	offset, _ := partionOffsetManager.NextOffset()
	partionConsumer, err := consumer.ConsumePartition(consts.TopicMessage, 0, offset)
	if err != nil {
		return err
	}

	consumerMessages := partionConsumer.Messages()

	select {
	case <-stopConsumeSignal:
		stopped <- struct{}{}
		return nil
	case kMessage := <-consumerMessages:
		nowTS := time.Now().Unix()
		var message models.Message
		message.CreateTimestamp = nowTS
		//
		message.Status = consts.MessageStatusNotSent
		err = json.Unmarshal(kMessage.Value, &message)
		if err != nil {
			// log
		}

		err = datasource.CreateMessage(message)
		if err != nil {
			// log
		}
		offset++
		partionOffsetManager.MarkOffset(offset, "")

		recipient := models.Recipient{
			PhoneNumber:     message.PhoneNumber,
			Name:            message.Name,
			CreateTimestamp: nowTS,
		}
		err = datasource.CreateRecipient(recipient)
		if err != nil {
			// log
		}
	}

	return nil
}
