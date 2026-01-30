package tree

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

// TreeItem represents a file or folder in the UI list.
type TreeItem struct {
	Path     string
	FullPath string
	IsDir    bool
	Depth    int
}

func (i TreeItem) FilterValue() string { return i.FullPath }
func (i TreeItem) Description() string { return "" }
func (i TreeItem) Title() string {
	indent := strings.Repeat("  ", i.Depth)
	icon := getIcon(i.Path, i.IsDir)
	return fmt.Sprintf("%s%s %s", indent, icon, i.Path)
}

// Build converts a list of file paths into a compacted, sorted tree list.
func Build(paths []string) []list.Item {
	// FIX 1: Initialize root as a directory so logic works,
	// but we won't compact the root itself.
	root := &node{
		children: make(map[string]*node),
		isDir:    true,
	}

	// 1. Build the raw tree structure
	for _, path := range paths {
		parts := strings.Split(path, "/")
		current := root
		for i, part := range parts {
			if _, exists := current.children[part]; !exists {
				isDir := i < len(parts)-1
				fullPath := strings.Join(parts[:i+1], "/")
				current.children[part] = &node{
					name:     part,
					fullPath: fullPath,
					children: make(map[string]*node),
					isDir:    isDir,
				}
			}
			current = current.children[part]
		}
	}

	// FIX 2: Do NOT compact the root node itself (which would hide top-level folders).
	// Instead, compact each top-level child individually.
	for _, child := range root.children {
		compact(child)
	}

	// 3. Flatten to list items
	var items []list.Item
	flatten(root, 0, &items)
	return items
}

// -- Helpers --

type node struct {
	name     string
	fullPath string
	children map[string]*node
	isDir    bool
}

// compact recursively merges directories that contain only a single directory child.
func compact(n *node) {
	if !n.isDir {
		return
	}

	// Compact children first (bottom-up traversal)
	for _, child := range n.children {
		compact(child)
	}

	// Logic: If I am a directory, and I have exactly 1 child, and that child is also a directory...
	if len(n.children) == 1 {
		var child *node
		for _, c := range n.children {
			child = c
			break
		}

		if child.isDir {
			// Merge child into parent
			// e.g. "internal" + "ui" becomes "internal/ui"
			n.name = filepath.Join(n.name, child.name)
			n.fullPath = child.fullPath
			n.children = child.children // Inherit grandchildren
		}
	}
}

func flatten(n *node, depth int, items *[]list.Item) {
	keys := make([]string, 0, len(n.children))
	for k := range n.children {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		a, b := n.children[keys[i]], n.children[keys[j]]
		// Folders first
		if a.isDir && !b.isDir {
			return true
		}
		if !a.isDir && b.isDir {
			return false
		}
		return a.name < b.name
	})

	for _, k := range keys {
		child := n.children[k]
		*items = append(*items, TreeItem{
			Path:     child.name,
			FullPath: child.fullPath,
			IsDir:    child.isDir,
			Depth:    depth,
		})

		if child.isDir {
			flatten(child, depth+1, items)
		}
	}
}

func getIcon(name string, isDir bool) string {
	if isDir {
		return " "
	}
	ext := filepath.Ext(name)
	switch strings.ToLower(ext) {
	case ".go":
		return " "
	case ".js", ".ts", ".tsx":
		return " "
	case ".svelte":
		return " "
	case ".md":
		return " "
	case ".json":
		return " "
	case ".yml", ".yaml":
		return " "
	case ".html":
		return " "
	case ".css":
		return " "
	case ".git":
		return " "
	case ".dockerfile":
		return " "
	default:
		return " "
	}
}
