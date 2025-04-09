package http_server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/dimfu/mbinder/models"
)

const (
	PORT     = ":8090"
	TEMP_DIR = "images"
)

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
		name := filepath.Base(item.Path)
		dstPath := filepath.Join(p, name)
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

func Run(items []*models.Item) {
	imgDir, err := temp(items)
	if err != nil {
		panic(err)
	}
	fs := http.FileServer(http.Dir(imgDir))
	http.Handle("/images/", http.StripPrefix("/images/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		for _, item := range items {
			var tags []string
			for _, tag := range item.Tags {
				tags = append(tags, tag.Name)
			}
			tagsToStr := strings.Join(tags, ", ")
			fmt.Fprintf(w, "<div><h2>%s</h2><img src='/images/%s' width='300'/><span>%s</span></div>", item.Path, item.Path, tagsToStr)
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
