package main

// Src from https://elithrar.github.io/article/approximating-html-template-inheritance/

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var templates map[string]*template.Template

// Load templates on program initialisation
func InitTemplates(templatesDir string) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	fmt.Println("searching for Templates in : ", path.Join(templatesDir+"/layouts/*.tmpl"))

	layouts, err := filepath.Glob(path.Join(templatesDir + "/layouts/*.tmpl"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("searching for Templates in : ", path.Join(templatesDir+"/includes/*.tmpl"))
	includes, err := filepath.Glob(path.Join(templatesDir + "/includes/*.tmpl"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(layouts, includes)

	// Generate our templates map from our layouts/ and includes/ directories
	for _, layout := range layouts {
		files := append(includes, layout)
		templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
	}

	fmt.Println(templates)
}

func GetTemplate(templateName string) (tmpl *template.Template) {
	// Ensure the template exists in the map.
	tmpl, ok := templates[templateName]
	if !ok {
		Check(fmt.Errorf("The template %s does not exist.", templateName))
	}
	return
}

// renderTemplate is a wrapper around template.ExecuteTemplate.
func RenderTemplate(w http.ResponseWriter, templateName string, data map[string]interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return GetTemplate(templateName).ExecuteTemplate(w, "base", data)
}

func RenderTemplateToFile(fileName, templateName string, data map[string]interface{}) error {
	fd, err := os.Create(fileName)
	Check(err)
	defer fd.Close()
	return GetTemplate(templateName).ExecuteTemplate(fd, "base", data)
}
