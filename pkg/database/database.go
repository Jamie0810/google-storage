package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"gitlab.silkrode.com.tw/team_golang/kbc/km/storage/config"
	"gitlab.silkrode.com.tw/team_golang/kbc/km/storage/model"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitDB(cnf config.Config) (db *gorm.DB, err error) {
	var connInfo string

	if cnf.Database.InstanceName == "" {
		connInfo = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			cnf.Database.User,
			cnf.Database.Password,
			cnf.Database.Host,
			cnf.Database.Port,
			cnf.Database.DBName,
		)
	} else {
		connInfo = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?charset=utf8mb4&parseTime=true&loc=UTC&time_zone=UTC",
			cnf.Database.User,
			cnf.Database.Password,
			cnf.Database.InstanceName,
			cnf.Database.DBName,
		)
	}

	db, err = gorm.Open("mysql", connInfo)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Images{})
	return db, nil
}
