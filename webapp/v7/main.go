package main

import (
	//"fmt"
	"fmt"
	"github.com/bmizerany/pat"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
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
	mux.Get("/books/:page", http.HandlerFunc(PostHandler))
	mux.Get("/books/:page/", http.HandlerFunc(PostHandler))
	mux.Get("/reading/:page/:bPage", http.HandlerFunc(ReadHandler))
	mux.Get("/reading/:page/", http.HandlerFunc(ReadHandler))
	mux.Get("/reading/:page", http.HandlerFunc(ReadHandler))
	mux.Get("/reading/", http.HandlerFunc(ReadHandler))
	mux.Get("/reading/:page/:bPage/", http.HandlerFunc(ReadHandler))
	mux.Get("/", http.HandlerFunc(PostHandler))
	http.Handle("/", mux)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
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
		return
	}
	post, status, err := Load(post_md, 0)
	if err != nil {
		errorHandler(w, r, status)
		return
	}
	if err := post_template.ExecuteTemplate(w, "layout", post); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}
}

func ReadHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	// Извлекаем параметр
	// Например, в http://127.0.0.1:3000/p1 page = "p1"
	// в http://127.0.0.1:3000/ page = ""
	page := params.Get(":page")
	bpage := params.Get(":bPage")
	// Путь к файлу (без расширения)
	// Например, posts/p1
	p := path.Join("books", page)
	var post_md string
	var pbook int64
	if page != "" {
		// если page не пусто, то считаем, что запрашивается файл
		// получим posts/p1.md
		post_md = p + ".fb2"
		if bpage != "" {

			pbook, _ = strconv.ParseInt(bpage, 10, 32)
			if pbook < 0 {
				pbook = 0
			}
		} else {
			pbook = 0
		}
	} else {
		// если page пусто, то выдаем главную
		if err := first_template.ExecuteTemplate(w, "layout", nil); err != nil {
			log.Println(err.Error())
			errorHandler(w, r, 500)
		}
		return
	}
	post, status, err := Load(post_md, pbook)
	if err != nil {
		errorHandler(w, r, status)
		return
	}
	if err := read_template.ExecuteTemplate(w, "layout", post); err != nil {
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

type Person struct {
	Firstname string
	Lastname  string
}

type Book struct {
	Title      string
	Author     Person
	Annotation template.HTML
	Body       template.HTML
	Image      string
	Code       string
	Bpage      int64
}

func Load(md string, pNum int64) (Book, int, error) {
	info, err := os.Stat(md)
	if err != nil {
		if os.IsNotExist(err) {
			// файл не существует
			return Book{}, http.StatusNotFound, err
		}
	}
	if info.IsDir() {
		// не файл, а папка
		return Book{}, http.StatusNotFound, fmt.Errorf("dir")
	}
	fileread, _ := ioutil.ReadFile(md)
	file := string(fileread)
	title := tagR(file, "book-title")
	author := Person{Firstname: tagR(tagR(file, "author"), "first-name"), Lastname: tagR(tagR(file, "author"), "last-name")}
	annotation := tagR(file, "annotation")
	body := tagR(file, "body")
	//body = strings.Join(strings.Split(body, "\n")[pNum:pNum+200], "\n")
	body = body[pNum*4000 : (pNum+1)*4000]
	//body = body[:strings.LastIndex(body, "</p>")+4]
	image := tagR(file, "binary")
	code := encode(file)
	book := Book{title, author, template.HTML(annotation), template.HTML(body), image, code, pNum}
	return book, 200, nil

}

func tagR(file string, tag string) string {
	ind1 := strings.Index(file, "<"+tag)
	ind2 := strings.Index(file, "</"+tag+">")
	if ind1 == -1 || ind2 == -1 {
		return "-1"
	}
	file1 := []byte(file[ind1:])

	for i, v := range file1 {
		if v == '>' {
			ind1 += i + 1
			break
		}
	}

	return (file[ind1:ind2])
}
func encode(file string) string {
	ind1 := strings.Index(file, "encoding=")
	ind2 := strings.Index(file, "?>")
	if ind1 == -1 || ind2 == -1 {
		return "utf-8"
	}

	return (file[ind1+10 : ind2])
}
