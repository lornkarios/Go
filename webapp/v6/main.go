package main

import (
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/lornkarios/webapp/v6/parser"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	// компилируем шаблоны, если не удалось, то выходим
	first_template = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "main.html")))
	post_template  = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "book.html")))

	error_template = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "error.html")))
)

func main() {
	// для отдачи сервером статичных файлов из папки public/static
	fs := http.FileServer(http.Dir("./public/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	mux := pat.New()
	mux.Get("/:page", http.HandlerFunc(postHandler))
	mux.Get("/:page/", http.HandlerFunc(postHandler))
	mux.Get("/", http.HandlerFunc(postHandler))
	http.Handle("/", mux)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	// Извлекаем параметр
	// Например, в http://127.0.0.1:3000/p1 page = "p1"
	// в http://127.0.0.1:3000/ page = ""
	page := params.Get(":page")
	// Путь к файлу (без расширения)
	// Например, posts/p1
	p := path.Join("books", page)
	var post_md string
	if page != "" {
		// если page не пусто, то считаем, что запрашивается файл
		// получим posts/p1.md
		post_md = p + ".fb2"
	} else {
		// если page пусто, то выдаем главную
		if err := first_template.ExecuteTemplate(w, "layout", nil); err != nil {
			log.Println(err.Error())
			errorHandler(w, r, 500)
		}
	}
	post, status, err := parser.Load(post_md)
	if err != nil {
		errorHandler(w, r, status)
		return
	}
	if err := post_template.ExecuteTemplate(w, "layout", post); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if err := error_template.ExecuteTemplate(w, "layout", map[string]interface{}{"Error": http.StatusText(status), "Status": status}); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
