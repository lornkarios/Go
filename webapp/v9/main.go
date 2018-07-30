package main

import (
	"./source/daemon"
	"github.com/bmizerany/pat"
	"log"
	"net/http"
)

func main() {
	// для отдачи сервером статичных файлов из папки public/static
	fs := http.FileServer(http.Dir("./public/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	mux := pat.New()

	mux.Get("/books/:page", http.HandlerFunc(daemon.PostHandler))
	mux.Get("/books/:page/", http.HandlerFunc(daemon.PostHandler))
	mux.Get("/", http.HandlerFunc(daemon.PostHandler))

	mux.Get("/library/", http.HandlerFunc(daemon.BookHandler))
	mux.Get("/library", http.HandlerFunc(daemon.BookHandler))

	mux.Get("/reading/:page/:bPage", http.HandlerFunc(daemon.ReadHandler))
	mux.Get("/reading/:page/", http.HandlerFunc(daemon.ReadHandler))
	mux.Get("/reading/:page", http.HandlerFunc(daemon.ReadHandler))
	mux.Get("/reading/", http.HandlerFunc(daemon.ReadHandler))
	mux.Get("/reading/:page/:bPage/", http.HandlerFunc(daemon.ReadHandler))

	mux.Get("/add", http.HandlerFunc(daemon.AddHandler))
	mux.Get("/adder", http.HandlerFunc(daemon.AddHandler))
	mux.Post("/add2", http.HandlerFunc(daemon.AddHandler))

	mux.Get("/about1/", http.HandlerFunc(daemon.AboutHandler))
	mux.Get("/about2/", http.HandlerFunc(daemon.AboutHandler))
	mux.Get("/about3/", http.HandlerFunc(daemon.AboutHandler))
	mux.Get("/about4/", http.HandlerFunc(daemon.AboutHandler))

	mux.Get("/search/", http.HandlerFunc(daemon.SearchHandler))

	mux.Get("/reg/", http.HandlerFunc(daemon.RegisHandler))
	mux.Get("/reg2/", http.HandlerFunc(daemon.RegisHandler))
	mux.Get("/regis1/", http.HandlerFunc(daemon.RegisHandler))
	mux.Get("/regis1v/", http.HandlerFunc(daemon.RegisHandler))
	mux.Get("/regis2/", http.HandlerFunc(daemon.RegisHandler))
	mux.Get("/regis3/", http.HandlerFunc(daemon.RegisHandler))
	http.Handle("/", mux)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
