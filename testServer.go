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

// type Result struct {
// 	ReturnCode int `json:"returncode"`
// }

type Book struct {
	Name       string    `json:"name"`
	Author     string    `json:"author"`
	Pages      int       `json:"pages"`
	Year       int       `json:"year"`
	CreateTime time.Time `json:"createtime"`
}

const BaseUrl = "/api/0.1/"
const BookKind = "Book"
const BookRoot = "Book Root"
const BookMaxPages = 1000

var BookName = []string{"AAA", "BBB", "CCC", "DDD", "EEE", "FFF", "GGG", "HHH", "III", "JJJ"}
var BookAuthor = []string{"AuthorA", "AuthorB", "AuthorC", "AuthorD", "AuthorE", "AuthorF", "AuthorG", "AuthorH", "AuthorI", "AuthorJ"}

func init() {
	rand.Seed(time.Now().Unix())
	http.HandleFunc(BaseUrl, rootPage)
	http.HandleFunc(BaseUrl+"queryAll", queryAll)
	http.HandleFunc(BaseUrl+"storeTen", storeTen)
	http.HandleFunc(BaseUrl+"deleteAll", deleteAll)
	http.HandleFunc(BaseUrl+"books", books)
}

func rootPage(rw http.ResponseWriter, req *http.Request) {
	//
}

func books(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		queryAll(rw, req)
	case "POST":
		storeTen(rw, req)
	case "DELETE":
		deleteAll(rw, req)
	default:
		queryAll(rw, req)
	}
}

func queryAll(rw http.ResponseWriter, req *http.Request) {
	// Get all entities
	var dst []Book
	r := 0
	c := appengine.NewContext(req)
	_, err := datastore.NewQuery(BookKind).Order("Pages").GetAll(c, &dst)
	if err != nil {
		log.Println(err)
		r = 1
	}

	// Return status. WriteHeader() must be called before call to Write
	if r == 0 {
		rw.WriteHeader(http.StatusOK)
	} else {
		http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Return body
	encoder := json.NewEncoder(rw)
	if err = encoder.Encode(dst); err != nil {
		log.Println(err, "in encoding result", dst)
	} else {
		log.Printf("QueryAll() returns %d items\n", len(dst))
	}
}

func storeTen(rw http.ResponseWriter, req *http.Request) {
	// Store 10 random entities
	r := 0
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
			r = 1
			break
		}
	}

	// Return status. WriteHeader() must be called before call to Write
	if r == 0 {
		rw.WriteHeader(http.StatusCreated)
	} else {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func deleteAll(rw http.ResponseWriter, req *http.Request) {
	// Delete root entity after other entities
	r := 0
	c := appengine.NewContext(req)
	pKey := datastore.NewKey(c, BookKind, BookRoot, 0, nil)
	if keys, err := datastore.NewQuery(BookKind).KeysOnly().GetAll(c, nil); err != nil {
		log.Println(err)
		r = 1
	} else if err := datastore.DeleteMulti(c, keys); err != nil {
		log.Println(err)
		r = 1
	} else if err := datastore.Delete(c, pKey); err != nil {
		log.Println(err)
		r = 1
	}

	// Return status. WriteHeader() must be called before call to Write
	if r == 0 {
		rw.WriteHeader(http.StatusOK)
	} else {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
