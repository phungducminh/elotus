package file

import "context"

type Sender interface {
	Upload(context.Context, *UploadFileRequest) (*UploadFileResponse, error)
}

type sender struct {
}

func NewSender() Sender {
	return &sender{}
}

func (s *sender) Upload(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error) {
	return &UploadFileResponse{}, nil
}
