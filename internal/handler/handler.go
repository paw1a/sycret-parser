package handler

import (
	"fmt"
	"github.com/beevik/etree"
	"net/http"
)

type SycretAPIRequest struct {
	Text     string `json:"text"`
	RecordID string `json:"recordid"`
}

type SycretAPIResponse struct {
	Result            int    `json:"result"`
	ResultDescription string `json:"resultdescription"`
	ResultData        string `json:"resultdata"`
}

type DocParserRequest struct {
	URLTemplate string `json:"URLTemplate"`
	RecordID    string `json:"RecordID"`
}

type DocParserResponse struct {
	URLWord string `json:"URLWord"`
}

func DocEndpoint(w http.ResponseWriter, r *http.Request) {
	doc := etree.NewDocument()

	if err := doc.ReadFromFile("forma_025u.doc"); err != nil {
		panic(err)
	}

	for _, elem := range doc.FindElements("//text") {
		attr := elem.SelectAttr("field").Value
		textElem := elem.SelectElement("r").SelectElement("t")
		textElem.SetText(fmt.Sprintf("%s ", attr))
	}
	doc.WriteToFile("out.doc")
}
