package model

import "github.com/jinzhu/gorm"

type Images struct {
	gorm.Model
	XID         string `gorm:"column:xid;type:varchar(255);" json:"xid" form:"xid"`
	FileName    string `gorm:"column:file_name;type:varchar(255);" json:"fileName" form:"fileName"`
	URL         string `gorm:"column:url;type:varchar(255);" json:"url" form:"url"`
	ContentType string `gorm:"column:content_type;type:varchar(255);" json:"contentType" form:"contentType"`
	AccountID   string `gorm:"column:account_id;type:varchar(255);" json:"accountID" form:"accountID"`
}
