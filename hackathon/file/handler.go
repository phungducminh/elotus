package file

import (
	"encoding/json"
	"net/http"

	. "elotus.com/hackathon/pkg/logutil/httputil"
	"elotus.com/hackathon/server"
	"go.uber.org/zap"
)

const MaxUploadSize = 8 << 20

type UploadFileRequest struct {
	FileName string
	Data     []byte
}

type UploadFileResponse struct {
}

type FileHandler struct {
	sender Sender
	lg     *zap.Logger
}

func NewFileHandler(s *server.Server) *FileHandler {
	return &FileHandler{
		sender: NewSender(),
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
		h.lg.Warn("file too large")
		ResponseBadRequest(w, "FILE_TOO_LARGE", "file exceeds 8MB")
		return
	}

	var req UploadFileRequest

	resp, err := h.sender.Upload(r.Context(), &req)
	if err != nil {
		ResponseInternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
