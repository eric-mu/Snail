// Copyright 2021 Yongqi Mu. All rights reserved.
// Use of this source code is governed by an apache
// license that can be found in the LICENSE file.

// Package ast - Abstract syntax tree definition.
package ast

// Stat - code statement
type Stat interface{}

// EmptyStat - `;`
type EmptyStat struct{}

// ContinueStat - continue
type ContinueStat struct {
	Line int
}

// BreakStat - continue
type BreakStat struct {
	Line int
}

type FuncCallStat = FuncCallExp

type WhileStat struct {
	Exp   Exp
	Block *Block
}
