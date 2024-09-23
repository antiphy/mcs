package dal

import "github.com/antiphy/mcs/models"

func (ds *datasource) CreateMessage(message models.Message) error {
	// tobe implement
	return nil
}

func (ds *datasource) QueryMessage(ID int64) (models.Message, error) {
	// tobe implement
	return models.Message{}, nil
}
func (ds *datasource) QueryMessageList(campaignID int) ([]models.Message, error) {
	// tobe implement
	return nil, nil
}

func (ds *datasource) UpdateMessageStatus(messageID int64, status int8) error {
	// tobe implement
	return nil
}
