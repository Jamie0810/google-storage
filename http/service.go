package service

import (
	"context"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
	"github.com/rs/xid"
	"gitlab.silkrode.com.tw/golang/errors"
)

func (g *GCPService) uploadFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	// -----check if there's something wrong with the file-----
	file, err := header.Open()
	if err != nil {
		return "", errors.Wrap(err, "os.Open Err")
	}
	defer file.Close()

	// -----upload the object with storage.Writer-----
	ctx, cancel := context.WithTimeout(g.ctx, time.Second*50)
	defer cancel()

	objectName := xid.New().String()
	writer := g.bucket.Object(objectName).NewWriter(ctx)
	if _, err = io.Copy(writer, file); err != nil {
		return "", errors.Wrap(err, "io.Copy Err")
	}

	if err := writer.Close(); err != nil {
		return "", errors.Wrap(err, "writer.Close Err")
	}

	g.logger.Info().Msgf("The file has been uploaded successfully.")
	return objectName, nil
}

func (g *GCPService) downloadFile(objectName string) (*storage.Reader, error) {
	ctx, cancel := context.WithTimeout(g.ctx, time.Second*50)
	defer cancel()

	reader, err := g.bucket.Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to download the object.")
	}

	g.logger.Info().Msgf("The file has been downloaded successfully.")
	return reader, nil
}
