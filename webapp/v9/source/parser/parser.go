package parser

import (
	"../database"
	"bytes"
	"encoding/base64"
	"gopkg.in/mgo.v2/bson"
	"image/jpeg"
	"os"
	"strings"
)

func Download(file string) error {
	title := tagR(file, "book-title")
	if database.ExistDb(title) || title == "" {
		return nil
	}
	author := database.Person{Firstname: tagR(tagR(file, "author"), "first-name"), Lastname: tagR(tagR(file, "author"), "last-name")}
	annotation := tagR(file, "annotation")
	body := formatBodyAbzac(tagR(file, "body"), 1200)

	image := tagR(file, "binary")
	if image != "" {
		bs64ToImg(title, image)
	}
	code := encode(file)
	book := &database.BookDB{bson.NewObjectId(), title, author.Firstname, author.Lastname, annotation, body, code}
	return database.Unload(book)

}

func formatBodyAbzac(body string, abzac int) string {
	body1 := strings.Split(body, "")
	b := 0
	vis := false
	p := false
	a := false
	for i, v := range body1 {

		if v == "<" {

			if (body1[i+1] == "p") && (body1[i+2] == ">") {
				p = true
			}
			if (body1[i+1] == "/") && (body1[i+2] == "p") {
				p = false
			}

			if (body1[i+1] == "a") && ((body1[i+2] == " ") || (body1[i+2] == ">")) {
				a = true
			}
			if (body1[i+1] == "/") && (body1[i+2] == "a") {
				a = false
			}

			vis = false
		}

		if vis {

			if p && !a && (b >= abzac) {
				b = 0
				body1[i] += "</p><antonKovalev><p>"
			} else {
				b++
			}

		}

		if v == ">" {

			vis = true
		}

	}
	return strings.Join(body1, "")
}

func tagR(file string, tag string) string {
	ind1 := strings.Index(file, "<"+tag)
	ind2 := strings.Index(file, "</"+tag+">")
	if ind1 == -1 || ind2 == -1 {
		return ""
	}
	ind1 += strings.Index(file[ind1:], ">") + 1
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

func bs64ToImg(title string, bs string) {

	unbased, err := base64.StdEncoding.DecodeString(bs)
	if err != nil {
		panic("Cannot decode b64")
	}

	r := bytes.NewReader(unbased)
	im, err := jpeg.Decode(r)
	if err != nil {
		panic("Bad jpeg")
	}

	f, err := os.OpenFile("public/static/images/books/"+title+".jpg", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic("Cannot open file")
	}

	jpeg.Encode(f, im, nil)
}
