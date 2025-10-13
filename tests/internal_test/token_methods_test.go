package internal_test

import (
	"jwt_clean/internal/auth"
	"jwt_clean/internal/service"
	"testing"
	"time"
)

var testSecret = []byte("test-secret-key")

func TestNewAuthService(t *testing.T) {
	s := service.NewAuthService(testSecret)
	if s == nil {
		t.Fatal("Expected service to be created, got nil")
	}
}

func TestGenerateToken(t *testing.T) {
	s := service.NewAuthService(testSecret)
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "valid username",
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "empty username",
			username: "",
			wantErr:  false,
		},
		{
			name:     "Special characters in username",
			username: "user@example.com",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := s.GenerateToken(tt.username)

			if (err!=nil) != tt.wantErr{
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && token == "" {
				t.Error("Expected a token, got empty string")
			}
		})
	}
}

func TestGenerateAccessToken(t *testing.T) {
	s:=service.NewAuthService(testSecret)
	token, err := s.GenerateAccessToken("testuser")
	if err!=nil{
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	if token==""{
		t.Error("expected token to be generated, got empty string")
	}

	claims, err:=s.ParseToken(token)
	if err!=nil{
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims.Subject!="testuser"{
		t.Errorf("expected subject 'testuser', got %s", claims.Subject)
	}

	if claims.TokenType!=auth.AccessTokenType{
		t.Errorf("expected token type %v, got %v", auth.AccessTokenType, claims.TokenType)
	}

	expectedExpiry:=time.Now().Add(15*time.Minute)
	acutalExpiry:=claims.ExpiresAt.Time
	diff:=acutalExpiry.Sub(expectedExpiry).Abs()

	if diff>time.Second*5{
		t.Errorf("expected expiry around %v, got %v", expectedExpiry, acutalExpiry)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	s:=service.NewAuthService(testSecret)
	token, err := s.GenerateRefreshToken("testuser")
	if err!=nil{
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	if token==""{
		t.Error("expected token to be generated, got empty string")
	}

	claims, err:=s.ParseToken(token)
	if err!=nil{
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims.Subject!="testuser"{
		t.Errorf("expected subject 'testuser', got %s", claims.Subject)
	}

	if claims.TokenType!=auth.RefreshTokenType{
		t.Errorf("expected token type %v, got %v", auth.RefreshTokenType, claims.TokenType)
	}

	expectedExpiry:=time.Now().Add(7*24*time.Hour)
	acutalExpiry:=claims.ExpiresAt.Time
	diff:=acutalExpiry.Sub(expectedExpiry).Abs()

	if diff>time.Second*5{
		t.Errorf("expected expiry around %v, got %v", expectedExpiry, acutalExpiry)
	}
}

func TestGenerateTokenPair(t *testing.T) {
	s:=service.NewAuthService(testSecret)

	tests:=[]struct {
		name string
		username string
		wantErr bool
	}{
		{
			name: "valid username",
			username: "testuser",
			wantErr: false,
		},
		{
			name: "empty username",
			username: "",
			wantErr: false,
		},
	}

	for _, tt:=range tests{
		t.Run(tt.name, func(t *testing.T) {
			pair, err:=s.GenerateTokenPair(tt.username)
			if (err!=nil) != tt.wantErr{
				t.Errorf("GenerateTokenPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr{
				return
			}

			if pair.AccessToken==""{
				t.Error("expected access token to be generated, got empty string")
			}

			if pair.RefreshToken==""{
				t.Error("expected refresh token to be generated, got empty string")
			}

			if pair.ExpiresIn!=900{
				t.Errorf("expected expires in 900, got %d", pair.ExpiresIn)
			}

			if pair.TokenType!="Bearer"{
				t.Errorf("expected token type 'Bearer', got %s", pair.TokenType)
			}

			accessClaims, err:=s.ParseToken(pair.RefreshToken)
			if err!=nil{
				t.Errorf("failed to parse access token: %v", err)
			} else if accessClaims.TokenType!=auth.AccessTokenType{
				t.Errorf("expected access token type %v, got %v", auth.AccessTokenType, accessClaims.TokenType)
			}

			refreshClaims, err:=s.ParseToken(pair.RefreshToken)
			if err!=nil{
				t.Errorf("failed to parse refresh token: %v", err)
			} else if refreshClaims.TokenType!=auth.RefreshTokenType{
				t.Errorf("expected refresh token type %v, got %v", auth.RefreshTokenType, refreshClaims.TokenType)
			}
		})
	}
}

