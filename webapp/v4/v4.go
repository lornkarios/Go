package main

import (
	//"errors"

	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("v4hello.html", "v4brain.html", "v4page.html"))
var validPath = regexp.MustCompile("^/(v4hello|v4brain|v4page)/([a-zA-Z0-9]+)$")

type Person struct {
	Login    string
	Password string
	Info     string
}
type PersonFile []Person

func (p *Person) save() error {
	filename := "Persons.txt"

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	s := p.Login + ":" + p.Password + "|"
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
		if v == ':' {
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

func loadInfo(login string) (string, error) {
	filename := "pinfo.txt"

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	s1 := ""
	s2 := ""
	for _, v := range body {
		if v == '|' {
			s2 = s1
			s1 = ""
		} else {
			if v == '/' {
				if s2 == login {

					return s1, nil
				}

				s1 = ""
				s2 = ""
			} else {
				s1 += string(v)
			}
		}

	}
	return "", nil
}

func checkLogin(login string) (bool, error) {
	filename := "Persons.txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, err
	}

	s1 := ""
	s2 := ""
	for _, v := range body {
		if v == ':' {
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

func checkPerson(p *Person) (bool, error) {
	filename := "Persons.txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, err
	}

	s1 := ""
	s2 := ""
	for _, v := range body {
		if v == ':' {
			s2 = s1
			s1 = ""
		} else {
			if v == '|' {
				if s1 == p.Password && s2 == p.Login {
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
	b, _ := checkLogin(login)
	//b1, _ := checkLogin((login)[1:])
	if login != "" && b {
		p, _ := load(login)
		renderTemplate(w, "v4hello", p)
		return
	}
	if login != "" || login == "|" {
		http.Redirect(w, r, "/v4hello/", http.StatusFound)
		return
	}
	//if b1 && login[1] == '|'{

	//}
	p := Person{Login: "", Password: "", Info: ""}
	renderTemplate(w, "v4hello", &p)

}

func pageHandler(w http.ResponseWriter, r *http.Request, login string) {
	b, _ := checkLogin(login)

	if login != "" && b {
		p, _ := load(login)
		p.Info, _ = loadInfo(login)
		renderTemplate(w, "v4page", p)
		return
	}

	http.Redirect(w, r, "/v4hello/", http.StatusFound)

}

func brainHandler(w http.ResponseWriter, r *http.Request, login string) {

	Login1 := r.FormValue("login")
	Password := r.FormValue("password")
	p := &Person{Login: Login1, Password: Password, Info: ""}
	b, _ := checkPerson(p)
	if !b {

		err := p.save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//Login = "|"+ Login

		renderTemplate(w, "v4brain", p)

		return
	}

	http.Redirect(w, r, "/v4page/"+Login1, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			if (r.URL.Path)[len("/v4hello/"):] != "" && (r.URL.Path)[:len("/v4hello/")] == "/v4hello/" {
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
	http.HandleFunc("/v4brain/", makeHandler(brainHandler))
	http.HandleFunc("/v4page/", makeHandler(brainHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
