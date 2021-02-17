// Copyright 2021 Yongqi Mu. All rights reserved.
// Use of this source code is governed by an apache
// license that can be found in the LICENSE file.

// Package lexer - Token constants definition.
package lexer

// Token kinds
const (
	TokenEOF = iota // end-of-file

	TokenSepComma  // ,
	TokenSepDot    // .
	TokenSepDomain // ::
	TokenSepLParen // (
	TokenSepRParen // )
	TokenSepLBrack // [
	TokenSepRBrack // ]
	TokenSepLCurly // {
	TokenSepRCurly // }

	TokenOpAssign   // =
	TokenOpBNot     // ~ (or bxor)
	TokenOpAdd      // +
	TokenOpMinus    // - (or sub/unm)
	TokenOpMul      // *
	TokenOpPackDict // **
	TokenOpDiv      // /
	TokenOpIDiv     // //
	TokenOpPow      // ^
	TokenOpMod      // %
	TokenOpBand     // &
	TokenOpBor      // |
	TokenOpLT       // <
	TokenOpLE       // <=
	TokenOpGT       // >
	TokenOpGE       // >=
	TokenOpEQ       // ==
	TokenOpNE       // !=
	TokenOpAnd      // and
	TokenOpOr       // or
	TokenOpNot      // not

	TokenKwBreak    // break
	TokenKwContinue // continue
	TokenKwElse     // else
	TokenKwElseIf   // elseif
	TokenKwFalse    // false
	TokenKwFor      // for
	TokenKwFunction // func
	TokenKwClass    // class
	TokenKwIf       // if
	TokenKwIn       // in
	TokenKwNil      // nil
	TokenKwReturn   // return
	TokenKwTrue     // true
	TokenKwWhile    // while

	TokenIdentifier    // identifier
	TokenLiteralNumber // number literal
	TokenLiteralString // string literal
	TokenLiteralBytes  // bytes literal

	TokenOpUNM         = TokenOpMinus
	TokenOpSub         = TokenOpMinus
	TokenOpBXor        = TokenOpBNot
	TokenOpPackTuple   = TokenOpMul
	TokenOpUnpackTuple = TokenOpMul
	TokenOpUnpackDict  = TokenOpPackDict
)

// content keywords to token kind constants
var keywords = map[string]int{
	"and":      TokenOpAnd,
	"break":    TokenKwBreak,
	"continue": TokenKwContinue,
	"else":     TokenKwElse,
	"elseif":   TokenKwElseIf,
	"false":    TokenKwFalse,
	"for":      TokenKwFor,
	"func":     TokenKwFunction,
	"class":    TokenKwClass,
	"if":       TokenKwIf,
	"in":       TokenKwIn,
	"nil":      TokenKwNil,
	"not":      TokenOpNot,
	"or":       TokenOpOr,
	"return":   TokenKwReturn,
	"true":     TokenKwTrue,
	"while":    TokenKwWhile,
}
