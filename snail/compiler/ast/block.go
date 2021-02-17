// Copyright 2021 Yongqi Mu. All rights reserved.
// Use of this source code is governed by an apache
// license that can be found in the LICENSE file.

// Package ast - Abstract syntax tree definition.
package ast

/*
Block - code block

Chunk ::= Block
// type Chunk *Block
Block ::= Stat* Retstat?
Retstat ::= return Exps?
Exps ::= Exp (‘,’ Exp)*
*/
type Block struct {
	LastLine int
	Stats    []Stat
	RetExps  []Exp
}
