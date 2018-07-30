package database

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
)

type Book struct {
	Title      string
	Author     Person
	Annotation template.HTML
	Body       template.HTML
	Code       string
	Bpage      int64
	Body1      string
}

type Person struct {
	Firstname string
	Lastname  string
}

type BookDB struct {
	Id          bson.ObjectId `bson:"_id"`
	Title       string        `bson:"title"`
	AuthorName  string        `bson:"authorName"`
	AuthorLName string        `bson:"authorLName"`
	Annotation  string        `bson:"annotation"`
	Body        string        `bson:"body"`
	Code        string        `bson:"code"`
}

type User struct {
	Id       bson.ObjectId `bson:"_id"`
	Login    string        `bson:"login"`
	Password string        `bson:"password"`
}

type Zakladka struct {
	Id    bson.ObjectId `bson:"_id"`
	User  string        `bson:"name"`
	Title string        `bson:"title"`
	Page  int           `bson:"page"`
}

func startDb(col string) mgo.Collection {
	session, _ := mgo.Dial("mongodb://127.0.0.1")
	mainDB := session.DB("library")
	return *mainDB.C(col)
}

func Load(md string, pNum int64) (Book, int, error) {
	session, _ := mgo.Dial("mongodb://127.0.0.1")
	defer session.Close()
	mainDB := session.DB("library")
	colBooks := mainDB.C("books")

	query := bson.M{
		"title": md,
	}

	library := []BookDB{}
	colBooks.Find(query).All(&library)
	author := Person{library[0].AuthorName, library[0].AuthorLName}
	book := Book{library[0].Title, author, template.HTML(library[0].Annotation), template.HTML(library[0].Body), library[0].Code, pNum, library[0].Body}

	return book, 200, nil

}

func LoadUser(md string) (*User, int, error) {
	session, _ := mgo.Dial("mongodb://127.0.0.1")
	defer session.Close()
	mainDB := session.DB("library")
	colBooks := mainDB.C("users")

	query := bson.M{
		"login": md,
	}

	library := []User{}
	colBooks.Find(query).All(&library)

	user := User{library[0].Id, library[0].Login, library[0].Password}

	return &user, 200, nil

}

func LoadFromDb() []string {
	session, _ := mgo.Dial("mongodb://127.0.0.1")
	defer session.Close()
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
func ExistDb(title string) bool {
	session, _ := mgo.Dial("mongodb://127.0.0.1")
	defer session.Close()
	mainDB := session.DB("library")
	colBooks := mainDB.C("books")
	query := bson.M{"title": title}
	N, _ := colBooks.Find(query).Count()
	return !(N == 0)

}

func UserExistDb(login string, password string) (int, error, error) {
	session, _ := mgo.Dial("mongodb://127.0.0.1")
	defer session.Close()
	mainDB := session.DB("library")
	colBooks := mainDB.C("users")
	query1 := bson.M{"login": login}
	query2 := bson.M{"login": login, "password": password}
	N, err := colBooks.Find(query1).Count()
	N1, err2 := colBooks.Find(query2).Count()
	p := 0
	if N > 0 {
		p = 1
	}
	if N1 > 0 {
		p = 2
	}
	return p, err, err2
}

func Unload(book *BookDB) error {
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

func UnloadUser(user *User) error {
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// получаем коллекцию
	productCollection := session.DB("library").C("users")
	err = productCollection.Insert(user)
	if err != nil {
		return err
	}
	return nil
}
