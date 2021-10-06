# mkvisitor

Given

```go
package example

type (
	Node struct{}
	Leaf struct{}
)
```

run `mkvisitor -type "Node,Leaf"` then generate

```go
package example

import "fmt"

type Visitor interface {
	VisitNode(*Node)
	VisitLeaf(*Leaf)
}

func (s *Node) Accept(v Visitor) { v.VisitNode(s) }
func (s *Leaf) Accept(v Visitor) { v.VisitLeaf(s) }

type VisitorDefault struct{}

func (s *VisitorDefault) VisitNode(_ *Node) {}
func (s *VisitorDefault) VisitLeaf(_ *Leaf) {}
func VisitSwitch(visitor Visitor, v interface{}) {
	switch v := v.(type) {
	case *Node:
		visitor.VisitNode(v)
	case *Leaf:
		visitor.VisitLeaf(v)
	default:
		panic(fmt.Sprintf("VisitSwitch cannot switch %#v", v))
	}
}
```

in visitor_mkvisitor.go in the same directory.
