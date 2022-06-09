package doc

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/paw1a/sycret-parser/internal/api"
	"time"
)

func GenerateDoc(docData []byte, recordID string) (string, error) {
	docTree := etree.NewDocument()

	if err := docTree.ReadFromBytes(docData); err != nil {
		panic(err)
	}

	for _, elem := range docTree.FindElements("//text") {
		attr := elem.SelectAttr("field").Value
		textElem := elem.SelectElement("r").SelectElement("t")
		newText, err := api.GetUserField(attr, recordID)
		if err != nil {
			return "", err
		}
		textElem.SetText(fmt.Sprintf("%s ", newText))
	}

	filename := time.Now().Format("2006-01-02 15-04-05")
	docTree.WriteToFile(fmt.Sprintf("%s.doc", filename))

	return filename, nil
}
