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

type Visitor interface {
	VisitNode(*Node)
	VisitLeaf(*Leaf)
}

func (s *Node) Accept(v Visitor) { v.VisitNode(s) }
func (s *Leaf) Accept(v Visitor) { v.VisitLeaf(s) }

type VisitorDefault struct{}

func (s *VisitorDefault) VisitNode(_ *Node) {}
func (s *VisitorDefault) VisitLeaf(_ *Leaf) {}
```

in visitor_mkvisitor.go in the same directory.
