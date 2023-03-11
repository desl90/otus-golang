package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Task struct {
		ID    int    `validate:"min:1|max:500"`
		title string `validate:"min:6|max:255"`
		grade int
	}

	Month struct {
		Summer int `validate:"min:3|max:5|in:3,4,5"`
	}

	Bad struct {
		ErrTagValue int `validate:"min:1|max:s5|in:3,5"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			777,
			ErrWrongType,
		},
		{
			Task{ID: 300, title: "Test title", grade: -1},
			nil,
		},
		{
			Task{ID: -1, title: "Test task", grade: -1},
			errors.New("Id: wrong value less min"),
		},
		{
			Task{ID: 1000, title: "Test task", grade: -1},
			errors.New("Id: wrong value more max"),
		},
		{
			Month{Summer: 4},
			nil,
		},
		{
			Bad{ErrTagValue: 3},
			errors.New("strconv.Atoi: parsing \"s5\": invalid syntax"),
		},
		{
			App{Version: "ios15"},
			nil,
		},
		{
			App{Version: "linux6"},
			errors.New("Version: wrong value length"),
		},
		{
			Response{Code: 200, Body: "ok"},
			nil,
		},
		{
			Response{Code: 301, Body: "Moved Permanently"},
			errors.New("Code: wrong value not in"),
		},
		{
			Token{Header: []byte("a"), Payload: []byte("b"), Signature: []byte("c")},
			nil,
		},
		{
			User{
				ID:     "111111111111111111111111111111111111",
				Name:   "Jon Doe",
				Age:    40,
				Email:  "test@test.ru",
				Role:   "admin",
				Phones: []string{"79990001122"},
				meta:   json.RawMessage{},
			},
			nil,
		},
		{
			User{
				ID:     "1",
				Name:   "Jon Doe",
				Age:    40,
				Email:  "test@test.ru",
				Role:   "admin",
				Phones: []string{"79990001122"},
				meta:   json.RawMessage{},
			},
			errors.New("ID: wrong value length"),
		},
		{
			User{
				ID:     "111111111111111111111111111111111111",
				Name:   "Jon Doe",
				Age:    40,
				Email:  "~test@test.ru",
				Role:   "admin",
				Phones: []string{"79990001122"},
				meta:   json.RawMessage{},
			},
			errors.New("Email: wrong value regexp"),
		},
		{
			User{
				ID:     "111111111111111111111111111111111111",
				Name:   "Jon Doe",
				Age:    40,
				Email:  "test@test.ru",
				Role:   "user",
				Phones: []string{"79990001122"},
				meta:   json.RawMessage{},
			},
			errors.New("Role: wrong value not in"),
		},
		{
			User{
				ID:     "111111111111111111111111111111111111",
				Name:   "Jon Doe",
				Age:    40,
				Email:  "test@test.ru",
				Role:   "admin",
				Phones: []string{"79990001122", "79990001133", "79990001144"},
				meta:   json.RawMessage{},
			},
			nil,
		},
		{
			User{
				ID:     "111111111111111111111111111111111111",
				Name:   "Jon Doe",
				Age:    40,
				Email:  "test@test.ru",
				Role:   "admin",
				Phones: []string{"799900011", "79990001133", "79990001144"},
				meta:   json.RawMessage{},
			},
			errors.New("Phones: wrong value length value 799900011"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			result := Validate(tt.in)

			if tt.expectedErr == nil {
				require.Equal(t, tt.expectedErr, result)
			} else {
				require.Equal(t, tt.expectedErr.Error(), result.Error())
			}

			_ = tt
		})
	}
}
