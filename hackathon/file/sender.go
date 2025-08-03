package file

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

var (
	ErrInvalidFileHeader = fmt.Errorf("invalid file header")
	ErrImageFileOnly     = fmt.Errorf("image file only")
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
	lg *zap.Logger
}

func NewSender(lg *zap.Logger) Sender {
	return &sender{
		lg: lg,
	}
}

func (s *sender) Upload(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error) {
	if req.Header == nil {
		return nil, ErrInvalidFileHeader
	}

	if _, ok := AcceptedImagesContentTypes[req.Header.Header.Get("Content-Type")]; !ok {
		return nil, ErrImageFileOnly
	}

	return &UploadFileResponse{}, nil
}
