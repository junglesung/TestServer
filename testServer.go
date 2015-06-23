package main

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Result struct {
	ReturnCode int `json:"returncode"`
}

type Book struct {
	Name       string    `json:"name"`
	Author     string    `json:"author"`
	Pages      int       `json:"pages"`
	Year       int       `json:"year"`
	CreateTime time.Time `json:"createtime"`
}

const BookKind = "Book"
const BookRoot = "Book Root"

var BookName = []string{"AAA", "BBB", "CCC", "DDD", "EEE", "FFF", "GGG", "HHH", "III", "JJJ"}
var BookAuthor = []string{"AuthorA", "AuthorB", "AuthorC", "AuthorD", "AuthorE", "AuthorF", "AuthorG", "AuthorH", "AuthorI", "AuthorJ"}

const BookMaxPages = 1000

func init() {
	rand.Seed(time.Now().Unix())
	http.HandleFunc("/", rootPage)
	http.HandleFunc("/queryAll", queryAll)
	http.HandleFunc("/storeTen", storeTen)
	http.HandleFunc("/deleteAll", deleteAll)
}

func rootPage(rw http.ResponseWriter, req *http.Request) {
	//
}

func queryAll(rw http.ResponseWriter, req *http.Request) {
	//
}

func storeTen(rw http.ResponseWriter, req *http.Request) {
	r := Result{0}
	c := appengine.NewContext(req)
	pKey := datastore.NewKey(c, BookKind, BookRoot, 0, nil)
	for i := 0; i < 10; i++ {
		v := Book{
			Name:       BookName[i],
			Author:     BookAuthor[i],
			Pages:      rand.Intn(BookMaxPages),
			Year:       rand.Intn(time.Now().Year()),
			CreateTime: time.Now(),
		}
		if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, BookKind, pKey), &v); err != nil {
			log.Println(err)
			r.ReturnCode = 1
			break
		}
	}
	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(r); err != nil {
		log.Println(err, "in encoding return code", r.ReturnCode)
	} else {
		log.Println("StoreTen() returns", r.ReturnCode)
	}
}

func deleteAll(rw http.ResponseWriter, req *http.Request) {
	//
}
