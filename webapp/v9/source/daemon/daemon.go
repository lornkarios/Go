package daemon

import (
	"../database"
	"../parser"
	//	"fmt"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

var (
	mainBook = &database.Book{Title: ""}
	mainUser = &database.User{Login: ""}
	// компилируем шаблоны, если не удалось, то выходим
	first_template  = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates/dop", "main1.html")))
	regis_template  = template.Must(template.ParseFiles(path.Join("templates", "index.html"), path.Join("templates/dop", "main2.html")))
	post_template   = template.Must(template.ParseFiles(path.Join("templates", "index1.html"), path.Join("templates/dop", "book.html")))
	read_template   = template.Must(template.ParseFiles(path.Join("templates", "index1.html"), path.Join("templates/dop", "reading.html")))
	error_template  = template.Must(template.ParseFiles(path.Join("templates", "index1.html"), path.Join("templates/dop", "error.html")))
	books_template  = template.Must(template.ParseFiles(path.Join("templates", "index1.html"), path.Join("templates/dop", "library.html")))
	add_template    = template.Must(template.ParseFiles(path.Join("templates", "index1.html"), path.Join("templates/dop", "add.html")))
	about1_template = template.Must(template.ParseFiles(path.Join("templates", "index1.html"), path.Join("templates/about", "about1.html")))
	about2_template = template.Must(template.ParseFiles(path.Join("templates", "index1.html"), path.Join("templates/about", "about2.html")))
	about3_template = template.Must(template.ParseFiles(path.Join("templates", "index1.html"), path.Join("templates/about", "about3.html")))
	about4_template = template.Must(template.ParseFiles(path.Join("templates", "index1.html"), path.Join("templates/about", "about4.html")))
	search_template = template.Must(template.ParseFiles(path.Join("templates", "index1.html"), path.Join("templates/dop", "search.html")))
)

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

		if mainUser.Login != "" {
			var p struct {
				Login string
				Code  string
				Title string
			}
			p.Login = mainUser.Login
			p.Code = "utf-8"
			p.Title = ""
			if err := first_template.ExecuteTemplate(w, "layout", p); err != nil {
				log.Println(err.Error())
				errorHandler(w, r, 500)
			}
			return
		}
		if err := first_template.ExecuteTemplate(w, "layout", nil); err != nil {
			log.Println(err.Error())
			errorHandler(w, r, 500)
		}
		return
	}
	var post *database.Book
	var err error
	var status int
	post = mainBook
	if mainBook.Title == "" || mainBook.Title != page {
		*post, status, err = database.Load(page, 0)
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
			if mainBook.Body1 != "" && len(strings.Split(mainBook.Body1, "<antonKovalev>"))-1 < int(pbook) {
				http.Redirect(w, r, "/reading/"+page+"/"+strconv.Itoa(len(strings.Split(mainBook.Body1, "<antonKovalev>"))-1), http.StatusMovedPermanently)
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
	var post *database.Book
	var err error
	var status int

	mainBook.Bpage = pbook

	post = mainBook

	if mainBook.Title == "" || mainBook.Title != page {
		*post, status, err = database.Load(page, pbook)

		if err != nil {
			errorHandler(w, r, status)
			return
		}
		mainBook = post
	}

	post.Body = template.HTML(strings.Split(mainBook.Body1, "<antonKovalev>")[post.Bpage])

	if err := read_template.ExecuteTemplate(w, "layout", *post); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}

}
func BookHandler(w http.ResponseWriter, r *http.Request) {
	titles := database.LoadFromDb()
	var p struct {
		Boddy template.HTML
		Code  string
		Title string
	}
	s := ""
	for _, v := range titles {
		//s += "<a href=\"" + "/books/" + v + "\" class=\"list-group-item\" >" + v + "</a>"
		s += "<li class=\"list-group-item\"><a href=\"/books/" + v + "\">" + v + "</a></li>"
	}
	p.Boddy = template.HTML(s)
	p.Code = "UTF-8"
	p.Title = ""
	if err := books_template.ExecuteTemplate(w, "layout", p); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	titles := database.LoadFromDb()
	var p struct {
		Boddy   template.HTML
		Request string
		Code    string
		Title   string
	}

	p.Request = r.FormValue("search")
	titles1 := make([]string, 0)
	for _, v := range titles {
		if strings.Index(strings.ToLower(v), strings.ToLower(p.Request)) != -1 {
			titles1 = append(titles1, v)
		}
	}
	s := ""
	for _, v := range titles1 {
		//s += "<a href=\"" + "/books/" + v + "\" class=\"list-group-item\" >" + v + "</a>"
		s += "<li class=\"list-group-item\"><a href=\"/books/" + v + "\">" + v + "</a></li>"
	}
	p.Boddy = template.HTML(s)
	p.Code = "UTF-8"
	p.Title = ""
	if err := search_template.ExecuteTemplate(w, "layout", p); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Add")
	var p struct {
		Boddy template.HTML
		Code  string
		Title string
	}
	s := ""
	p.Code = "UTF-8"
	p.Title = ""
	if r.URL.Path == "/add2" {
		var file string
		var err error
		if file, err = openFB(w, r); err != nil {
			log.Println(err.Error())
			errorHandler(w, r, 500)
		}
		if file == ".fb2" {
			http.Redirect(w, r, "/adder", http.StatusMovedPermanently)
			return

		}
		if err = parser.Download(file); err != nil {
			log.Println(err.Error())
			errorHandler(w, r, 500)
			return
		}
		http.Redirect(w, r, "/add", http.StatusMovedPermanently)
		return

	}
	if r.URL.Path == "/adder" {
		s = "<li style = \"color:#7d512d; \">Выберите файл формата fb2!!!</li>"
	}
	p.Boddy = template.HTML(s)
	//http.PostForm("/add", data)
	//resp, _ := http.PostForm("/add", url.Values{"login_login": {"Value"}, "login_password": {"123"}})
	if err := add_template.ExecuteTemplate(w, "layout", p); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Add")

	if r.URL.Path == "/about1/" {
		if err := about1_template.ExecuteTemplate(w, "layout", nil); err != nil {
			log.Println(err.Error())
			errorHandler(w, r, 500)
		}
		return

	}

	if r.URL.Path == "/about2/" {
		if err := about2_template.ExecuteTemplate(w, "layout", nil); err != nil {
			log.Println(err.Error())
			errorHandler(w, r, 500)
		}
		return

	}

	if r.URL.Path == "/about3/" {
		if err := about3_template.ExecuteTemplate(w, "layout", nil); err != nil {
			log.Println(err.Error())
			errorHandler(w, r, 500)
		}
		return

	}
	//http.PostForm("/add", data)
	//resp, _ := http.PostForm("/add", url.Values{"login_login": {"Value"}, "login_password": {"123"}})
	if err := about4_template.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}
}

