package expressions

import (
	"bytes"
	"errors"
	"unicode"
)

type TokenType uint8

const (
	Unknown TokenType = iota
	Identifier
	Parameter
	Between
	GreaterThan
	AttributeExists
	AttributeNotExists
	Size
)

type Token struct {
	Literal        string
	TokenType      TokenType
	ContainedValue string
}

var (
	validKeywords = map[string]TokenType{
		"between":              Between,
		"attribute_exists":     AttributeExists,
		"attribute_not_exists": AttributeNotExists,
		"size":                 Size,
	}
)

func ConsumeString(chars []string, idx int) (*Token, int, error) {
	literal := bytes.Buffer{}
	seenStartChar := false
	for {
		if idx == len(chars) {
			break
		}
		charRune := []rune(chars[idx])[0]
		// TODO Check if Identifiers can start with underscore
		if unicode.IsDigit(charRune) || charRune == '_' && !seenStartChar {
			return nil, 0, errors.New("expected a char as initial part of token")
		}

		if unicode.IsLetter(charRune) {
			seenStartChar = true
			literal.WriteString(chars[idx])
		} else if unicode.IsDigit(charRune) || charRune == '_' {
			literal.WriteString(chars[idx])
		} else if charRune == '(' {

			// TODO - Logic for contains, size, type
			containedValue, err := consumeContainedValue(chars, idx+1)
			if err != nil {
				return nil, 0, err
			}
			tokenLit := literal.String()
			tokenType, ok := validKeywords[tokenLit]
			if !ok {
				return nil, 0, errors.New("invalid token type")
			}
			length := len(tokenLit) + len(containedValue) + 2
			return &Token{TokenType: tokenType, Literal: tokenLit, ContainedValue: containedValue}, length, nil
		} else {
			break
		}
		idx++
	}
	tokenLit := literal.String()
	return &Token{TokenType: Identifier, Literal: tokenLit, ContainedValue: ""}, len(tokenLit), nil
}

func consumeContainedValue(chars []string, idx int) (string, error) {
	result := bytes.Buffer{}
	var lastRune rune
	for {
		if idx == len(chars) {
			break
		}

		charRune := []rune(chars[idx])[0]
		lastRune = charRune
		if charRune == ')' {
			break
		}
		// TODO - ':' should only be allowed at the start of a contained value
		if unicode.IsDigit(charRune) || unicode.IsLetter(charRune) || charRune == ':' {
			result.WriteRune(charRune)
		} else {
			return "", errors.New("unexpected character")
		}
		idx++
	}
	if lastRune != ')' {
		return "", errors.New("unexpected final value - expected bracket")
	}
	return result.String(), nil
}

func ConsumeParameter(chars []string, idx int) (*Token, int, error) {
	literal := bytes.Buffer{}
	literal.WriteString(":")
	idx += 1

	seenStartChar := false
	for {
		if idx == len(chars) {
			break
		}
		charRune := []rune(chars[idx])[0]
		if unicode.IsDigit(charRune) && !seenStartChar {
			return nil, 0, errors.New("expected a char as initial part of parameter")
		}

		if unicode.IsLetter(charRune) {
			seenStartChar = true
			literal.WriteString(chars[idx])
			idx++
		} else if unicode.IsDigit(charRune) {
			literal.WriteString(chars[idx])
			idx++
		} else {
			break
		}
	}
	return &Token{TokenType: Parameter, Literal: literal.String()}, len(literal.String()), nil
}

func ParseConditionalExpression(expression string) []*Token {
	/*var result []Token
	chars := strings.Split(expression, "")
	idx := 0

	*/
	return nil
}
