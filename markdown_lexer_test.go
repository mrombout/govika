package vika

import (
	"bufio"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexTokensIndividually(t *testing.T) {
	testCases := []struct {
		line           string
		expectedTokens []token
	}{
		{"---", []token{preambleDelimiter{}}},
		{"# Issue title", []token{header1Title{Content: "Issue title"}}},
		{"\r\n", []token{emptyLine{}}},
		{"> Lorum ipsum", []token{commentContent{Content: "Lorum ipsum"}}},
		{"~ Jane Doe", []token{commentAuthor{Author: "Jane Doe"}}},
		{"key: value", []token{unknownContent{Content: "key: value"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.line, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(testCase.line))
			tokens, err := lex(scanner)

			if err != nil {
				t.Fatalf("expected result to be nil, but got %t", err)
			}
			if len(tokens) != len(testCase.expectedTokens) {
				t.Fatalf("expected to have lexed exactly %d tokens, but was %d (%t)", len(testCase.expectedTokens), len(tokens), tokens)
			}
			for i := 0; i < len(testCase.expectedTokens); i++ {
				actualToken := tokens[i]
				expectedToken := testCase.expectedTokens[i]

				if !reflect.DeepEqual(expectedToken, actualToken) {
					t.Errorf("expected token %d to equal %t, but was %t", i, expectedToken, actualToken)
				}
			}
		})
	}
}

type ErroringScanner struct {
}

func (ErroringScanner) Err() error {
	return errors.New("this error is expected")
}

func (ErroringScanner) Scan() bool { return false }

func (ErroringScanner) Text() string { return "" }

func TestLexErrorDuringLexing(t *testing.T) {
	scanner := ErroringScanner{}

	_, err := lex(scanner)

	if scanner.Err().Error() != err.Error() {
		t.Errorf("expected Lex to return error %t , but was %t", scanner.Err(), err)
	}
}

func TestIsPreambleDelimiterValid(t *testing.T) {
	testCases := []struct {
		name       string
		validInput string
	}{{"plain", "---"}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			result := isPreambleDelimiter(tc.validInput)

			// assert
			assert.True(t, result)
		})
	}
}

func TestLexPreambleDelimiterValid(t *testing.T) {
	testCases := []struct {
		name          string
		validInput    string
		expectedToken preambleDelimiter
	}{{"plain", `---
`, preambleDelimiter{}}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			result := lexPreambleDelimiter(tc.validInput)

			// assert
			assert.Equal(t, tc.expectedToken, result)
		})
	}
}

func TestIsHeader1TitleValid(t *testing.T) {
	testCases := []struct {
		name       string
		validInput string
	}{{"plain", `# Some header`}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			result := isHeader1Title(tc.validInput)

			// assert
			assert.True(t, result)
		})
	}
}

func TestLexHeader1TitleValid(t *testing.T) {
	testCases := []struct {
		name          string
		validInput    string
		expectedToken header1Title
	}{{"plain", `# Some header`, header1Title{Content: "Some header"}}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			result := lexHeader1Title(tc.validInput)

			// assert
			assert.Equal(t, tc.expectedToken, result)
		})
	}
}

func TestIsEmptyLineValid(t *testing.T) {
	testCases := []struct {
		name       string
		validInput string
	}{{"plain", ``}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			result := isEmptyLine(tc.validInput)

			// assert
			assert.True(t, result)
		})
	}
}

func TestLexEmptyLineValid(t *testing.T) {
	testCases := []struct {
		name          string
		validInput    string
		expectedToken emptyLine
	}{{"plain", ``, emptyLine{}}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			result := lexEmptyLine(tc.validInput)

			// assert
			assert.Equal(t, tc.expectedToken, result)
		})
	}
}

func TestIsCommentContentValid(t *testing.T) {
	testCases := []struct {
		name       string
		validInput string
	}{{"plain", `> This is considered comment content`}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			result := isCommentContent(tc.validInput)

			// assert
			assert.True(t, result)
		})
	}
}

func TestLexCommentContentValid(t *testing.T) {
	testCases := []struct {
		name          string
		validInput    string
		expectedToken commentContent
	}{{"plain", `> This is considered comment content`, commentContent{Content: "This is considered comment content"}}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			result := lexCommentContent(tc.validInput)

			// assert
			assert.Equal(t, tc.expectedToken, result)
		})
	}
}

func TestIsCommentAuthorValid(t *testing.T) {
	testCases := []struct {
		name       string
		validInput string
	}{{"plain", `~ John Doe`}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			result := isCommentAuthor(tc.validInput)

			// assert
			assert.True(t, result)
		})
	}
}

func TestLexCommentAuthorValid(t *testing.T) {
	testCases := []struct {
		name          string
		validInput    string
		expectedToken commentAuthor
	}{{"plain", `~ John Doe`, commentAuthor{Author: "John Doe"}}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			result := lexCommentAuthor(tc.validInput)

			// assert
			assert.Equal(t, tc.expectedToken, result)
		})
	}
}

func TestLexUnknownContentValid(t *testing.T) {
	testCases := []struct {
		name          string
		validInput    string
		expectedToken unknownContent
	}{
		{"IssueDescription", `Unknown issue description content.`, unknownContent{Content: "Unknown issue description content."}},
		{"YamlPreamble", `key: value`, unknownContent{Content: "key: value"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			result := lexUnknownContent(tc.validInput)

			// assert
			assert.Equal(t, tc.expectedToken, result)
		})
	}
}
