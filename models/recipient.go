package models

type Recipient struct {
	ID              int64
	PhoneNumber     string // uniq key
	Name            string
	CreateTimestamp int64
	UpdateTimestamp int64
}
