package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"elotus.com/hackathon/auth"
	"elotus.com/hackathon/file"
	"elotus.com/hackathon/server"
	"elotus.com/hackathon/storage"
)

func TestRegisterThenLoginThenUpload(t *testing.T) {
	opts := make([]server.Option, 0)
	opts = append(opts, server.WithUploadFileDir(t.TempDir()))
	s, err := server.NewServer(opts...)
	s.Storage = storage.NewRecorder()
	if err != nil {
		t.Fatal(err)
	}

	registerHandler := auth.NewRegisterHandler(s)
	http.Handle("/api/auth/register", registerHandler)
	loginHandler := auth.NewLoginHandler(s)
	http.Handle("/api/auth/login", loginHandler)
	fileHandler := file.NewFileHandler(s)
	http.Handle("/api/file/upload", auth.AuthMiddleware(s, fileHandler))
	// TODO: @replace this with our serveMux
	ss := httptest.NewServer(http.DefaultServeMux)
	defer ss.Close()

	if err := register(ss); err != nil {
		t.Fatal(err)
	}

	resp, err := login(ss)
	if err != nil {
		t.Fatal(err)
	}

	p, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	var login auth.LoginResponse
	json.Unmarshal(p, &login)


	buffer := &bytes.Buffer{}
	multipartW := multipart.NewWriter(buffer)
	w, err := multipartW.CreateFormFile("data", "/Users/phungducminh/Pictures/pexels-pixabay-358532.jpg")
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("/Users/phungducminh/Pictures/pexels-pixabay-358532.jpg")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteTo(w)
	if err != nil {
		t.Fatal(err)
	}

	multipartW.Close()

	req, err := http.NewRequest("POST", ss.URL+"/api/file/upload", buffer)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", login.AccessToken))
	req.Header.Add("Content-Type", multipartW.FormDataContentType())
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	p, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("upload response: %s\n", p)
}

func login(ss *httptest.Server) (*http.Response, error) {
	p, err := json.Marshal(&auth.LoginRequest{
		Username: "john",
		Password: "123456",
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(ss.URL+"/api/auth/login", "application/json", bytes.NewReader(p))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func register(s *httptest.Server) error {
	p, err := json.Marshal(&auth.RegisterRequest{
		Username: "john",
		Password: "123456",
	})
	resp, err := http.Post(s.URL+"/api/auth/register", "application/json", bytes.NewReader(p))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
