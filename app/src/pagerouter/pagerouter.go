package pagerouter

import (
	"com/mapify/trie"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type PageRouter struct {
	root     string
	pageTrie *trie.TemplateTrieNode
}

var contentTypeMap = map[string]string{
	"css": "text/css",
	"js":  "text/javascript",
}

func NewPageRouter(root string) (*PageRouter, error) {
	trie := trie.NewTemplateTrieNode()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		log.Printf("Found template: %s", path)
		return trie.AddTemplateTrieNode(path, "/page.html")
	})
	if err != nil {
		return nil, err
	}
	return &PageRouter{root: root, pageTrie: trie}, nil
}

func (router *PageRouter) Run(port int16) error {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		log.Printf("Request for path: %s", path)
		tmpl, vals := router.GetTemplate(path)

		if tmpl == nil {
			log.Printf("Error: Could not find template for path %s", path)
			http.Error(w, "Not Found Error", http.StatusNotFound)
			return
		}

		_, fileType, found := strings.Cut(path, ".")
		if found {
			if contentType, ok := contentTypeMap[fileType]; ok {
				w.Header().Add("Content-Type", contentType)
			}
			log.Println(w.Header()["Content-Type"])
		}

		err := tmpl.Execute(w, vals)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	log.Printf("Server started on :%d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

}

func (router *PageRouter) GetTemplate(path string) (*template.Template, trie.DynamicValuesMap) {
	dynamicValues := make(trie.DynamicValuesMap)
	path, query, found := strings.Cut(path, "?")
	if found {
		for _, param := range strings.Split(query, "&") {
			k, v, f := strings.Cut(param, "=")
			if f {
				dynamicValues[k] = v
			}
		}
	}
	templatePath := router.root + path
	return router.pageTrie.GetTemplate(templatePath, dynamicValues)
}

func (router *PageRouter) Print() {
	router.pageTrie.Print("")
}
