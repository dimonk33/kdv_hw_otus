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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "1",
				Name:   "test",
				Age:    1,
				Email:  "test@test.ru",
				Role:   "test",
				Phones: []string{"11111111111"},
				meta:   json.RawMessage(""),
			},
			ValidationErrors{
				ValidationError{"ID", ErrValidateLenCond},
				ValidationError{"Age", ErrValidateMinCond},
				ValidationError{"Role", ErrValidateArrayCond},
			},
		},
		{
			1,
			ErrNotStruct,
		},
		{
			App{
				Version: "12345",
			},
			ValidationErrors{},
		},
		{
			Token{
				Header:    []byte{1, 2},
				Payload:   []byte{3, 4},
				Signature: []byte{5, 6},
			},
			ValidationErrors{},
		},
		{
			Response{
				Code: 200,
				Body: "",
			},
			ValidationErrors{},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			var validationErrors, expectedErrors ValidationErrors
			if errors.As(err, &validationErrors) && errors.As(tt.expectedErr, &expectedErrors) {
				for i, e := range validationErrors {
					res := errors.Is(e.Err, (expectedErrors)[i].Err)
					require.True(t, res)
				}
			} else {
				require.Equal(t, tt.expectedErr, err)
			}
		})
	}

	t.Run("int validate", func(t *testing.T) {
		conditionStr := "in:1,10,100|min:0|max:99"
		err := validateInt(1, conditionStr)
		require.Nil(t, err)
		err = validateInt(-1, conditionStr)
		require.NotNil(t, err)
		err = validateInt(100, conditionStr)
		require.NotNil(t, err)
	})

	t.Run("string validate 1", func(t *testing.T) {
		conditionStr := "in:test1,test_2,test_55|len:6|regexp:^[[:alpha:]]{4}_\\d{1,}$"
		err := validateString("test_2", conditionStr)
		require.Nil(t, err)
		err = validateString("test1", conditionStr)
		require.NotNil(t, err)
		err = validateString("test_55", conditionStr)
		require.NotNil(t, err)
	})
}
