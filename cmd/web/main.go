package main

import (
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(&neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Запуск веб-сервера на http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs *neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		_, err := nfs.fs.Open(index)
		if err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}
