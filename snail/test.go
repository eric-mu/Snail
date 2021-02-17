package main

import (
	"fmt"
	"io/ioutil"
	"os"
	lex "snail/compiler/lexer"
)

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		testLexer(string(data), os.Args[1])
	}
}

func testLexer(chunk, chunkName string) {
	lexer := lex.NewLexer(chunk, chunkName)
	for {
		line, kind, token := lexer.NextToken()
		fmt.Printf("line %d: [%d][%-10s] %s\n", line, kind, kindToName(kind), token)
		if kind == lex.TokenEOF {
			break
		}
	}
}

func kindToName(kind int) string {
	switch {
	case kind < lex.TokenSepComma:
		return "EOF"
	case kind <= lex.TokenSepRCurly:
		return "separator"
	case kind <= lex.TokenOpNot:
		return "operator"
	case kind <= lex.TokenKwWhile:
		return "keyword"
	case kind == lex.TokenIdentifier:
		return "identifier"
	case kind == lex.TokenLiteralNumber:
		return "number"
	case kind == lex.TokenLiteralString:
		return "string"
	case kind == lex.TokenLiteralBytes:
		return "bytes"
	default:
		return "other"
	}
}
