package parser

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

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
}

func Load(md string) (Book, int, error) {
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
	image := tagR(file, "binary")
	book := Book{title, author, template.HTML(annotation), template.HTML(body), image}
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
