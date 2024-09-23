package models

type Campaign struct {
	ID                 int
	CampaignID         string // uniq index, the client user specified campaign id
	Path               string
	Template           string
	ScheduledTimestamp int64
	Status             int8
	CreateTimestamp    int64
	UpdateTimestamp    int64
}

func (c *Campaign) TableName() string {
	return "tb_campaign"
}
