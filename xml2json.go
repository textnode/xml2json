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

package xml2json

import (
	"encoding/xml"
	"github.com/textnode/jsonStreamer"
	"io"
)

type Frame struct {
	collectedText []byte
}

func NewFrame() *Frame {
	return &Frame{collectedText: make([]byte, 0, 0)}
}

func (self *Frame) AddText(text []byte) {
	capacity := cap(self.collectedText)
	currentSize := len(self.collectedText)
	requiredSize := currentSize + len(text)
	if requiredSize > capacity {
		newCollectedText := make([]byte, requiredSize, requiredSize*2)
		copy(newCollectedText[0:currentSize], self.collectedText[0:currentSize])
		self.collectedText = newCollectedText
	}
	copy(self.collectedText[currentSize:requiredSize], text)
}

type Xml2Json struct {
	textKey     string
	childrenKey string
	frames      []*Frame
}

func NewXml2Json(textKey string, childrenKey string) *Xml2Json {
	var newXml2Json *Xml2Json = &Xml2Json{textKey: textKey, childrenKey: childrenKey, frames: make([]*Frame, 1, 10)}
	newXml2Json.frames[0] = NewFrame()
	return newXml2Json
}

func (self *Xml2Json) Transform(in io.Reader, out io.Writer) (err error) {
	var decoder *xml.Decoder = xml.NewDecoder(in)
	var encoder *jsonStreamer.JsonStreamer = jsonStreamer.NewJsonStreamer(out)

	var token xml.Token
	token, err = decoder.Token()

	for ; err == nil; token, err = decoder.Token() {
		var currentFrameIndex int = len(self.frames) - 1
		var currentFrame *Frame = self.frames[currentFrameIndex]

		switch specific := token.(type) {
		case xml.StartElement:
			err = encoder.StartObject()
			if err != nil {
				return
			}

			self.frames = append(self.frames, NewFrame())
			currentFrameIndex = len(self.frames) - 1
			currentFrame = self.frames[currentFrameIndex]

			err = encoder.WriteKey(specific.Name.Local)
			if err != nil {
				return
			}
			err = encoder.StartObject()
			if err != nil {
				return
			}

			for _, attr := range specific.Attr {
				err = encoder.WriteKey(attr.Name.Local)
				if err != nil {
					return
				}
				err = encoder.WriteStringValue(attr.Value)
				if err != nil {
					return
				}
			}

			err = encoder.WriteKey(self.childrenKey)
			if err != nil {
				return
			}
			err = encoder.StartArray()
			if err != nil {
				return
			}

		case xml.EndElement:
			err = encoder.EndArray() //close children
			if err != nil {
				return
			}

			if 0 < len(currentFrame.collectedText) {
				err = encoder.WriteKey(self.textKey)
				if err != nil {
					return
				}
				err = encoder.WriteStringValue(string(currentFrame.collectedText))
				if err != nil {
					return
				}
			}
			self.frames = self.frames[:len(self.frames)-1]
			err = encoder.EndObject()
			if err != nil {
				return
			}
			err = encoder.EndObject()
			if err != nil {
				return
			}
		case xml.CharData:
			currentFrame.AddText(specific)
			//case xml.Comment:
			//case xml.ProcInst:
			//case xml.Directive:
		}
	}
	return
}