func RegisHandler(w http.ResponseWriter, r *http.Request) {
	var p struct {
		User    *database.User
		Message template.HTML
		Boddy   template.HTML
		Code    string
		Title   string
	}

	p.User = &database.User{}
	p.Code = "utf-8"
	p.Title = ""
	p.Boddy = "<form name =\"login\" action=\"/reg2/\" class=\"form-horizontal main_form\" method=\"get\"><br><div class=\"form-group\"><div class=\"col-sm-offset-1 col-sm-22\"><input type=\"text\" class=\"form-control\" id=\"inputEmail3\" placeholder=\"Логин\" name=\"login\"></div></div><div class=\"form-group\"><div class=\" col-sm-offset-1 col-sm-22\"><input type=\"password\" class=\"form-control\" id=\"inputPassword3\" name=\"password\" placeholder=\"Пароль\"></div></div><div class=\"form-group\"><div class=\"col-sm-offset-1 col-sm-10\"><button type=\"submit\" class=\"btn libBtn\">Войти</button></div><div class=\"row\"><div class=\"col-sm-offset-2 col-sm-10\"><br></div><div class=\"col-sm-offset-2 col-sm-10\"><a href=\"/regis3/\" class=\"text-nowrap\">Регистрация</a></div></div></div></form>"

	if (r.URL.Path == "/reg/") || (r.URL.Path == "/reg2/") {
		p.User.Login = r.FormValue("login")

		p.User.Password = r.FormValue("password")

		var mon int
		var err, err2 error
		if mon, err, err2 = database.UserExistDb(p.User.Login, p.User.Password); (err != nil) || (err2 != nil) {
			if err != nil {
				log.Println(err.Error())
			}
			if err2 != nil {
				log.Println(err2.Error())
			}
			errorHandler(w, r, 500)
			return
		}

		if mon == 1 {
			if r.URL.Path == "/reg/" {
				http.Redirect(w, r, "/regis1/", http.StatusMovedPermanently)
				return
			}
			if r.URL.Path == "/reg2/" {

				http.Redirect(w, r, "/regis1v/", http.StatusMovedPermanently)
				return
			}
		}

		if mon == 2 {

			mainUser, _, _ = database.LoadUser(p.User.Login)
			p.Message = ""
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return

		}
		if mon == 0 {
			if r.URL.Path == "/reg/" {
				p.User.Id = bson.NewObjectId()
				if err := database.UnloadUser(p.User); err != nil {
					log.Println(err.Error())
					errorHandler(w, r, 500)
					return
				}

				http.Redirect(w, r, "/regis2/", http.StatusMovedPermanently)
				return
			}
			if r.URL.Path == "/reg2/" {

				http.Redirect(w, r, "/regis3/", http.StatusMovedPermanently)
				return
			}

		}
	}
	if r.URL.Path == "/regis1/" {
		p.Message = "Такой логин уже есть!!!"
	}
	if r.URL.Path == "/regis1v/" {
		p.Message = "Не верный пароль!!!"
	}
	if r.URL.Path == "/regis2/" {
		p.Message = "Вы успешно зарегистрировались"
	}
	if r.URL.Path == "/regis3/" {
		p.Message = "Регистрация"
		p.Boddy = "<form name =\"login\" action=\"/reg/\" class=\"form-horizontal main_form\" method=\"get\"><br><div class=\"form-group\"><div class=\"col-sm-offset-1 col-sm-22\"><input type=\"text\" class=\"form-control\" id=\"inputEmail3\" placeholder=\"Логин\" name=\"login\"></div></div><div class=\"form-group\"><div class=\" col-sm-offset-1 col-sm-22\"><input type=\"password\" class=\"form-control\" id=\"inputPassword3\" name=\"password\" placeholder=\"Пароль\"></div></div><div class=\"form-group\"><div class=\"col-sm-offset-1 col-sm-10\"><button type=\"submit\" class=\"btn libBtn\">Регистрация</button></div></div></form>"
	}
	if err := regis_template.ExecuteTemplate(w, "layout", p); err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
	}
	return

}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if err := error_template.ExecuteTemplate(w, "layout", map[string]interface{}{"Error": http.StatusText(status), "Status": status}); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
func openFB(w http.ResponseWriter, r *http.Request) (string, error) {

	r.ParseMultipartForm(1024 * 1024 * 5)
	file, handler, err := r.FormFile("my_file")
	if err != nil {
		log.Println(err.Error())
		errorHandler(w, r, 500)
		return "", err
	}
	if !(handler.Header["Content-Type"][0] == "application/x-fictionbook+xml") {
		return ".fb2", nil
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
