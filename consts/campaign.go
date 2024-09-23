package consts

const (
	CampaignPathTypeFile = "file"
	CampaignPathTypeUrl  = "url"
)

const (
	// just created in db
	CampaignStatusCreated int8 = 1
	// sending message
	CampaignStatusProcessing int8 = 1
	// send all message completed
	CampaignStatusCompleted int8 = 2
)
