package main

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

type getJishoApi struct {
	client c
}

type c struct {
	http.Client
}

var pingedTimesTracker = 0

const address = "127.0.0.1:8082"

func (api *getJishoApi) search(keyword string) (*http.Response, error) {
	urlx := "https://jisho.org/api/v1/search/words?keyword"
	response, e := api.client.Get(strings.Join([]string{urlx, url.QueryEscape(keyword)}, "="))
	if e != nil {
		return nil, e
	}
	return response, nil
}

func getIndex(wr http.ResponseWriter, r *http.Request) {
	http.Redirect(wr, r, "/kanji/1", 302)
	//-------OLD ROUTE----------
	/*t, e := template.ParseFiles("index.html")
	if e != nil {
		log.Println(e.Error())
	}
	t.Execute(wr, nil)*/
	return
}

func getDefinition(wr http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		wr.WriteHeader(http.StatusForbidden)
		fmt.Fprint(wr, "Forbidden!")
		return
	}
	api := &getJishoApi{}
	url := strings.SplitAfter(r.RequestURI, "=")
	response, e := api.search(url[1])
	b, e := ioutil.ReadAll(response.Body)
	if e != nil {
		log.Println(e.Error())
	}
	fmt.Println(string(b))
	wr.Header().Set("Content-Type", "application/json")
	wr.Write(b)
}

func getKanji(wr http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		wr.WriteHeader(http.StatusForbidden)
		fmt.Fprint(wr, "Forbidden!")
		return
	}
	url := strings.SplitAfter(r.RequestURI, "/kanji/")
	t, e := template.ParseFiles("kanji.html")
	if e != nil {
		log.Println(e.Error())
	}
	i, e := strconv.Atoi(url[1])
	if e != nil {
		log.Println(e.Error())
	}
	t.Execute(wr, dump.Entry[i])
}

func listen() *http.Server {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", getIndex)
	http.HandleFunc("/kanji/", getKanji)
	http.HandleFunc("/search", getDefinition)
	server := &http.Server{}
	go func() {
		log.Println(http.ListenAndServe(address, server.Handler))
	}()
	return server
}

func main() {
	loadDefinitions()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	t := time.Second * 90
	s := listen()
	ctx, cfunc := context.WithTimeout(context.Background(), t)
	defer cfunc()
	<-sig
	log.Println("Server shutting down.")
	log.Println(s.Shutdown(ctx))
}
