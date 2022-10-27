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
		ID     string `json:"id" validate:"len:10"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		// meta   json.RawMessage
	}
)

var tests = []struct {
	in          interface{}
	expectedErr ValidationErrors
}{
	{
		in: User{
			ID:     "10",
			Name:   "Andrey",
			Age:    21,
			Email:  "shabandrew@mail.ru",
			Role:   "admin",
			Phones: []string{"89045267897", "12345678"},
		},
		expectedErr: nil,
	},
	{
		in: User{
			ID:     "11",
			Name:   "Коля",
			Age:    50,
			Email:  "kolan@gmail.com",
			Role:   "admin",
			Phones: []string{"89045267897", "12345678", "12345678"},
		},
		expectedErr: nil,
	},
	{
		in: User{
			ID:    "46mira46",
			Name:  "Mira",
			Age:   21,
			Email: "marathebest@mail.ru",
			Role:  "admin",
		},
		expectedErr: nil,
	},
	{
		in: User{
			ID:    "10",
			Name:  "Andrey",
			Age:   21,
			Email: "shabandrew@mail.ru",
			Role:  "stuff",
		},
		expectedErr: nil,
	},
	{
		in: User{
			ID:    "10",
			Name:  "Andrey",
			Age:   21,
			Email: "shabandrew@mail.ru",
			Role:  "stuff",
		},
		expectedErr: nil,
	},
	{
		in: User{
			ID:    "012345678910",
			Name:  "Andrey",
			Age:   21,
			Email: "shabandrew@mail.ru",
			Role:  "stuff",
		},
		expectedErr: []ValidationError{{
			Field: "ID",
			Err:   ErrStringLen,
		}},
	},
	{
		in: User{
			ID:     "10",
			Name:   "Andrey",
			Age:    10,
			Email:  "shabandrew@mail.ru",
			Role:   "admin",
			Phones: []string{"89045267897", "12345678"},
		},
		expectedErr: []ValidationError{{
			Field: "Age",
			Err:   ErrIntMin,
		}},
	},
	{
		in: User{
			ID:     "10",
			Name:   "Andrey",
			Age:    100,
			Email:  "shabandrew@mail.ru",
			Role:   "admin",
			Phones: []string{"89045267897", "12345678"},
		},
		expectedErr: []ValidationError{{
			Field: "Age",
			Err:   ErrIntMax,
		}},
	},
	{
		in: User{
			ID:     "10",
			Name:   "Andrey",
			Age:    21,
			Email:  "@mail.ru",
			Role:   "admin",
			Phones: []string{"89045267897", "12345678"},
		},
		expectedErr: []ValidationError{{
			Field: "Email",
			Err:   ErrStringRegexp,
		}},
	},
	{
		in: User{
			ID:     "10",
			Name:   "Andrey",
			Age:    21,
			Email:  "shabandrew@mail.ru",
			Role:   "student",
			Phones: []string{"89045267897", "12345678"},
		},
		expectedErr: []ValidationError{{
			Field: "Role",
			Err:   ErrStringIn,
		}},
	},
}

func TestValidate(t *testing.T) {
	for i, test := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			test := test
			t.Parallel()

			require.Equal(t, test.expectedErr, Validate(test.in))
		})
	}
}
