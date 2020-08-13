package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gitlab.silkrode.com.tw/team_golang/kbc/km/storage/config"
	http "gitlab.silkrode.com.tw/team_golang/kbc/km/storage/http"
	"gitlab.silkrode.com.tw/team_golang/kbc/km/storage/pkg/log"
)

func InitRouter(config config.Config, logger *log.Logger, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	svc, err := http.InitGCPService(config, logger, db)
	if err != nil {
		logger.Fatal().Msgf("Failed init GCP Service.", err)
	}

	r.POST("/images/:accountID", svc.UploadImage)
	r.GET("/images/:xid", svc.DownloadImage)
	r.GET("/images", svc.GetImageList)
	r.DELETE("/images/:id", svc.DeleteImage)

	return r
}
