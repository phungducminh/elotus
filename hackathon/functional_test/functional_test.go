package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

	registerReq := &auth.RegisterRequest{
		Username: "john",
		Password: "123456",
	}
	w := &bytes.Buffer{}
	encoder := json.NewEncoder(w)
	encoder.Encode(registerReq)
	resp, err := http.Post(ss.URL+"/api/auth/register", "application/json", w)
	if err != nil {
		t.Fatal(err)
	}
	w.Reset()

	registerResp, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("register response: %s\n", registerResp)

	loginReq := &auth.LoginRequest{
		Username: "john",
		Password: "123456",
	}
	encoder.Encode(loginReq)
	resp, err = http.Post(ss.URL+"/api/auth/login", "application/json", w)
	if err != nil {
		t.Fatal(err)
	}
	w.Reset()

	loginResp, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var loginRespVal auth.LoginResponse
	decoder := json.NewDecoder(bytes.NewBuffer(loginResp))
	decoder.Decode(&loginRespVal)
	fmt.Printf("login response: %v\n", loginRespVal)

	uploadReq, err := http.NewRequest("POST", ss.URL+"/api/file/upload", nil)
	if err != nil {
		t.Fatal(err)
	}
	uploadReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", loginRespVal.AccessToken))
	query := uploadReq.URL.Query()
	// TODO: @replace with other picture
	query.Add("data", "")
	uploadReq.URL.RawQuery = query.Encode()
	fmt.Printf("request: %v\n", uploadReq)
	resp, err = http.DefaultClient.Do(uploadReq)
	if err != nil {
		t.Fatal(err)
	}

	uploadResp, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("upload response: %s\n", uploadResp)
}
