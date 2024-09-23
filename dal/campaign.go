package dal

import (
	"github.com/antiphy/mcs/dal/errors"
	"github.com/antiphy/mcs/models"
)

func (ds *datasource) CreateCampaign(campaign *models.Campaign) error {
	err := ds.mysqlDatasource.CreateCampaign(campaign)
	if err != nil {
		return err
	}

	return ds.memoryDatasource.CreateCampaign(campaign)
}

func (ds *datasource) QueryCampaign(campaignID string) (models.Campaign, error) {
	campaign, err := ds.memoryDatasource.QueryCampaign(campaignID)
	if err == errors.ErrRecordNotFound {
		campaign, dbErr := ds.mysqlDatasource.QueryCampaign(campaignID)
		if dbErr == nil {
			ds.memoryDatasource.CreateCampaign(&campaign)
		}
		return campaign, dbErr
	}
	return campaign, err
}

func (ds *datasource) QueryCampaignList(status int8) ([]models.Campaign, error) {
	return ds.mysqlDatasource.QueryCampaignList(status)
}
