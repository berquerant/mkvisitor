package main

import "log"

type (
	Tree interface {
		IsTree()
	}
	Vertex struct {
		Value  string
		Leaves []Tree
	}
	Leaf struct {
		Value string
	}
)

func (*Vertex) IsTree() {}
func (*Leaf) IsTree()   {}

type TravelVisitor struct {
	values []string
}

func (s *TravelVisitor) add(v string)      { s.values = append(s.values, v) }
func (s *TravelVisitor) VisitLeaf(v *Leaf) { s.add(v.Value) }
func (s *TravelVisitor) VisitVertex(v *Vertex) {
	s.add(v.Value)
	for _, leaf := range v.Leaves {
		switch leaf := leaf.(type) {
		case *Vertex:
			s.VisitVertex(leaf)
		case *Leaf:
			s.VisitLeaf(leaf)
		default:
			log.Fatal("unknown tree")
		}
	}
}

type LeafVisitor struct {
	VisitorDefault
	value string
}

func (s *LeafVisitor) VisitLeaf(v *Leaf) { s.value = v.Value }

func main() {
	var (
		tv = &TravelVisitor{values: []string{}}
		t1 = &Vertex{
			Value: "v1",
			Leaves: []Tree{
				&Leaf{Value: "l1"},
				&Vertex{
					Value: "v2",
					Leaves: []Tree{
						&Leaf{Value: "l2"},
						&Leaf{Value: "l3"},
					},
				},
				&Leaf{Value: "l4"},
			},
		}
	)
	t1.Accept(tv)
	if len(tv.values) != 6 {
		log.Fatalf("want 6 values but got %d values", len(tv.values))
	}
	for i, w := range []string{"v1", "l1", "v2", "l2", "l3", "l4"} {
		if tv.values[i] != w {
			log.Fatalf("want %v at %d but got %v", w, i, tv.values[i])
		}
	}
	var (
		lv = &LeafVisitor{}
		t2 = &Leaf{Value: "L"}
	)
	t1.Accept(lv)
	if lv.value != "" {
		log.Fatalf("want empty but got %s", lv.value)
	}
	t2.Accept(lv)
	if lv.value != "L" {
		log.Fatalf("want L but got %s", lv.value)
	}
}
