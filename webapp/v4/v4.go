package main

import (
	//"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	//"fmt"
)

var templates = template.Must(template.ParseFiles("v4hello.html"))
var validPath = regexp.MustCompile("^/(v4hello)/([a-zA-Z0-9]+)$")

type Person struct {
	Login    string
	Password string
}
type PersonFile []Person

func (p *Person) save() error {
	filename := "Persons.txt"

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	s := p.Login + " " + p.Password + "|"
	s1 := string(body) + s

	return ioutil.WriteFile(filename, []byte(s1), 0600)
}

func load(login string) (*Person, error) {
	filename := "Persons.txt"

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	s1 := ""
	s2 := ""
	for _, v := range body {
		if v == ' ' {
			s2 = s1
			s1 = ""
		} else {
			if v == '|' {
				if s2 == login {

					return &Person{Login: login, Password: s1}, nil
				}

				s1 = ""
				s2 = ""
			} else {
				s1 += string(v)
			}
		}

	}
	return nil, nil
}
func checkPerson(login string) (bool, error) {
	filename := "Persons.txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, err
	}

	s1 := ""
	s2 := ""
	for _, v := range body {
		if v == ' ' {
			s2 = s1
			s1 = ""
		} else {
			if v == '|' {
				if s2 == login {
					return true, nil
				}

				s1 = ""
				s2 = ""
			} else {
				s1 += string(v)
			}
		}

	}
	return false, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Person) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func viewHandler(w http.ResponseWriter, r *http.Request, login string) {
	b, _ := checkPerson(login)
	if login != "" && b {
		p, _ := load(login)
		renderTemplate(w, "v4hello", p)
		return
	}
	if login != "" || login == "|" {
		http.Redirect(w, r, "/v4hello/", http.StatusFound)
		return
	}
	p := Person{Login: "", Password: ""}
	renderTemplate(w, "v4hello", &p)

}

/*
func editHandler(w http.ResponseWriter, r *http.Request, title string) {

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {

	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
*/
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			if (r.URL.Path)[len("/v4hello/"):] != "" {
				fn(w, r, "|")
				return
			}
			fn(w, r, "")
			return
		}
		fn(w, r, m[2])
	}
}

func main() {

	http.HandleFunc("/v4hello/", makeHandler(viewHandler))
	//http.HandleFunc("/edit/", makeHandler(editHandler))
	//http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))

}
