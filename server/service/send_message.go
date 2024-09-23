package service

import (
	"encoding/json"

	"github.com/antiphy/mcs/consts"
	"github.com/antiphy/mcs/messagesender"
	"github.com/antiphy/mcs/models"
	"github.com/antiphy/mcs/server/utils"
	"go.uber.org/zap"
)

var stopSendMessageSignal = make(chan struct{}, 1)

func StopSendMessage() {
	stopSendMessageSignal <- struct{}{}
}

func SendMessage(datasource models.Datasource, messageSender messagesender.MessageSender, logger zap.Logger, campaignChan <-chan models.Campaign, stopped chan struct{}) {
	select {
	case <-stopSendMessageSignal:
		stopped <- struct{}{}
		return
	case campaign := <-campaignChan:
		// TODO:
		// 1. control the count of goroutine in case too much memory cost
		// 2. wait the running goroutine to exit while stop signal acceptted
		go sendCampaignMessage(datasource, messageSender, logger, campaign)
	}
}

// TODO: add stop signal to monitor service stop
func sendCampaignMessage(datasource models.Datasource, messageSender messagesender.MessageSender, logger zap.Logger, campaign models.Campaign) {
	messages, _ := datasource.QueryMessageList(campaign.ID)
	instanceID := utils.GetInstanceID()
	_ = messages
	if instanceID == 0 {
		// its single instance
		// foreach message , send
	} else {
		instanceCount := utils.GetAllInstanceCount()
		// its mutil instance
		// only send message when message id%instanceCount == instanceID -1
		for i := range messages {
			if messages[i].ID%int64(instanceCount) == int64(instanceID)-1 {
				if messages[i].Status == consts.MessageStatusSent {
					continue
				}
				bts, _ := json.Marshal(messages[i])
				err := messageSender.SendMsg(bts, messages[i].PhoneNumber, messages[i].Name)
				// handle err
				_ = err
				datasource.UpdateMessageStatus(messages[i].ID, consts.MessageStatusSent)
			}
		}
	}
}
