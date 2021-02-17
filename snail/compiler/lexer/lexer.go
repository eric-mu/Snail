// Copyright 2021 Yongqi Mu. All rights reserved.
// Use of this source code is governed by an apache
// license that can be found in the LICENSE file.

// Package lexer - Lexer definition.
package lexer

import (
	"fmt"
	"regexp"
	"strings"
)

var reNumber = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9+])?|^[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)
var reIdentifier = regexp.MustCompile(`[_a-zA-Z][_0-9a-zA-Z]*`)

// Lexer definition
type Lexer struct {
	chunk         string
	chunkName     string
	line          int
	nextToken     string
	nextTokenKind int
	nextTokenLine int
}

// NewLexer - lexer constructor
func NewLexer(chunk, chunkName string) *Lexer {
	return &Lexer{chunk, chunkName, 1, "", 0, 0}
}

// NextToken - return next token information
func (lexer *Lexer) NextToken() (line, kind int, token string) {
	if lexer.nextTokenLine > 0 { // hit cache
		line = lexer.nextTokenLine
		kind = lexer.nextTokenKind
		token = lexer.nextToken
		lexer.line = lexer.nextTokenLine
		lexer.nextTokenLine = 0
		return
	}

	lexer.skipWhiteSpaces()
	if len(lexer.chunk) == 0 {
		return lexer.line, TokenEOF, "EOF"
	}

	c := lexer.chunk[0]
	switch c {
	case ',':
		lexer.forward(1)
		return lexer.line, TokenSepComma, ","
	case '(':
		lexer.forward(1)
		return lexer.line, TokenSepLParen, "("
	case ')':
		lexer.forward(1)
		return lexer.line, TokenSepRParen, ")"
	case '[':
		lexer.forward(1)
		return lexer.line, TokenSepLBrack, "["
	case ']':
		lexer.forward(1)
		return lexer.line, TokenSepRBrack, "]"
	case '{':
		lexer.forward(1)
		return lexer.line, TokenSepLCurly, "{"
	case '}':
		lexer.forward(1)
		return lexer.line, TokenSepRCurly, "}"
	case '+':
		lexer.forward(1)
		return lexer.line, TokenOpAdd, "+"
	case '-':
		lexer.forward(1)
		return lexer.line, TokenOpMinus, "-"
	case '^':
		lexer.forward(1)
		return lexer.line, TokenOpPow, "^"
	case '%':
		lexer.forward(1)
		return lexer.line, TokenOpMod, "%"
	case '&':
		lexer.forward(1)
		return lexer.line, TokenOpBand, "&"
	case '|':
		lexer.forward(1)
		return lexer.line, TokenOpBor, "|"
	case ':':
		if !lexer.nextStringIs("::") {
			lexer.syntaxError("Invalid character `:`")
		}
		lexer.forward(2)
		return lexer.line, TokenSepDomain, "::"
	case '!':
		if !lexer.nextStringIs("!=") {
			lexer.syntaxError("Invalid character `!`")
		}
		lexer.forward(2)
		return lexer.line, TokenOpNE, "!="
	case '*':
		if lexer.nextStringIs("**") {
			lexer.forward(2)
			return lexer.line, TokenOpPackDict, "**"
		}
		lexer.forward(1)
		return lexer.line, TokenOpMul, "*"
	case '=':
		if lexer.nextStringIs("==") {
			lexer.forward(2)
			return lexer.line, TokenOpEQ, "=="
		}
		lexer.forward(1)
		return lexer.line, TokenOpAssign, "="
	case '<':
		if lexer.nextStringIs("<=") {
			lexer.forward(2)
			return lexer.line, TokenOpLE, "<="
		}
		lexer.forward(1)
		return lexer.line, TokenOpLT, "<"
	case '>':
		if lexer.nextStringIs(">=") {
			lexer.forward(2)
			return lexer.line, TokenOpGE, ">="
		}
		lexer.forward(1)
		return lexer.line, TokenOpGT, ">"
	case '.':
		if len(lexer.chunk) == 1 || !isDigit(lexer.chunk[1]) {
			lexer.forward(1)
			return lexer.line, TokenSepDot, "."
		}
	case '\'':
		if lexer.nextStringIs("'''") {
			return lexer.line, TokenLiteralString, lexer.scanLongString('\'')
		}
		return lexer.line, TokenLiteralString, lexer.scanShortString('\'')
	case '"':
		if lexer.nextStringIs("\"\"\"") {
			return lexer.line, TokenLiteralString, lexer.scanLongString('"')
		}
		return lexer.line, TokenLiteralString, lexer.scanShortString('"')
	}

	if lexer.nextStringIs("b'") || lexer.nextStringIs("b\"") {
		quote := lexer.chunk[1]
		return lexer.line, TokenLiteralBytes, string(lexer.scanBytes(quote))
	}
	if c == '.' || isDigit(c) {
		token := lexer.scanNumber()
		return lexer.line, TokenLiteralNumber, token
	}
	if c == '_' || isLetter(c) {
		token := lexer.scanIdentifier()
		if kind, ok := keywords[token]; ok {
			return lexer.line, kind, token // keyword
		}
		return lexer.line, TokenIdentifier, token // identifier
	}

	lexer.syntaxError(fmt.Sprintf("Unexpected symbol near '%q'", c))
	return
}

