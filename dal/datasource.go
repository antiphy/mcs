package dal

import (
	"github.com/antiphy/mcs/models"
	"gorm.io/gorm"
)

type datasource struct {
	mysqlDatasource  models.Datasource
	memoryDatasource models.Datasource
}

func NewDatasource(db *gorm.DB) models.Datasource {
	return nil
}
