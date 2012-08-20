package main

import (
	"log"
	"os"
	"net/http"
	"github.com/textnode/xml2json"
)

func main() {
	log.Println("Started")

	var err error = nil

	resp, err := http.Get("http://search.twitter.com/search.rss?q=from%3Atextnode")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	out, err := os.Create("file.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	defer out.Close()

	var x2j *xml2json.Xml2Json = xml2json.NewXml2Json("obfuscatedTextKey", "obfuscatedChildrenKey")
	err = x2j.Transform(resp.Body, out)

	log.Println("Finished: ", err)
}
