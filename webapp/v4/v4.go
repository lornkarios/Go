package main

import (
	//"errors"
	//"html/template"
	"io/ioutil"
	//"log"
	//"net/http"
	//"regexp"
)

//var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
//var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Person struct {
	Login    string
	Password string
}
type PersonFile []Person

func (p *PersonFile) save() error {
	filename := "Persons.txt"
	s := ""
	for i := range *p {
		s += (*p)[i].Login + " " + (*p)[i].Password + "|"
	}

	return ioutil.WriteFile(filename, []byte(s), 0600)
}
func loadPerson() (*PersonFile, error) {
	filename := "Persons.txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var p PersonFile
	s1 := ""
	s2 := ""
	for _, v := range body {
		if v == ' ' {
			s2 = s1
			s1 = ""
		} else {
			if v == '|' {
				var p1 Person = Person{Login: s2, Password: s1}
				p = append(p, p1)
				s1 = ""
				s2 = ""
			} else {
				s1 += string(v)
			}
		}

	}
	return &p, nil
}

/*
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {

	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)

}

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

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
*/
func main() {

	persons, _ := loadPerson()
	(*persons)[0].Login = "Anton"
	(*persons)[0].Password = "1"
	(*persons)[2].Login = "Roma"
	(*persons)[2].Password = "2"
	persons.save()
	//http.HandleFunc("/view/", makeHandler(viewHandler))
	//http.HandleFunc("/edit/", makeHandler(editHandler))
	//http.HandleFunc("/save/", makeHandler(saveHandler))

	//log.Fatal(http.ListenAndServe(":8080", nil))

}
