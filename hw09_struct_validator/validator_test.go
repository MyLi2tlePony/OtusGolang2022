package hw09structvalidator

import (
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
		// meta   json.RawMessage
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

var tests = []struct {
	in          interface{}
	expectedErr error
}{
	{
		in: &User{
			ID:    "10",
			Name:  "Andrey",
			Age:   21,
			Email: "shabandrew@mail.ru",
			Role:  "main role",
		},
		expectedErr: nil,
	},
	{
		in: &User{
			ID:    "11",
			Name:  "Коля",
			Age:   50,
			Email: "kolan@gmail.com",
			Role:  "new role",
		},
		expectedErr: nil,
	},
	{
		in: &User{
			ID:    "46mira46",
			Name:  "Mira",
			Age:   21,
			Email: "marathebest@mail.ru",
			Role:  "no role",
		},
		expectedErr: nil,
	},
	{
		in: &User{
			ID:    "10",
			Name:  "Andrey",
			Age:   21,
			Email: "shabandrew@mail.ru",
			Role:  "kjbvaskjdvckjladklvn role",
		},
		expectedErr: nil,
	},
	{
		in: &User{
			ID:    "10",
			Name:  "Andrey",
			Age:   21,
			Email: "shabandrew@mail.ru",
			Role:  "some role",
		},
		expectedErr: nil,
	},
}

func TestValidate(t *testing.T) {
	for i, test := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			test := test
			t.Parallel()

			require.Equal(t, nil, Validate(test))
		})
	}
}
