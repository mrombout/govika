package vika

import (
	"strings"
)

// Scanner provides a convenient interface for the lexer.
// It is a subset of bufio.Scanner.
type Scanner interface {
	Err() error
	Scan() bool
	Text() string
}

type token interface {
}

type preambleDelimiter struct {
}

type header1Title struct {
	Content string
}

type emptyLine struct {
}

type commentContent struct {
	Content string
}

type commentAuthor struct {
	Author string
}

type unknownContent struct {
	Content string
}

// Lex lexes an issue into logical tokens that makes parsing easier.Lex
//
// - PREAMBLE_DELIMITER
// - YAML_CONTENT
// - HEADER1_TITLE
// - EMPTY_LINE
// - DESCRIPTION_PARAGRAPH
// - COMMENT_CONTENT
// - COMMENT_AUTHOR
func lex(scanner Scanner) ([]token, error) {
	tokens := []token{}

	for scanner.Scan() {
		line := scanner.Text()

		var currentToken token
		switch {
		case isPreambleDelimiter(line):
			currentToken = lexPreambleDelimiter(line)
		case isHeader1Title(line):
			currentToken = lexHeader1Title(line)
		case isEmptyLine(line):
			currentToken = lexEmptyLine(line)
		case isCommentContent(line):
			currentToken = lexCommentContent(line)
		case isCommentAuthor(line):
			currentToken = lexCommentAuthor(line)
		default:
			currentToken = lexUnknownContent(line)
		}

		if currentToken != nil {
			tokens = append(tokens, currentToken)
		}
	}

	if err := scanner.Err(); err != nil {
		return tokens, err
	}

	return tokens, nil
}

func isPreambleDelimiter(line string) bool {
	return line == "---"
}

func lexPreambleDelimiter(line string) preambleDelimiter {
	return preambleDelimiter{}
}

func isHeader1Title(line string) bool {
	return strings.HasPrefix(line, "# ")
}

func lexHeader1Title(line string) header1Title {
	return header1Title{
		Content: line[2:],
	}
}

func isEmptyLine(line string) bool {
	return line == ""
}

func lexEmptyLine(line string) emptyLine {
	return emptyLine{}
}

func isCommentContent(line string) bool {
	return strings.HasPrefix(line, "> ")
}

func lexCommentContent(line string) commentContent {
	return commentContent{
		Content: line[2:],
	}
}

func isCommentAuthor(line string) bool {
	return strings.HasPrefix(line, "~ ")
}

func lexCommentAuthor(line string) commentAuthor {
	return commentAuthor{
		Author: line[2:],
	}
}

// lexUnknownContent handles content that is considered too cumbersome to
// process by the lexer.
//
// The issue description handled by the parser because it depends on the context
// whether it's a valid issue description or a syntax error.
//
// The yaml preamble is handled by the parser because YAML is simply everything
// between the first and second preamble delimiter.
func lexUnknownContent(line string) unknownContent {
	return unknownContent{
		Content: line,
	}
}
