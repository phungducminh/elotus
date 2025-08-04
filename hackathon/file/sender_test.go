package file

import (
	"context"
	"mime/multipart"
	"net/textproto"
	"strings"
	"testing"

	"elotus.com/hackathon/auth"
	"elotus.com/hackathon/storage"
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
		dir := t.TempDir()
		sender := NewSender(zaptest.NewLogger(t), storage.NewRecorder(), dir)
		req := &UploadFileRequest{
			File: strings.NewReader("Hello World"),
			Header: &multipart.FileHeader{
				Filename: "file.txt",
				Header: textproto.MIMEHeader{
					"Content-Type": []string{tt.contentType},
				},
				Size: 1000,
			},
		}
		ctx := auth.WithUserId(context.Background(), "123")
		_, err := sender.Upload(ctx, req)
		if err != tt.expect {
			t.Errorf("Upload(), expect=%v, actual=%v", tt.expect, err)
		}
	}
}
