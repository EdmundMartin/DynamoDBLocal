package expressions

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestConsumeParameter(t *testing.T) {

	testCases := []struct {
		Name          string
		Value         string
		WantErr       bool
		ExpectedToken *Token
	}{
		{
			Name:    "A valid parameter",
			Value:   ":category1",
			WantErr: false,
			ExpectedToken: &Token{
				Literal:   ":category1",
				TokenType: Parameter,
			},
		},
		{
			Name:    "A invalid parameter",
			Value:   ":1dog",
			WantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tok, _, err := ConsumeParameter(strings.Split(tc.Value, ""), 0)

			if !tc.WantErr {
				assert.NoError(t, err)
				assert.Equal(t, tc.ExpectedToken, tok)
			}
			if tc.WantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestConsumeString(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         string
		WantErr       bool
		ExpectedToken *Token
	}{
		{
			Name:    "A valid identifier",
			Value:   "price",
			WantErr: false,
			ExpectedToken: &Token{
				Literal:        "price",
				TokenType:      Identifier,
				ContainedValue: "",
			},
		},
		{
			Name:    "A valid Attribute Exists",
			Value:   "attribute_exists(price)",
			WantErr: false,
			ExpectedToken: &Token{
				Literal:        "attribute_exists",
				TokenType:      AttributeExists,
				ContainedValue: "price",
			},
		},
		{
			Name:    "A valid not attribute Exists",
			Value:   "attribute_not_exists(price)",
			WantErr: false,
			ExpectedToken: &Token{
				Literal:        "attribute_not_exists",
				TokenType:      AttributeNotExists,
				ContainedValue: "price",
			},
		},
		{
			Name:    "A valid size token",
			Value:   "size(price)",
			WantErr: false,
			ExpectedToken: &Token{
				Literal:        "size",
				TokenType:      Size,
				ContainedValue: "price",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			tok, _, err := ConsumeString(strings.Split(tc.Value, ""), 0)

			if !tc.WantErr {
				assert.NoError(t, err)
				assert.Equal(t, tc.ExpectedToken, tok)
			}
		})
	}
}
