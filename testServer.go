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
const BookMaxPages = 1000

var BookName = []string{"AAA", "BBB", "CCC", "DDD", "EEE", "FFF", "GGG", "HHH", "III", "JJJ"}
var BookAuthor = []string{"AuthorA", "AuthorB", "AuthorC", "AuthorD", "AuthorE", "AuthorF", "AuthorG", "AuthorH", "AuthorI", "AuthorJ"}

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
	// Get all entities
	var dst []Book
	c := appengine.NewContext(req)
	_, err := datastore.NewQuery(BookKind).Order("Pages").GetAll(c, &dst)
	if err != nil {
		log.Println(err)
	}

	// Return
	encoder := json.NewEncoder(rw)
	if err = encoder.Encode(dst); err != nil {
		log.Println(err, "in encoding result", dst)
	} else {
		log.Printf("QueryAll() returns %d items\n", len(dst))
	}

	// Status
	rw.WriteHeader(http.StatusOK)
}

func storeTen(rw http.ResponseWriter, req *http.Request) {
	// Store 10 random entities
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

	// Return
	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(r); err != nil {
		log.Println(err, "in encoding return code", r.ReturnCode)
	} else {
		log.Println("StoreTen() returns", r.ReturnCode)
	}

	// Status
	rw.WriteHeader(http.StatusCreated)
}

func deleteAll(rw http.ResponseWriter, req *http.Request) {
	// Delete root entity after other entities
	r := Result{0}
	c := appengine.NewContext(req)
	pKey := datastore.NewKey(c, BookKind, BookRoot, 0, nil)
	if keys, err := datastore.NewQuery(BookKind).KeysOnly().GetAll(c, nil); err != nil {
		log.Println(err)
		r.ReturnCode = 1
	} else if err := datastore.DeleteMulti(c, keys); err != nil {
		log.Println(err)
		r.ReturnCode = 1
	} else if err := datastore.Delete(c, pKey); err != nil {
		log.Println(err)
		r.ReturnCode = 1
	}

	// Return
	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(r); err != nil {
		log.Println(err, "in encoding return code", r.ReturnCode)
	} else {
		log.Println("DeleteAll() returns", r.ReturnCode)
	}

	// Status
	rw.WriteHeader(http.StatusOK)
}
