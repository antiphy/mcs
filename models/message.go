package models

type Message struct {
	ID              int64
	Name            string // ??? store in mysql in case same phone number has different name
	PhoneNumber     string // index uniq index with CampaignID
	CampaignID      int64  // index
	Status          int8
	CreateTimestamp int64
	UpdateTimestamp int64
}
