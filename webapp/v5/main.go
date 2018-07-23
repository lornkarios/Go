package main

import (
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type Post struct {
	Title string
	Body  template.HTML
}

var (
	// компилируем шаблоны, если не удалось, то выходим
	post_template = template.Must(template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "post.html")))

	error_template = template.Must(template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "error.html")))
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
	p := path.Join("posts", page)
	var post_md string
	if page != "" {
		// если page не пусто, то считаем, что запрашивается файл
		// получим posts/p1.md
		post_md = p + ".md"
	} else {
		// если page пусто, то выдаем главную
		post_md = p + "/index.md"
	}
	post, status, err := load_post(post_md)
	if err != nil {
		errorHandler(w, r, status)
		return
	}
	if err := post_template.ExecuteTemplate(w, "layout", post); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if err := error_template.ExecuteTemplate(w, "layout", map[string]interface{}{"Error": http.StatusText(status), "Status": status}); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
		return
	}
}

// Загружает markdown-файл и конвертирует его в HTML
// Возвращает объект типа Post
// Если путь не существует или является каталогом, то возвращаем ошибку
func load_post(md string) (Post, int, error) {
	info, err := os.Stat(md)
	if err != nil {
		if os.IsNotExist(err) {
			// файл не существует
			return Post{}, http.StatusNotFound, err
		}
	}
	if info.IsDir() {
		// не файл, а папка
		return Post{}, http.StatusNotFound, fmt.Errorf("dir")
	}
	fileread, _ := ioutil.ReadFile(md)
	lines := strings.Split(string(fileread), "\n")
	title := string(lines[0])
	body := strings.Join(lines[1:len(lines)], "\n")
	body = string(blackfriday.MarkdownCommon([]byte(body)))
	post := Post{title, template.HTML(body)}
	return post, 200, nil
}
