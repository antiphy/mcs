package service

import (
	"time"

	"github.com/antiphy/mcs/consts"
	"github.com/antiphy/mcs/models"
	"go.uber.org/zap"
)

func CampaignMonitor(datasource models.Datasource, campaignChan chan<- models.Campaign, logger zap.Logger) {
	ticker := time.NewTicker(1 * time.Second)

	for {
		now := <-ticker.C
		campaigns, err := datasource.QueryCampaignList(consts.CampaignStatusCreated)
		if err != nil {
			// log
		}
		for i := range campaigns {
			if campaigns[i].ScheduledTimestamp >= now.Unix() {
				campaignChan <- campaigns[i]
			}
		}
	}

}
