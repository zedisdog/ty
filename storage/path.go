package storage

import (
	"fmt"
	"path/filepath"
	"strings"
)

// NewPath 根据给定根路径创建Path实例
func NewPath(root string) *Path {
	abs, err := filepath.Abs(strings.TrimRight(root, "/"))
	if err != nil {
		panic(err)
	}
	return &Path{
		root: abs,
	}
}

type Path struct {
	root string
}

// Concat concatenate root path and given path
func (p Path) Concat(path string) string {
	if path == "" {
		return p.root
	}
	return fmt.Sprintf("%s/%s", p.root, strings.TrimLeft(path, "/"))
}

// Dir return dir name of path
func (p Path) Dir(path string) string {
	return filepath.Dir(p.Concat(path))
}
