package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var dump entryDump

type KanjiFile struct {
	Kanji []string
}

type senses struct {
	EnDef         []string `json:"english_definitions"`
	PartsOfSpeech []string `json:"parts_of_speech"`
	Info          []string `json:"info"`
	Restrictions  []string `json:"restrictions"`
	SeeAlso       []string `json:"see_also"`
}

type jp struct {
	Word    string `json:"word"`
	Reading string `json:"reading"`
}

// actually only index 0 of data array.
type data struct {
	Slug     string   `json:"slug"`
	JLPT     []string `json:"jlpt"`
	Japanese []jp     `json:"japanese"`
	Senses   []senses `json:"senses"`
}
type meta struct {
	Status int `json:"meta"`
}

type definition struct {
	Meta meta   `json:"meta"`
	Data []data `json:"data"`
}

type entry struct {
	Kanji      string
	Definition definition
}
type entryDump struct {
	Entry []entry
}

func migrate() {
	source, e := os.Open("assets/aozora.json")
	defer source.Close()
	if e != nil {
		log.Println("Error opening file 1", e.Error())
	}
	b, e := ioutil.ReadAll(source)
	if e != nil {
		log.Println("Error reading from source", e.Error())
	}
	var k KanjiFile
	if e := json.Unmarshal(b, &k.Kanji); e != nil {
		log.Println("Error unmarshaling kanji", e.Error())
	}
	definitions, e := os.OpenFile("definitions.json", 1, os.ModeAppend)
	defer definitions.Close()
	if e != nil {
		log.Println("Error opening file 2", e.Error())
	}
	api := &getJishoApi{}
	/***Start loop***/
	var dump entryDump
	for _, v := range k.Kanji {
		log.Println(v)
		var ek entry
		response, e := api.search(v)
		b, e := ioutil.ReadAll(response.Body)
		if e != nil {
			log.Println("Search error >>>", e.Error())
		}
		fmt.Println(string(b))
		ek.Kanji = v
		if e := json.Unmarshal(b, &ek.Definition); e != nil {
			log.Println("Error unmarshaling definition", e.Error())
		}
		dump.Entry = append(dump.Entry, ek)
	}
	b, e = json.Marshal(&dump)
	if e != nil {
		log.Println(e.Error())
	}
	if _, e := definitions.Write(b); e != nil {
		log.Println(e.Error())
	}
	definitions.Close()
	source.Close()
}

func loadDefinitions() {
	f, e := os.Open("definitions.json")
	if e != nil {
		log.Println(e.Error())
	}
	b, e := ioutil.ReadAll(f)
	if e != nil {
		log.Println(e.Error())
		return
	}

	if e = json.Unmarshal(b, &dump); e != nil {
		log.Println(e.Error())
		return
	}
	log.Println("Kanji definitions loaded.")
	//log.Println("Testing:", dump.Entry[1500].Definition.Data[0].Japanese[0].Reading)
}

/*
func main() {
	//migrate()

		api := &getJishoApi{}
		r, e := api.search("ー")
		if e != nil {
			log.Println("Search error >>>", e.Error())
		}
		b, e := ioutil.ReadAll(r.Body)
		if e != nil {
			log.Println("Read error >>>", e.Error())
		}
		fmt.Println(string(b))

}
*/
