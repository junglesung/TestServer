package main

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Book struct {
	Name       string    `json:"name"`
	Author     string    `json:"author"`
	Pages      int       `json:"pages"`
	Year       int       `json:"year"`
	CreateTime time.Time `json:"createtime"`
}

const BookKind = "Book"
const BookRoot = "Book Root"
const BookName = ["AAA", "BBB", "CCC", "DDD", "EEE", "FFF", "GGG", "HHH", "III", "JJJ"]
const BookAuthor = ["AuthorA", "AuthorB", "AuthorC", "AuthorD", "AuthorE", "AuthorF", "AuthorG", "AuthorH", "AuthorI", "AuthorJ"]
const BookMaxPages = 1000

func init() {
	rand.Seed(time.Now().Unix())
	http.HandleFunc("/queryAll", queryAll)
	http.HandleFunc("/storeAll", queryAll)
	http.HandleFunc("/deleteAll", queryAll)
}

func queryAll(rw http.ResponseWriter, req *http.Request) {
	//
}

func storeTen(rw http.ResponseWriter, req *http.Request) {
	//
	c := appengine.NewContext(req)
	pKey := datastore.NewKey(c, BookKind, BookRoot, 0, nil)
	for i := 0; i < 10; i++ {
		v := Book{
			Name: BookName[i], 
			Author: BookAuthor[i], 
			Pages: rand.Intn(BookMaxPages)
			Year: rand.Intn(time.Now().Year())
			CreatTime: time.Now()
		}
		if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, BookKind, pKey), &v); err != nil {
			//
		}
	}
}

func deleteAll(rw http.ResponseWriter, req *http.Request) {
	//
}
