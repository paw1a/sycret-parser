package app

import (
	"fmt"
	"github.com/beevik/etree"
)

func Run() {
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
