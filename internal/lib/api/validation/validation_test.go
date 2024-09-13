package validation

import (
	"testing"

	todoconfig "github.com/sabbatD/srest-api/internal/lib/todoConfig"
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
			wantErr: true,
		},
		{
			name: "spaces",
			args: userConfig.User{
				Login:    "      ",
				Username: "      ",
				Password: "      ",
				Email:    "      ",
			},
			wantErr: true,
		},
		{
			name: "to much",
			args: userConfig.User{
				Login:    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Username: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Password: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Email:    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()-_=+[{]};:',<.>/?",
			},
			wantErr: true,
		},
		{
			name: "symbs",
			args: userConfig.User{
				Login:    "`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Username: "`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Password: "`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Email:    "`~!@#$%^&*()-_=+[{]};:',<.>/?",
			},
			wantErr: true,
		},
		{
			name: "hitry",
			args: userConfig.User{
				Login:    "Login",
				Username: "Username",
				Password: "`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Email:    "`~!@#$%^&*",
			},
			wantErr: true,
		},
		{
			name: "one-empty",
			args: userConfig.User{
				Login:    "Login",
				Username: "",
				Password: "`~!@#$%^&*()-_=+[{]};:',<.>/?",
				Email:    "`~!@#$%^&*",
			},
			wantErr: true,
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

func TestValidateStruct2(t *testing.T) {
	InitValidator()

	tests := []struct {
		name    string
		args    todoconfig.TodoRequest
		wantErr bool
	}{
		{
			name: "normal",
			args: todoconfig.TodoRequest{
				Title:  "todo",
				IsDone: "true",
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: todoconfig.TodoRequest{
				Title:  "",
				IsDone: "false",
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: todoconfig.TodoRequest{
				Title:  "todo",
				IsDone: "",
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: todoconfig.TodoRequest{
				Title:  "",
				IsDone: "",
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: todoconfig.TodoRequest{
				Title:  "",
				IsDone: "abc",
			},
			wantErr: true,
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
