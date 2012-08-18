// Copyright 2012 Darren Elwood <darren@textnode.com> http://www.textnode.com @textnode
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"os"
	"github.com/textnode/xml2json"
)

func main() {
	log.Println("Started")

	var err error = nil

	in, err := os.Open("file.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	out, err := os.Create("file.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	defer out.Close()

	var x2j *xml2json.Xml2Json = xml2json.NewXml2Json("obfuscatedTextKey", "obfuscatedChildrenKey")
	err = x2j.Transform(in, out)

	log.Println("Finished: ", err)
}
