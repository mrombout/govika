package vika

import (
	"bufio"
	"bytes"
	"errors"
	"reflect"
)

type tokenStack struct {
	tokens []token
}

func (t *tokenStack) peek() *token {
	if len(t.tokens) <= 0 {
		return nil
	}
	return &t.tokens[0]
}

func (t *tokenStack) pop() *token {
	if len(t.tokens) <= 0 {
		return nil
	}

	token := t.tokens[0]
	t.tokens = t.tokens[1:]

	return &token
}

func parseMarkdownIssue(issue *Issue, content []byte) error {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	tokens, err := lex(scanner)
	if err != nil {
		return err
	}

	stack := tokenStack{
		tokens: tokens,
	}

	if err := parsePreamble(&stack, issue); err != nil {
		return err
	}

	return nil
}

func parsePreamble(stack *tokenStack, issue *Issue) error {
	if _, err := acceptToken(stack, preambleDelimiter{}); err != nil {
		return err
	}

	token, err := acceptToken(stack, unknownContent{})
	if err != nil {
		return err
	}
	content := (*token).(unknownContent).Content
	err = parseYmlIssue(issue, []byte(content))
	if err != nil {
		return err
	}

	if _, err := acceptToken(stack, preambleDelimiter{}); err != nil {
		return err
	}

	return nil
}

func parseIssueTitle(stack *tokenStack, issue *Issue) error {
	token, err := acceptToken(stack, header1Title{})
	if err != nil {
		return err
	}
	issue.Title = (*token).(header1Title).Content

	if _, err := acceptToken(stack, emptyLine{}); err != nil {
		return err
	}

	return nil
}

func parseIssueDescription(stack *tokenStack, issue *Issue) error {
	for token := stack.peek(); !isToken(stack, commentHeader{}); token = stack.peek() {
		if val, ok := (*token).(unknownContent); ok {
			issue.Description += val.Content + "\n"
		} else if _, ok := (*token).(emptyLine); ok {
			issue.Description += "\n"
		}

		stack.pop()
	}

	issue.Description = issue.Description[:len(issue.Description)-1]

	return nil
}

func parseComments(stack *tokenStack) error {
	if _, err := acceptToken(stack, commentHeader{}); err != nil {
		return err
	}
	if _, err := acceptToken(stack, emptyLine{}); err != nil {
		return err
	}
	return nil
}

func parseComment(stack *tokenStack) error {
	comment := Comment{}

	for token := stack.peek(); !isToken(stack, commentAuthor{}); token = stack.peek() {
		if val, ok := (*token).(commentContent); ok {
			comment.Message += val.Content + "\n"
		} else if _, ok := (*token).(emptyLine); ok {

		}

		stack.pop()
	}

	token, err := acceptToken(stack, emptyLine{})
	if err != nil {
		return err
	}
	comment.Author = (*token).(commentAuthor).Author

	return nil
}

func acceptToken(stack *tokenStack, tokenType interface{}) (*token, error) {
	if !isToken(stack, tokenType) {
		return nil, errors.New("unexpected token")
	}

	return stack.pop(), nil
}

func isToken(stack *tokenStack, tokenType interface{}) bool {
	if len(stack.tokens) <= 0 {
		return false
	}
	token := stack.peek()

	typeA := reflect.TypeOf(*token)
	typeB := reflect.TypeOf(tokenType)

	return typeA == typeB
}
