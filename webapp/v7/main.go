package main

import (
	//"fmt"
	"github.com/bmizerany/pat"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io"
	//"io/ioutil"
	"log"
	"net/http"
	//"os"
	"path"
	"strconv"
	"strings"
)

type pef struct {
	URL string
}
type abos interface{}

var (
	mainBook = &Book{Title: ""}
	// компилируем шаблоны, если не удалось, то выходим
	first_template = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "main.html")))
	post_template  = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "book.html")))
	read_template  = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "reading.html")))
	error_template = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "error.html")))
	books_template = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "library.html")))
	add_template   = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates", "add.html")))
)

func main() {
	// для отдачи сервером статичных файлов из папки public/static
	fs := http.FileServer(http.Dir("./public/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	mux := pat.New()

	mux.Get("/books/:page", http.HandlerFunc(PostHandler))
	mux.Get("/books/:page/", http.HandlerFunc(PostHandler))
	mux.Get("/", http.HandlerFunc(PostHandler))

	mux.Get("/library/", http.HandlerFunc(BookHandler))
	mux.Get("/library", http.HandlerFunc(BookHandler))

	mux.Get("/reading/:page/:bPage", http.HandlerFunc(ReadHandler))
	mux.Get("/reading/:page/", http.HandlerFunc(ReadHandler))
	mux.Get("/reading/:page", http.HandlerFunc(ReadHandler))
	mux.Get("/reading/", http.HandlerFunc(ReadHandler))
	mux.Get("/reading/:page/:bPage/", http.HandlerFunc(ReadHandler))

	mux.Get("/add", http.HandlerFunc(AddHandler))
	mux.Post("/add2", http.HandlerFunc(AddHandler))

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

	if page == "" {
		// если page пусто, то выдаем главную
		if err := first_template.ExecuteTemplate(w, "layout", nil); err != nil {
			log.Println(err.Error())
			errorHandler(w, r, 500)
		}
		return
	}
	var post *Book
	var err error
	var status int
	post = mainBook
	if mainBook.Title == "" || mainBook.Title != page {
		*post, status, err = Load(page, 0)
		if err != nil {
			errorHandler(w, r, status)
			return
		}
		mainBook = post
	}

	if err := post_template.ExecuteTemplate(w, "layout", *post); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}
}

func ReadHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	//loadFromDb()
	// Извлекаем параметр
	// Например, в http://127.0.0.1:3000/p1 page = "p1"
	// в http://127.0.0.1:3000/ page = ""

	bpage := params.Get(":bPage")
	page := params.Get(":page")

	// Путь к файлу (без расширения)
	// Например, posts/p1

	var pbook int64
	if page != "" {
		// если page не пусто, то считаем, что запрашивается файл
		// получим posts/p1.md
		if strings.Index(r.URL.Path, "sl") != -1 {
			pbook, _ = strconv.ParseInt(r.URL.Path[strings.Index(r.URL.Path, page)+len(page)+1:strings.Index(r.URL.Path, "sl")-1], 10, 32)
			http.Redirect(w, r, "/reading/"+page+"/"+strconv.Itoa(int(pbook+1)), http.StatusMovedPermanently)
			return
		}
		if strings.Index(r.URL.Path, "pr") != -1 {
			pbook, _ = strconv.ParseInt(r.URL.Path[strings.Index(r.URL.Path, page)+len(page)+1:strings.Index(r.URL.Path, "pr")-1], 10, 32)
			http.Redirect(w, r, "/reading/"+page+"/"+strconv.Itoa(int(pbook-1)), http.StatusMovedPermanently)
			return
		}

		if bpage != "" {

			pbook, _ = strconv.ParseInt(bpage, 10, 32)
			if pbook < 0 {
				http.Redirect(w, r, "/reading/"+page, http.StatusMovedPermanently)
				return
			}
			if mainBook.Body1 != "" && len(mainBook.Body1)/4000 < int(pbook) {
				http.Redirect(w, r, "/reading/"+page+"/"+strconv.Itoa(len(mainBook.Body1)/4000), http.StatusMovedPermanently)
				return
			}

		} else {
			pbook = 0
		}
	} else {
		// если page пусто, то выдаем главную
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	var post *Book
	var err error
	var status int

	mainBook.Bpage = pbook

	post = mainBook

	if mainBook.Title == "" || mainBook.Title != page {
		*post, status, err = Load(page, pbook)

		if err != nil {
			errorHandler(w, r, status)
			return
		}
		mainBook = post
	}
	if int(post.Bpage) == len(mainBook.Body1)/4000 {
		post.Body = template.HTML(post.Body1[post.Bpage*4000 : int(post.Bpage)*4000+len(post.Body1)%4000])

	} else {

		post.Body = template.HTML(post.Body1[post.Bpage*4000 : (post.Bpage+1)*4000])
	}
	if err := read_template.ExecuteTemplate(w, "layout", *post); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}
}
func BookHandler(w http.ResponseWriter, r *http.Request) {
	titles := LoadFromDb()
	var p struct {
		Boddy template.HTML
		Code  string
		Title string
	}
	s := ""
	for _, v := range titles {
		s += "<a href=\"" + "/books/" + v + "\" class=\"list-group-item\">" + v + "</a>"
	}
	p.Boddy = template.HTML(s)
	p.Code = "UTF-8"
	p.Title = ""
	if err := books_template.ExecuteTemplate(w, "layout", p); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Add")

	if r.URL.Path == "/add2" {
		var file string
		var err error
		if file, err = openFB(w, r); err != nil {
			log.Println(err.Error())
			errorHandler(w, r, 500)
		}
		if err = Download(file); err != nil {
			log.Println(err.Error())
			errorHandler(w, r, 500)
		}
		http.Redirect(w, r, "/add", http.StatusMovedPermanently)
		return

	}
	//http.PostForm("/add", data)
	//resp, _ := http.PostForm("/add", url.Values{"login_login": {"Value"}, "login_password": {"123"}})
	if err := add_template.ExecuteTemplate(w, "layout", nil); err != nil {
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
	Body1      string
}

type BookDB struct {
	Id          bson.ObjectId `bson:"_id"`
	Title       string        `bson:"title"`
	AuthorName  string        `bson:"authorName"`
	AuthorLName string        `bson:"authorLName"`
	Annotation  string        `bson:"annotation"`
	Body        string        `bson:"body"`
	Image       string        `bson:"image"`
	Code        string        `bson:"code"`
}

func Load(md string, pNum int64) (Book, int, error) {
	//info, err := os.Stat(md)
	//if err != nil {
	//	if os.IsNotExist(err) {
	// файл не существует
	//		return Book{}, http.StatusNotFound, err
	//	}
	//}

	//if info.IsDir() {
	// не файл, а папка
	//	return Book{}, http.StatusNotFound, fmt.Errorf("dir")
	//}

	session, _ := mgo.Dial("mongodb://127.0.0.1")
	mainDB := session.DB("library")
	colBooks := mainDB.C("books")
	query := bson.M{
		"title": md,
	}

	library := []BookDB{}
	colBooks.Find(query).All(&library)
	author := Person{library[0].AuthorName, library[0].AuthorLName}

	book := Book{library[0].Title, author, template.HTML(library[0].Annotation), template.HTML(library[0].Body), library[0].Image, library[0].Code, pNum, library[0].Body}
	return book, 200, nil

}

func Download(file string) error {
	title := tagR(file, "book-title")
	author := Person{Firstname: tagR(tagR(file, "author"), "first-name"), Lastname: tagR(tagR(file, "author"), "last-name")}
	annotation := tagR(file, "annotation")
	body := tagR(file, "body")
	//body = body[pNum*4000 : (pNum+1)*4000]
	image := tagR(file, "binary")
	code := encode(file)
	book := &BookDB{bson.NewObjectId(), title, author.Firstname, author.Lastname, annotation, body, image, code}

	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// получаем коллекцию
	productCollection := session.DB("library").C("books")
	err = productCollection.Insert(book)
	if err != nil {
		return err
	}
	return nil

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
func openFB(w http.ResponseWriter, r *http.Request) (string, error) {

	r.ParseMultipartForm(1024 * 1024 * 5)
	file, handler, err := r.FormFile("my_file")
	if err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
		return "", err
	}
	p := make([]byte, 1024*1024*10)

	defer file.Close()
	handler.Open()
	file.Seek(0, io.SeekStart)
	//var n int
	_, err = file.Read(p)
	if err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
		return "", err
	}
	file.Close()
	return (string(p)), nil
}
func LoadFromDb() []string {
	session, _ := mgo.Dial("mongodb://127.0.0.1")
	mainDB := session.DB("library")
	colBooks := mainDB.C("books")
	query := bson.M{}
	library := []BookDB{}
	colBooks.Find(query).All(&library)
	var titles []string
	for _, p := range library {

		titles = append(titles, p.Title)
	}
	return titles
}
