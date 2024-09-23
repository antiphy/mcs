package models

type Datasource interface {
	CreateCampaign(campaign *Campaign) error
	QueryCampaign(campaignID string) (Campaign, error)
	QueryCampaignList(status int8) ([]Campaign, error)

	CreateMessage(message Message) error
	QueryMessage(ID int64) (Message, error)
	QueryMessageList(campaignID int) ([]Message, error)
	UpdateMessageStatus(messageID int64, status int8) error

	CreateRecipient(recipient Recipient) error
}
