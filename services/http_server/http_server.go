package http_server

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/dimfu/mbinder/models"
)

const (
	PORT     = ":8090"
	TEMP_DIR = "images"
)

var (
	//go:embed templates *
	templateDir embed.FS
	templates   map[string]*template.Template
)

func init() {
	err := getAllTemplate()
	if err != nil {
		panic(err)
	}
}

func Run(items []*models.Item) {
	imgDir, err := temp(items)
	if err != nil {
		panic(err)
	}
	fs := http.FileServer(http.Dir(imgDir))
	http.Handle("/images/", http.StripPrefix("/images/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		err := renderTemplate(w, "index.tmpl", map[string]interface{}{
			"data": map[string]interface{}{
				"title": "Home",
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	server := &http.Server{Addr: PORT}
	go func() {
		fmt.Println("Starting server on http://localhost:8090")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("\nShutting down...")

	server.Shutdown(context.Background())

	os.RemoveAll(imgDir)
	fmt.Println("Server successfully shut down")
}

func temp(items []*models.Item) (string, error) {
	p, err := os.MkdirTemp("", TEMP_DIR)
	if err != nil {
		return "", err
	}
	for _, item := range items {
		f, err := os.Open(item.Path)
		if err != nil {
			return "", err
		}
		defer f.Close()
		name := path.Base(item.Path)
		dstPath := path.Join(p, name)
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return "", err
		}
		_, err = io.Copy(dstFile, f)
		if err != nil {
			return "", err
		}
		item.Path = name
	}
	return p, nil
}

func readFolder(dir string) ([]string, error) {
	var files []string
	entries, err := templateDir.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return files, nil
		}
		return files, err
	}

	for _, entry := range entries {
		fp := path.Join(dir, entry.Name())
		if entry.IsDir() {
			_files, err := readFolder(fp)
			if err != nil {
				return files, err
			}
			files = append(files, _files...)
		} else {
			files = append(files, fp)
		}

	}

	return files, nil
}

func getAllTemplate() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layouts, err := readFolder("templates/layouts")
	if err != nil {
		return err
	}

	includes, err := readFolder("templates/includes")
	if err != nil {
		return err
	}

	basePath := "templates/base.tmpl"

	for _, layout := range layouts {
		if path.Ext(layout) != ".tmpl" {
			continue
		}
		files := append(includes, layout, basePath)
		name := path.Base(layout)
		tmpl, err := template.New(name).ParseFS(templateDir, files...)
		if err != nil {
			return err
		}
		templates[name] = tmpl
	}

	return nil
}

func renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("template %s does not exists\n", name)
	}
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "base.tmpl", data)
	if err != nil {
		return err
	}

	// fmt.Println("Rendered HTML:", buf.String())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = buf.WriteTo(w)
	return err
}
