package doc

import (
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"github.com/paw1a/sycret-parser/internal/api"
)

var (
	ErrReadDocData = errors.New("failed to read doc data")
)

func GenerateDoc(docData []byte, recordID string) ([]byte, error) {
	docTree := etree.NewDocument()

	if err := docTree.ReadFromBytes(docData); err != nil {
		return nil, ErrReadDocData
	}

	for _, elem := range docTree.FindElements("//text") {
		fieldValue := elem.SelectAttr("field").Value
		textElem := elem.SelectElement("r").SelectElement("t")

		newText, err := api.GetUserField(fieldValue, recordID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate doc: %v", err)
		}

		textElem.SetText(fmt.Sprintf("%s ", newText))
	}

	return docTree.WriteToBytes()
}
