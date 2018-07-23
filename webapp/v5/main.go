package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

var (
	// компилируем шаблоны, если не удалось, то выходим
	post_template = template.Must(template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "post.html")))
)

func main() {
	// для отдачи сервером статичных файлов из папки public/static
	fs := http.FileServer(http.Dir("./public/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", postHandler)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// обработчик запросов
	if err := post_template.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
