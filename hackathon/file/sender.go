package file

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"elotus.com/hackathon/auth"
	"elotus.com/hackathon/storage"
	"elotus.com/hackathon/storage/query"
	"go.uber.org/zap"
)

var (
	ErrEmptyFile     = fmt.Errorf("file must not be empty")
	ErrImageFileOnly = fmt.Errorf("image file only")
)

var AcceptedImagesContentTypes = map[string]struct{}{
	"image/jpeg":    {},
	"image/png":     {},
	"image/gif":     {},
	"image/bmp":     {},
	"image/tiff":    {},
	"image/webp":    {},
	"image/svg+xml": {},
	"image/avif":    {},
}

type Sender interface {
	Upload(context.Context, *UploadFileRequest) (*UploadFileResponse, error)
}

type sender struct {
	lg      *zap.Logger
	storage storage.Storage
	dir     string
}

func NewSender(lg *zap.Logger, storage storage.Storage, dir string) Sender {
	return &sender{
		lg:      lg,
		storage: storage,
		dir:     dir,
	}
}

func (s *sender) Upload(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error) {
	if req.Header == nil || req.File == nil {
		return nil, ErrEmptyFile
	}

	contentType := req.Header.Header.Get("Content-Type")
	if _, ok := AcceptedImagesContentTypes[contentType]; !ok {
		return nil, ErrImageFileOnly
	}

	userId, err := strconv.ParseInt(auth.UserId(ctx), 10, 64)
	if err != nil {
		return nil, err
	}
	id, err := s.storage.InsertFile(&query.InsertFileParams{
		UserID:      userId,
		Filename:    req.Header.Filename,
		ContentType: contentType,
		Size:        int32(req.Header.Size),
	})
	if err != nil {
		return nil, err
	}

	// in order to make sure the file path is unique, filename will be file_id
	path := fmt.Sprintf("%s/%d", s.dir, id)
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = io.Copy(f, req.File)
	if err != nil {
		return nil, err
	}

	resp := &UploadFileResponse{
		FileId: id,
	}
	return resp, nil
}
