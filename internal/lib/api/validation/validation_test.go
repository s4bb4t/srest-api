package validation

import (
	"testing"

	"github.com/sabbatD/srest-api/internal/lib/userConfig"
)

func TestValidateStruct(t *testing.T) {
	InitValidator()

	tests := []struct {
		name    string
		args    userConfig.User
		wantErr bool
	}{
		{
			name: "normal",
			args: userConfig.User{
				Login:    "login",
				Username: "username",
				Password: "password",
				Email:    "email@example.com",
			},
		},
		{
			name: "empty",
			args: userConfig.User{
				Login:    "",
				Username: "",
				Password: "",
				Email:    "",
			},
		},
		{
			name: "spaces",
			args: userConfig.User{
				Login:    "      ",
				Username: "      ",
				Password: "      ",
				Email:    "      ",
			},
		},
		{
			name: "to much",
			args: userConfig.User{
				Login:    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Username: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Password: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Email:    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()-_=+[{]};:',<.>/?",
			},
		},
		{
			name: "symbs",
			args: userConfig.User{
				Login:    "`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Username: "`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Password: "`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Email:    "`~!@#$%^&*()-_=+[{]};:',<.>/?",
			},
		},
		{
			name: "hitry",
			args: userConfig.User{
				Login:    "Login",
				Username: "Username",
				Password: "`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Email:    "`~!@#$%^&*",
			},
		},
		{
			name: "one-empty",
			args: userConfig.User{
				Login:    "Login",
				Username: "",
				Password: "`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Email:    "`~!@#$%^&*",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateStruct(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
