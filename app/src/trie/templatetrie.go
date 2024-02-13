package trie

import (
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

type TemplateTrieNode struct {
	children   map[string]*TemplateTrieNode
	template   *template.Template
	dynamicKey string
}

const dynamicSegment = "*"

type DynamicValuesMap = map[string]string

func NewTemplateTrieNode() *TemplateTrieNode {
	return &TemplateTrieNode{
		children: make(map[string]*TemplateTrieNode),
	}
}

// Creates TemplateTrie from a path.
// Any segment that has square brackets will be parsed as a dynamic segment.
// No 2 files can be added for one directory path.
// No 2 dynamic segments can be added in the same directory.
// No 2 dynamic segments can be added with the same dynamic key in the same path.
// I.e. configs/[id]/config/page.html will be parsed for current node to
//
//	node *TemplateTrieNode{
//	  children: map[string]*TemplateTrieNode{
//	    "configs": *TemplateTrieNode{
//	      children: map[string]*TemplateTrieNode{
//	        "*": *TemplateTrieNode{
//	          children: map[string]*TemplateTrieNode{
//	            "config": *TemplateTrieNode{
//	               children: map[string]*TemplateTrieNode{}
//	               template: template.ParseFiles("configs/[id]/config/page.html")
//	               dynamicKey: ""}},
//	          template: nil,
//	          dynamicKey: "id"}},
//	      template: nil,
//	      dynamicKey: ""}},
//	  template: nil,
//	  dynamicKey: ""}
func (node *TemplateTrieNode) AddTemplateTrieNode(path string, trimSuffix string) error {
	templ, err := template.ParseFiles(path)
	if err != nil {
		return err
	}
	path = strings.TrimSuffix(path, trimSuffix)
	segments := strings.Split(path, string(filepath.Separator))
	return node.addTemplateNode(segments, templ)
}

// Retrieves template for path.
// Returns nil if no such template exists.
// Skips over empty path segments, i.e. `/////config` == `/config`
// Always expects `/` as a seperator
func (node *TemplateTrieNode) GetTemplate(path string, dynamicValues DynamicValuesMap) (*template.Template, DynamicValuesMap) {
	segments := strings.Split(path, "/")
	return node.getTemplate(segments, dynamicValues)
}

func (node *TemplateTrieNode) addTemplateNode(segments []string, templ *template.Template) error {
	getNextNode := func(segment string, dynamicKey string) *TemplateTrieNode {
		_, exists := node.children[segment]
		if !exists {
			node.children[segment] = &TemplateTrieNode{
				children:   make(map[string]*TemplateTrieNode),
				dynamicKey: dynamicKey,
			}
		}
		return node.children[segment]
	}

	if len(segments) == 1 {
		getNextNode(segments[0], "").template = templ
		return nil
	}
	dir := segments[0]
	isDynamic := strings.HasPrefix(dir, "[") && strings.HasSuffix(dir, "]")

	if isDynamic && len(dir) < 3 || len(dir) == 0 {
		return errors.New("TemplateTrieNode: directory name length too short")
	}

	if isDynamic {
		dynamicKey := dir[1 : len(dir)-1]
		getNextNode(dynamicSegment, dynamicKey).addTemplateNode(segments[1:], templ)
	} else {
		getNextNode(dir, "").addTemplateNode(segments[1:], templ)
	}
	return nil
}

func (node *TemplateTrieNode) getTemplate(segments []string, dynamicValues DynamicValuesMap) (*template.Template, DynamicValuesMap) {
	if len(segments) == 0 {
		return node.template, dynamicValues
	}
	if len(segments[0]) == 0 {
		return node.getTemplate(segments[1:], dynamicValues)
	}
	if next, exists := node.children[segments[0]]; exists {
		return next.getTemplate(segments[1:], dynamicValues)
	}
	if next, exists := node.children[dynamicSegment]; exists && len(next.dynamicKey) > 0 {
		dynamicValues[next.dynamicKey] = segments[0]
		return next.getTemplate(segments[1:], dynamicValues)
	}
	return nil, dynamicValues
}

func (node *TemplateTrieNode) Print(prefix string) {
	templateIndicator := ""
	if node.template != nil {
		templateIndicator = "Template"
	}

	if prefix != "" {
		fmt.Printf("%s %s\n", prefix, templateIndicator)
	}

	for key, child := range node.children {
		newPrefix := strings.Split(prefix, "-")[0] + "  -"
		displayKey := newPrefix + " " + key + ":"
		if child.dynamicKey != "" {
			displayKey = newPrefix + " " + fmt.Sprintf("[%s]:", child.dynamicKey)
		}

		child.Print(displayKey)
	}
}
