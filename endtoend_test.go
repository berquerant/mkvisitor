package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEndToEnd(t *testing.T) {
	b := newMkVisitor(t)
	defer b.close()

	for _, tc := range []*struct {
		title     string
		fileName  string
		typeNames []string
	}{
		{
			title:     "tree",
			fileName:  "tree.go",
			typeNames: []string{"Vertex", "Leaf"},
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			b.compileAndRun(t, tc.fileName, tc.typeNames)
		})
	}
}

type mkvisitor struct {
	dir, mkvisitor string
}

func newMkVisitor(t *testing.T) *mkvisitor {
	t.Helper()
	s := &mkvisitor{}
	s.init(t)
	return s
}

func (s *mkvisitor) init(t *testing.T) {
	t.Helper()
	dir, err := os.MkdirTemp("", "mkvisitor")
	if err != nil {
		t.Fatal(err)
	}
	mkvisitor := filepath.Join(dir, "mkvisitor")
	// build mkvisitor
	if err := run("go", "build", "-o", mkvisitor); err != nil {
		t.Fatal(err)
	}
	s.dir = dir
	s.mkvisitor = mkvisitor
}

func (s *mkvisitor) close() {
	os.RemoveAll(s.dir)
}

func (s *mkvisitor) compileAndRun(t *testing.T, fileName string, typeNames []string) {
	t.Helper()
	src := filepath.Join(s.dir, fileName)
	if err := copyFile(src, filepath.Join("testdata", fileName)); err != nil {
		t.Fatal(err)
	}
	mkvisitorSrc := filepath.Join(s.dir, fmt.Sprintf("%s_mkvisitor.go", strings.Split(fileName, ".")[0]))
	// run mkvisitor
	if err := run(s.mkvisitor,
		"-type", strings.Join(typeNames, ","),
		"-vType", "Visitor",
		"-vPrefix", "Visit",
		"-accept", "Accept",
		"-output", mkvisitorSrc, src); err != nil {
		t.Fatal(err)
	}
	// run testfile with generated file
	if err := run("go", "run", mkvisitorSrc, src); err != nil {
		t.Fatal(err)
	}
}

func run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func copyFile(to, from string) error {
	toFile, err := os.Create(to)
	if err != nil {
		return err
	}
	defer toFile.Close()
	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()
	_, err = io.Copy(toFile, fromFile)
	return err
}
