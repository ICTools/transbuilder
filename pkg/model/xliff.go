package model

import "encoding/xml"

type Xliff struct {
	XMLName xml.Name `xml:"xliff"`
	Version string   `xml:"version,attr"`
	File    []File   `xml:"file"`
}

type File struct {
	XMLName    xml.Name `xml:"file"`
	Original   string   `xml:"original,attr"`
	SourceLang string   `xml:"source-language,attr"`
	TargetLang string   `xml:"target-language,attr"`
	Body       Body     `xml:"body"`
}

type Body struct {
	XMLName   xml.Name    `xml:"body"`
	TransUnit []TransUnit `xml:"trans-unit"`
}

type TransUnit struct {
	XMLName xml.Name `xml:"trans-unit"`
	ID      string   `xml:"id,attr"`
	Source  string   `xml:"source"`
	Target  string   `xml:"target"`
}