// NextTokenAssertKind - Get next token and assert its kind
func (lexer *Lexer) NextTokenAssertKind(kind int) (line int, token string) {
	line, _kind, token := lexer.NextToken()
	if kind != _kind {
		err := fmt.Sprintf("Syntax error near '%s'", token)
		lexer.syntaxError(err)
	}
	return
}

// NextTokenAssertIdentifier - Get next token and assert it's an identifier
func (lexer *Lexer) NextTokenAssertIdentifier() (line int, token string) {
	return lexer.NextTokenAssertKind(TokenIdentifier)
}

// LookAhead - view the next token but not forward
func (lexer *Lexer) LookAhead() int {
	if lexer.nextTokenLine > 0 {
		return lexer.nextTokenKind
	}
	currentLine := lexer.line
	line, kind, token := lexer.NextToken()
	lexer.line = currentLine
	lexer.nextTokenLine = line
	lexer.nextTokenKind = kind
	lexer.nextToken = token
	return kind
}

func (lexer *Lexer) syntaxError(msg string) {
	err := fmt.Sprintf("%s line %d: %s", lexer.chunkName, lexer.line, msg)
	panic(err)
}

// Line - Return current line number of the lexer
func (lexer *Lexer) Line() int {
	return lexer.line
}

// TODO: forward only once
func (lexer *Lexer) skipWhiteSpaces() {
	for len(lexer.chunk) > 0 {
		if lexer.nextStringIs("//") {
			lexer.skipShortComment()
		} else if lexer.nextStringIs("/*") {
			lexer.skipLongComment()
		} else if lexer.nextStringIs("\r\n") || lexer.nextStringIs("\n\r") {
			lexer.forward(2)
			lexer.line++
		} else if isNewLine(lexer.chunk[0]) {
			lexer.forward(1)
			lexer.line++
		} else if isWhiteSpace(lexer.chunk[0]) {
			lexer.forward(1)
		} else {
			break
		}
	}
}

func (lexer *Lexer) skipShortComment() {
	lexer.forward(2)
	for len(lexer.chunk) > 0 && !isNewLine(lexer.chunk[0]) {
		lexer.forward(1)
	}
}

func (lexer *Lexer) skipLongComment() {
	lexer.forward(2)
	for len(lexer.chunk) > 0 && !lexer.nextStringIs("*/") {
		if lexer.nextStringIs("\r\n") || lexer.nextStringIs("\n\r") {
			lexer.forward(2)
			lexer.line++
		} else if isNewLine(lexer.chunk[0]) {
			lexer.forward(1)
			lexer.line++
		} else {
			lexer.forward(1)
		}
	}
	lexer.forward(2)
}

func (lexer *Lexer) nextStringIs(s string) bool {
	return strings.HasPrefix(lexer.chunk, s)
}

func (lexer *Lexer) forward(n int) {
	lexer.chunk = lexer.chunk[n:]
}

// TODO: escape
func (lexer *Lexer) scanShortString(quote byte) string {
	index := 0
	for i := 1; i < len(lexer.chunk); i++ {
		c := lexer.chunk[i]
		if isNewLine(c) {
			break
		}
		if c == quote {
			index = i
			break
		}
	}
	if index == 0 {
		lexer.syntaxError("Cannot found closing quotation mark for a short string")
	}
	str := lexer.chunk[1:index]
	lexer.forward(index + 1)
	return str
}

func (lexer *Lexer) scanLongString(quote byte) string {
	closingQuotes := string([]byte{quote, quote, quote})
	index := 0
	skipNextCharacter := false
	for i := 3; i < len(lexer.chunk); i++ {
		if skipNextCharacter {
			skipNextCharacter = false
			continue
		}
		c := lexer.chunk[i]
		if lexer.nextStringIs(closingQuotes) {
			index = i
			break
		}
		if lexer.nextStringIs("\r\n") || lexer.nextStringIs("\n\r") {
			skipNextCharacter = true
			lexer.line++
			continue
		} else if isNewLine(c) {
			lexer.line++
		}
	}
	if index == 0 {
		lexer.syntaxError("Cannot found closing quotation mark for a long string")
	}
	str := lexer.chunk[3:index]
	lexer.forward(index + 3)
	return str
}

func (lexer *Lexer) scanBytes(quote byte) []byte {
	bytes := []byte{}
	index := 0
	for i := 1; i < len(lexer.chunk); i++ {
		c := lexer.chunk[i]
		if c == quote {
			index = i
			break
		}
		bytes = append(bytes, c)
	}
	if index == 0 {
		lexer.syntaxError("Cannot found closing quotation mark for a short bytes")
	}
	lexer.forward(index + 1)
	return bytes
}

func (lexer *Lexer) scanNumber() string {
	return lexer.scan(reNumber)
}

func (lexer *Lexer) scanIdentifier() string {
	return lexer.scan(reIdentifier)
}

func (lexer *Lexer) scan(re *regexp.Regexp) string {
	if token := re.FindString(lexer.chunk); token != "" {
		lexer.forward(len(token))
		return token
	}
	panic("unreachable!")
}

func isWhiteSpace(c byte) bool {
	switch c {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}
	return false
}

func isNewLine(c byte) bool {
	return c == '\r' || c == '\n'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}
