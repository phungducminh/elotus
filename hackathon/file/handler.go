package file

import (
	"encoding/json"
	"net/http"

	"elotus.com/hackathon/server"
)

const MaxUploadSize = 8 << 20

type UploadFileRequest struct {
	FileName string
	Data     []byte
}

type UploadFileResponse struct {
}

type FileHandler struct {
	server *server.Server
	sender Sender
}

func NewFileHandler(s *server.Server) *FileHandler {
	return &FileHandler{
		sender: NewSender(),
	}
}

func (h *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	err := r.ParseMultipartForm(MaxUploadSize)
	if err != nil {
		// TODO: replace checking file size logic with better one
		w.Header().Set("Content-Type", "application/json")
		resp := &Response{
			Error: ErrorResponse{
				Code:    "FILE_TOO_LARGE",
				Message: "file exceeds 8MB",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var req UploadFileRequest

	resp, err := h.sender.Upload(r.Context(), &req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// TODO: handle duplication
type Response struct {
	Error ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
