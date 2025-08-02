package auth

import (
	"testing"

	"elotus.com/hackathon/storage"
)

func TestRegister(t *testing.T) {
	recorder := storage.NewRecorder()
	auth := NewAuth(recorder, []byte("SERCRET_KEY"))
	req := &RegisterRequest{
		Username: "john",
		Password: "A@123P",
	}
	resp, err := auth.Register(req)
	if err != nil {
		t.Fatalf("expect register successfully, username=%s, password=%s", req.Username, req.Password)
	}
	
	user, err := recorder.GetUserByUserName("john")
	if err != nil {
		t.Fatalf("expect found registered user, username=%s, password=%s, userId=%s",
			req.Username, req.Password, resp.UserId)
	}

	if user.HashedPassword == req.Password {
		t.Fatal("expect password must be hashed")
	}
}

func TestLogin(t *testing.T) {
	recorder := storage.NewRecorder()
	auth := NewAuth(recorder, []byte("SERCRET_KEY"))
	_, err := auth.Register(&RegisterRequest{
		Username: "john",
		Password: "A@123P",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = auth.Login(&LoginRequest{
		Username: "john",
		Password: "A@123P",
	})
	if err != nil {
		t.Fatalf("login with valid credential, expect not error, but error=%v", err)
	}

	_, err = auth.Login(&LoginRequest{
		Username: "john0",
		Password: "A@123P",
	})
	if err != ErrInvalidCredentials {
		t.Fatalf("login with invalid username, expect invalid credentals error, error=%v", err)
	}

	_, err = auth.Login(&LoginRequest{
		Username: "john",
		Password: "a@123P",
	})
	if err != ErrInvalidCredentials {
		t.Fatalf("login with invalid password, expect invalid credentals error, error=%v", err)
	}

	_, err = auth.Login(&LoginRequest{
		Username: "mayer",
		Password: "Ab@123B",
	})
	if err == nil {
		t.Fatal("expect a non-registered user won't be able to login")
	}
}
