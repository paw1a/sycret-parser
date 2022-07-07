package doc

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"github.com/jmoiron/sqlx"
	"github.com/paw1a/sycret-parser/internal/db"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

var (
	ErrReadDocData = errors.New("failed to read doc data")
)

func GenerateDoc(docData []byte, recordID string, conn *sqlx.DB) ([]byte, error) {
	docTree := etree.NewDocument()
	if err := docTree.ReadFromBytes(docData); err != nil {
		return nil, ErrReadDocData
	}

	rootUse := docTree.FindElement("//use")
	if rootUse == nil {
		return nil, fmt.Errorf("no root use in document")
	}

	table := rootUse.SelectAttr("table")
	if table == nil {
		return nil, fmt.Errorf("no table attr in root USE tag")
	}

	queryString := fmt.Sprintf("select * from %s where id=?", table.Value)

	rows, err := conn.Query(queryString, recordID)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed select for root USE tag")
	}

	objects, err := db.ScanSelectRows(rows)
	if err != nil {
		return nil, fmt.Errorf("root USE tag: %v", err)
	}

	if len(objects) > 1 {
		return nil, fmt.Errorf("multiple root use tags")
	}

	err = recursiveParse(rootUse, objects[0], conn)
	if err != nil {
		return nil, err
	}

	return docTree.WriteToBytes()
}

//func GenerateDoc(docData []byte, recordID string, conn *sqlx.DB) ([]byte, error) {
//	docTree := etree.NewDocument()
//
//	if err := docTree.ReadFromBytes(docData); err != nil {
//		return nil, ErrReadDocData
//	}
//
//	for _, elem := range docTree.FindElements("//text") {
//		fieldValue := elem.SelectAttr("field").Value
//		textElem := elem.SelectElement("r").SelectElement("t")
//
//		newText, err := api.GetUserField(fieldValue, recordID)
//		if err != nil {
//			return nil, fmt.Errorf("failed to generate doc: %v", err)
//		}
//
//		textElem.SetText(fmt.Sprintf("%s ", newText))
//	}
//
//	return docTree.WriteToBytes()
//}

//func Generate(docData []byte, recordID string) ([]byte, error) {
//	docTree := etree.NewDocument()
//
//	if err := docTree.ReadFromBytes(docData); err != nil {
//		return nil, ErrReadDocData
//	}
//
//	rootUse := docTree.FindElement("//PurchaseOrder")
//	//rootTableName := rootUse.SelectAttr("table").Value
//	//fmt.Printf("%s\n", rootTableName)
//
//	//for _, elem := range rootUse.ChildElements() {
//	//	rec(elem, 0)
//	//}
//	rec(rootUse, 0)
//
//	return docTree.WriteToBytes()
//}
//
//func rec(rootElem *etree.Element, n int) {
//	for _, elem := range rootElem.ChildElements() {
//		fmt.Printf("%s%s\n", strings.Repeat(" ", n*4), elem.Tag)
//		rec(elem, n+1)
//	}
//}

func recursiveParse(rootElem *etree.Element, object map[string]interface{}, conn *sqlx.DB) error {
	for _, elem := range rootElem.ChildElements() {
		if elem.Tag == "text" {
			fieldName := elem.SelectAttr("field")
			if fieldName == nil {
				return fmt.Errorf("no field attr in TEXT tag")
			}

			value, ok := object[strings.ToUpper(fieldName.Value)]
			textElem := elem.SelectElement("r").SelectElement("t")

			if ok {
				sr := strings.NewReader(fmt.Sprintf("%v ", value))
				tr := transform.NewReader(sr, charmap.Windows1251.NewDecoder())
				buf, err := ioutil.ReadAll(tr)

				if err != nil {
					return fmt.Errorf("encoding error")
				}

				encodedString := string(buf)
				textElem.SetText(fmt.Sprintf("%s ", encodedString))
			} else {
				textElem.SetText("")
			}

			return nil
		} else if elem.Tag == "use" {
			table := elem.SelectAttr("table")
			query := elem.SelectAttr("query")

			if table == nil && query == nil {
				return fmt.Errorf("no table and query attr in USE tag")
			}

			var rows *sql.Rows
			var cols []string
			var err error

			if table != nil {
				idName := strings.ToUpper(table.Value) + "ID"
				queryString := fmt.Sprintf("select * from %s where id=?", table.Value)
				rows, err = conn.Query(queryString, object[idName])
			}

			if query != nil {
				queryString := query.Value
				rows, err = conn.Query(queryString)
			}

			if err != nil {
				return fmt.Errorf("failed select for USE tag")
			}

			cols, _ = rows.Columns()

			var newObject map[string]interface{}

			for rows.Next() {
				columns := make([]interface{}, len(cols))
				columnPointers := make([]interface{}, len(cols))
				for i, _ := range columns {
					columnPointers[i] = &columns[i]
				}

				if err := rows.Scan(columnPointers...); err != nil {
					return fmt.Errorf("failed to scan fields for USE tag")
				}

				newObject = make(map[string]interface{})
				for i, colName := range cols {
					val := columnPointers[i].(*interface{})
					newObject[colName] = *val
				}
			}

			err = recursiveParse(elem, newObject, conn)
			if err != nil {
				return err
			}
		} else if elem.Tag == "list" {
			table := elem.SelectAttr("table")
			fkey := elem.SelectAttr("fkey")

			if table == nil || fkey == nil {
				return fmt.Errorf("no table or fkey attr in LIST tag")
			}

			queryString := fmt.Sprintf("select * from %s where %s=?", table.Value, fkey.Value)
			rows, err := conn.Query(queryString, object["ID"])
			cols, _ := rows.Columns()

			if err != nil {
				return fmt.Errorf("failed select for LIST tag")
			}

			var newObject map[string]interface{}

			var counter int
			for rows.Next() {
				columns := make([]interface{}, len(cols))
				columnPointers := make([]interface{}, len(cols))
				for i, _ := range columns {
					columnPointers[i] = &columns[i]
				}

				if err := rows.Scan(columnPointers...); err != nil {
					return fmt.Errorf("failed to scan fields for LIST tag")
				}

				newObject = make(map[string]interface{})
				for i, colName := range cols {
					val := columnPointers[i].(*interface{})
					newObject[colName] = *val
				}

				fmt.Printf("%v\n", newObject)

				if counter > 0 {
					copyTag := elem.Copy()
					elem.Parent().AddChild(copyTag)
					err = recursiveParse(copyTag, newObject, conn)
				} else {
					err = recursiveParse(elem, newObject, conn)
				}

				if err != nil {
					return err
				}
				counter++
			}
		} else {
			err := recursiveParse(elem, object, conn)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
