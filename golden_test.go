package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGolden(t *testing.T) {
	dir, err := os.MkdirTemp("", "mkvisitor")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	const (
		oneStructPrivate = `type first struct{}`
		oneStructPublic  = `type First struct{}`
		twoStructs       = `type First struct{}
type Second struct{}`
	)
	for _, tc := range []*struct {
		title        string
		input        string
		vName        string
		vPrefix      string
		acceptMethod string
		acceptors    []string
		output       string
	}{
		{
			title:        "two structs",
			input:        twoStructs,
			vName:        "Visitor",
			vPrefix:      "Visit",
			acceptMethod: "Accept",
			acceptors:    []string{"First", "Second"},
			output: `type Visitor interface {
  VisitFirst(*First)
  VisitSecond(*Second)
}

func (s *First) Accept(v Visitor)  { v.VisitFirst(s) }
func (s *Second) Accept(v Visitor) { v.VisitSecond(s) }

type VisitorDefault struct{}

func (s *VisitorDefault) VisitFirst(_ *First)   {}
func (s *VisitorDefault) VisitSecond(_ *Second) {}
func VisitSwitch(visitor Visitor, v interface{}) {
  switch v := v.(type) {
  case *First:
    visitor.VisitFirst(v)
  case *Second:
    visitor.VisitSecond(v)
  default:
    panic(fmt.Sprintf("VisitSwitch cannot switch %#v", v))
  }
}
`,
		},
		{
			title:        "public one from two structs",
			input:        twoStructs,
			vName:        "Visitor",
			vPrefix:      "Visit",
			acceptMethod: "Accept",
			acceptors:    []string{"First"},
			output: `type Visitor interface {
  VisitFirst(*First)
}

func (s *First) Accept(v Visitor) { v.VisitFirst(s) }

type VisitorDefault struct{}

func (s *VisitorDefault) VisitFirst(_ *First) {}
func VisitSwitch(visitor Visitor, v interface{}) {
  switch v := v.(type) {
  case *First:
    visitor.VisitFirst(v)
  default:
    panic(fmt.Sprintf("VisitSwitch cannot switch %#v", v))
  }
}
`,
		},
		{
			title:        "public one change names",
			input:        oneStructPublic,
			vName:        "OverVisitor",
			vPrefix:      "Over",
			acceptMethod: "Deny",
			acceptors:    []string{"First"},
			output: `type OverVisitor interface {
  OverFirst(*First)
}

func (s *First) Deny(v OverVisitor) { v.OverFirst(s) }

type OverVisitorDefault struct{}

func (s *OverVisitorDefault) OverFirst(_ *First) {}
func OverSwitch(visitor OverVisitor, v interface{}) {
  switch v := v.(type) {
  case *First:
    visitor.OverFirst(v)
  default:
    panic(fmt.Sprintf("OverSwitch cannot switch %#v", v))
  }
}
`,
		},
		{
			title:        "private one",
			input:        oneStructPrivate,
			vName:        "visitor",
			vPrefix:      "visit",
			acceptMethod: "accept",
			acceptors:    []string{"first"},
			output: `type visitor interface {
  visitFirst(*first)
}

func (s *first) accept(v visitor) { v.visitFirst(s) }

type visitorDefault struct{}

func (s *visitorDefault) visitFirst(_ *first) {}
func visitSwitch(visitor visitor, v interface{}) {
  switch v := v.(type) {
  case *first:
    visitor.visitFirst(v)
  default:
    panic(fmt.Sprintf("visitSwitch cannot switch %#v", v))
  }
}
`,
		},
		{
			title:        "public one",
			input:        oneStructPublic,
			vName:        "Visitor",
			vPrefix:      "Visit",
			acceptMethod: "Accept",
			acceptors:    []string{"First"},
			output: `type Visitor interface {
  VisitFirst(*First)
}

func (s *First) Accept(v Visitor) { v.VisitFirst(s) }

type VisitorDefault struct{}

func (s *VisitorDefault) VisitFirst(_ *First) {}
func VisitSwitch(visitor Visitor, v interface{}) {
  switch v := v.(type) {
  case *First:
    visitor.VisitFirst(v)
  default:
    panic(fmt.Sprintf("VisitSwitch cannot switch %#v", v))
  }
}
`,
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			input := fmt.Sprintf("package test\n%s", tc.input)
			file := fmt.Sprintf("%s.go", tc.title)
			absFile := filepath.Join(dir, file)
			if err := os.WriteFile(absFile, []byte(input), 0600); err != nil {
				t.Error(err)
			}
			g := Generator{}
			g.parsePackage([]string{absFile})
			g.generate(tc.acceptMethod, tc.vName, tc.vPrefix, tc.acceptors)
			got := string(g.format())
			assert.Equal(t, tc.output, strings.ReplaceAll(got, "\t", "  "))
		})
	}
}
