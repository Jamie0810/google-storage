package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gitlab.silkrode.com.tw/team_golang/kbc/km/storage/config"
	"gitlab.silkrode.com.tw/team_golang/kbc/km/storage/model"
	"gitlab.silkrode.com.tw/team_golang/kbc/km/storage/pkg/log"
	pagination "gitlab.silkrode.com.tw/team_golang/kbc/km/storage/utils"
	"google.golang.org/api/option"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
)

type GCPService struct {
	config     config.Config
	logger     *log.Logger
	db         *gorm.DB
	client     *storage.Client
	bucketName string
	bucket     *storage.BucketHandle
	w          io.Writer
	ctx        context.Context
}

func InitGCPService(config config.Config, logger *log.Logger, db *gorm.DB) (*GCPService, error) {
	ctx := context.Background()
	// -----init a client-----
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.GoogleStorage.StorageKeyFile))
	if err != nil {
		logger.Error().Msgf("Failed to create a client.", err)
		return nil, err
	}

	// -----check if a bucket exists-----
	bucket := client.Bucket(config.GoogleStorage.BucketName)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		// -----create a bucket in Google Storage-----
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		if err := bucket.Create(ctx, config.GoogleStorage.ProjectID, &storage.BucketAttrs{
			Location: "asia",
		}); err != nil {
			logger.Error().Msgf("Failed to create the storage bucket.", err)
			return nil, err
		}
		logger.Info().Msgf("A bucket '" + config.GoogleStorage.BucketName + "' has been created.")
	}

	// -----enable public access to the bucket-----
	policy, err := bucket.IAM().Policy(ctx)
	if err != nil {
		return nil, err
	}
	policy.InternalProto.Bindings = append(policy.InternalProto.Bindings, &iampb.Binding{
		Role:    "roles/storage.objectViewer",
		Members: []string{"allUsers"},
	})
	if err := bucket.IAM().SetPolicy(ctx, policy); err != nil {
		return nil, err
	}

	gcpService := GCPService{
		config:     config,
		logger:     logger,
		db:         db,
		client:     client,
		bucketName: config.GoogleStorage.BucketName,
		bucket:     bucket,
		ctx:        ctx,
	}
	return &gcpService, nil
}

func (g *GCPService) UploadImage(c *gin.Context) {
	// -----read the file-----
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("File err : %s", err.Error()))
		return
	}
	// -----upload the file to Google storage-----
	objectName, err := g.uploadFile(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload the image."})
		return
	}

	// -----store the URL into DB-----
	fileURL := g.config.GoogleStorage.URL + g.bucketName + "/" + objectName
	image := model.Images{
		URL:         fileURL,
		XID:         objectName,
		FileName:    header.Filename,
		ContentType: strings.Join(header.Header["Content-Type"], "; "),
		AccountID:   c.Param("accountID"),
	}

	if err := g.db.Create(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to store the URL into db."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Image xid": image.XID})
}

func (g *GCPService) GetImageList(c *gin.Context) {
	_page, _ := c.GetQuery("page")
	_item, _ := c.GetQuery("item")

	page, err := strconv.ParseInt(_page, 10, 64)
	if err == nil {
		fmt.Println("strconv.ParseInt Err")
	}

	item, err := strconv.ParseInt(_item, 10, 64)
	if err == nil {
		fmt.Println("strconv.ParseInt Err")
	}

	p := pagination.Pagination{
		Page:    page,
		PerPage: item,
	}

	limit, offset := p.CheckOrSetDefault().LimitAndOffset()

	if !g.config.IsAuthorized {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not allowed to get image list."})
		return
	}
	var images []model.Images

	if err := g.db.Offset(offset).Limit(limit).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Can not get image list."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": images})
}

func (g *GCPService) DownloadImage(c *gin.Context) {
	xid := c.Param("xid")
	image := model.Images{}

	if err := g.db.Where("xid = ?", xid).Find(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Can not find the image."})
		return
	}

	fileReader, err := g.downloadFile(image.XID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to download the image."})
		return
	}

	c.Status(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename="+image.FileName)
	c.Header("Content-Type", image.ContentType)
	c.Header("X-Secret", image.AccountID)
	c.Stream(func(w io.Writer) bool {
		_, err := io.Copy(w, fileReader)
		if err != nil {
			g.logger.Debug().Err(err).Msg("stream io.Copy error")
		}
		return false
	})
}

func (g *GCPService) DeleteImage(c *gin.Context) {
	if !g.config.IsAuthorized {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not allowed to delete any images."})
		return
	}

	var image model.Images

	if g.db.First(&image, c.Param("id")).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Can not fine the image"})
		return
	}

	if err := g.db.Unscoped().Delete(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "The image has been deleted."})
}
