package middleware

import (
	"github.com/gold-kou/go-housework/app/common"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var sharedTestToken string

func TestGenerateToken(t *testing.T) {
	type args struct {
		userName string
	}
	tests := []struct {
		name        string
		args        args
		environment string
		want        string
		wantErr     bool
	}{
		{
			name:        "success",
			args:        args{userName: common.TestUserName},
			environment: common.TestSecretKey,
			want:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			// set env
			tmp := common.SetTestEnv("JWT_SECRET_KEY", tt.environment)
			defer tmp()

			// test target
			got, err := GenerateToken(tt.args.userName)
			sharedTestToken = got

			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				// just checking HEADER
				assert.Equal(tt.want, strings.Split(got, ".")[0])
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	type args struct {
		tokenString string
	}
	tests := []struct {
		name        string
		args        args
		environment string
		want        *Auth
		wantErr     bool
		watnErrMsg  string
	}{
		{
			name:        "success",
			args:        args{tokenString: sharedTestToken},
			environment: common.TestSecretKey,
			want:        &Auth{UserName: common.TestUserName},
			wantErr:     false,
		},
		{
			name:        "fail(expired)",
			args:        args{tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM3MjIzNjIsImlhdCI6IjIwMjAtMDMtMDhUMTE6NTI6NDIuMjIxMjY2NCswOTowMCIsIm5hbWUiOiJ0ZXN0In0.FYMJmXo17aUhTpdaLifMovDQ0BiKSq8LnssLwxFvshI"},
			environment: common.TestSecretKey,
			want:        &Auth{UserName: common.TestUserName},
			wantErr:     true,
			watnErrMsg:  "token is expired",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set env
			tmp := common.SetTestEnv("JWT_SECRET_KEY", tt.environment)
			defer tmp()

			assert := assert.New(t)

			// test target
			got, err := VerifyToken(tt.args.tokenString)

			// assert
			if tt.wantErr {
				assert.Error(err)
				assert.EqualError(err, tt.watnErrMsg)
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}
