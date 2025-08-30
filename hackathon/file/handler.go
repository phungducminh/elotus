package file

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	. "elotus.com/hackathon/pkg/logutil/httputil"
	"elotus.com/hackathon/server"
	"go.uber.org/zap"
)

const MaxUploadSize = 8 << 20

type UploadFileRequest struct {
	File   io.Reader
	Header *multipart.FileHeader
}

type UploadFileResponse struct {
	FileId int64
}

type FileHandler struct {
	sender Sender
	lg     *zap.Logger
}

func NewFileHandler(s *server.Server) *FileHandler {
	return &FileHandler{
		sender: NewSender(s.Logger, s.Storage, s.Cfg.UploadFileDir),
		lg:     s.Logger,
	}
}

func (h *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ResponseMethodNotAllowed(w)
		return
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	err := r.ParseMultipartForm(MaxUploadSize)
	if err != nil {
		h.lg.Warn("file too large", zap.Error(err))
		ResponseBadRequest(w, "FILE_TOO_LARGE", "file exceeds 8MB")
		return
	}
	file, header, err := r.FormFile("data")
	if err != nil {
		ResponseBadRequest(w, "INVALID_FORM", "expect a form with field named 'data'")
		return
	}
	defer file.Close()

	h.lg.Info("file info", zap.Any("header", header), zap.Any("file", file))
	req := &UploadFileRequest{
		File:   file,
		Header: header,
	}
	resp, err := h.sender.Upload(r.Context(), req)
	if err == ErrImageFileOnly {
		ResponseBadRequest(w, "IMAGE_ONLY", "image only")
		return
	} else if err == ErrEmptyFile {
		ResponseBadRequest(w, "EMPTY_FILE", "file must not be empty")
		return
	}
	if err != nil {
		h.lg.Warn("failed to upload file", zap.Error(err))
		ResponseInternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
