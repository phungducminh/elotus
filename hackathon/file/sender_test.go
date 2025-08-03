package file

import (
	"context"
	"mime/multipart"
	"net/textproto"
	"testing"

	"go.uber.org/zap/zaptest"
)

func TestUpload(t *testing.T) {
	tests := []struct {
		desc        string
		contentType string
		expect      error
	}{
		{
			desc:        "image/jpeg",
			contentType: "image/jpeg",
			expect:      nil,
		},
		{
			desc:        "image/avif",
			contentType: "image/avif",
			expect:      nil,
		},
		{
			desc:        "application/octet-stream",
			contentType: "application/octet-stream",
			expect:      ErrImageFileOnly,
		},
	}

	for _, tt := range tests {
		sender := NewSender(zaptest.NewLogger(t))
		req := &UploadFileRequest{
			File: nil,
			Header: &multipart.FileHeader{
				Filename: "file.txt",
				Header: textproto.MIMEHeader{
					"Content-Type": []string{tt.contentType},
				},
				Size: 1000,
			},
		}
		_, err := sender.Upload(context.Background(), req)
		if err != tt.expect {
			t.Errorf("Upload(), expect=%v, actual=%v", tt.expect, err)
		}
	}
}
