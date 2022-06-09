package doc

import (
	"fmt"
	"github.com/beevik/etree"
	"os"
)

func GenerateDoc(docFile *os.File) (*os.File, error) {
	docTree := etree.NewDocument()

	if _, err := docTree.ReadFrom(docFile); err != nil {
		panic(err)
	}

	for _, elem := range docTree.FindElements("//text") {
		attr := elem.SelectAttr("field").Value
		textElem := elem.SelectElement("r").SelectElement("t")
		textElem.SetText(fmt.Sprintf("%s ", attr))
	}
	docTree.WriteToFile("out.doc")

	return docFile, nil
}
