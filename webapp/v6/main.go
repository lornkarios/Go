package main

import (
	//"fmt"
	"github.com/bmizerany/pat"
	"github.com/lornkarios/Go/webapp/v6/source/daemon"
	"github.com/lornkarios/Go/webapp/v6/source/parser"

	"html/template"
	//"io/ioutil"
	"log"
	"net/http"
	//"os"
	"path"
	//"strings"
)

var (
	// компилируем шаблоны, если не удалось, то выходим
	first_template = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "main.html")))
	post_template  = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "book.html")))
	read_template  = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "reading.html")))
	error_template = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "error.html")))
)

func main() {
	// для отдачи сервером статичных файлов из папки public/static
	fs := http.FileServer(http.Dir("./public/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	mux := pat.New()
	mux.Get("/books/:page", http.HandlerFunc(postHandler))
	mux.Get("/books/:page/", http.HandlerFunc(postHandler))
	mux.Get("/reading/:page", http.HandlerFunc(readHandler))
	mux.Get("/reading/:page/", http.HandlerFunc(readHandler))
	mux.Get("/", http.HandlerFunc(postHandler))
	http.Handle("/", mux)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
