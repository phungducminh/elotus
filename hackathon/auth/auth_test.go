package auth

import (
	"testing"
	"time"

	"elotus.com/hackathon/storage"
)

func TestRegister(t *testing.T) {
	recorder := storage.NewRecorder()
	auth := NewAuth(recorder, []byte("SECRET_KEY"), 60)
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

	_, err = auth.Register(&RegisterRequest{
		Username: "john",
		Password: "B@123P",
	})
	if err == nil {
		t.Fatalf("register duplicate username, expect register failed, username=%s, password=%s", req.Username, req.Password)
	}
}

func TestLogin(t *testing.T) {
	recorder := storage.NewRecorder()
	auth := NewAuth(recorder, []byte("SECRET_KEY"), 60)
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

func TestVerifyValidAndExpiredToken(t *testing.T) {
	recorder := storage.NewRecorder()
	auth := NewAuth(recorder, []byte("SECRET_KEY"), 2)
	registerResp, err := auth.Register(&RegisterRequest{
		Username: "john",
		Password: "A@123P",
	})
	if err != nil {
		t.Fatal(err)
	}

	loginResp, err := auth.Login(&LoginRequest{
		Username: "john",
		Password: "A@123P",
	})
	if err != nil {
		t.Fatal(err)
	}

	verifyResp, err := auth.Verify(&VerifyRequest{
		AccessToken: loginResp.AccessToken,
	})
	if err != nil {
		t.Fatal(err)
	}

	if verifyResp.UserId != registerResp.UserId {
		t.Fatalf("verify valid token, expect=%s, actual=%s", registerResp.UserId, verifyResp.UserId)
	}

	<-time.Tick(time.Second * 3)
	_, err = auth.Verify(&VerifyRequest{
		AccessToken: loginResp.AccessToken,
	})
	if err == nil {
		t.Fatal("expect token is expired")
	}
}

func TestVerifyInvalidToken(t *testing.T) {
	recorder := storage.NewRecorder()
	auth := NewAuth(recorder, []byte("SECRET_KEY"), 60)

	_, err := auth.Register(&RegisterRequest{
		Username: "john",
		Password: "A@123P",
	})
	if err != nil {
		t.Fatal(err)
	}

	loginResp, err := auth.Login(&LoginRequest{
		Username: "john",
		Password: "A@123P",
	})
	if err != nil {
		t.Fatal(err)
	}

	invalidToken := loginResp.AccessToken[:len(loginResp.AccessToken)-1]
	_, err = auth.Verify(&VerifyRequest{
		AccessToken: invalidToken,
	})
	if err == nil {
		t.Fatal("expect invalid token")
	}

	invalidAuth := NewAuth(recorder, []byte("ANOTHER_SECRET_KEY"), 60)
	_, err = invalidAuth.Verify(&VerifyRequest{
		AccessToken: loginResp.AccessToken,
	})
	if err == nil {
		t.Fatal("expect invalid token")
	}
}
