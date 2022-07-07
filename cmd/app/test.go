package main

import (
	"bytes"
	"fmt"
	_ "github.com/nakagami/firebirdsql"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

//func main() {
//	conn, err := db.NewDB(db.DBConnection{
//		DBServer: "45.132.17.160",
//		DBPath:   "C:/DemoServer/MOB_APP/beautyVialange.fdb",
//	})
//	defer conn.Close()
//
//	if err != nil {
//		fmt.Printf("%v", err)
//		return
//	}
//
//	//row, err := conn.Query("select id from visit where id=2000")
//	//row.Next()
//	//var id int
//	//row.Scan(&id)
//	//fmt.Printf("%d\n", id)
//
//	//row, err := conn.Query("select * from visit")
//
//	//	rows, _ := conn.Query("select * from visit where id=2000")
//	//	cols, _ := rows.Columns()
//	//
//	//	var m map[string]interface{}
//	//
//	//	for rows.Next() {
//	//		columns := make([]interface{}, len(cols))
//	//		columnPointers := make([]interface{}, len(cols))
//	//		for i, _ := range columns {
//	//			columnPointers[i] = &columns[i]
//	//		}
//	//
//	//		if err := rows.Scan(columnPointers...); err != nil {
//	//			fmt.Printf("%v", err)
//	//			return
//	//		}
//	//
//	//		m = make(map[string]interface{})
//	//		for i, colName := range cols {
//	//			val := columnPointers[i].(*interface{})
//	//			m[colName] = *val
//	//		}
//	//
//	//		fmt.Print(m["CLIENTID"])
//	//	}
//	//
//	//	//var visit []map[string]interface{}
//	//	//
//	//	//err = conn.Select(&visit, "select * from visit where id=10")
//	//	//if err != nil {
//	//	//	fmt.Printf("%v", err)
//	//	//	return
//	//	//}
//	//	//
//	//	//fmt.Printf("%v", visit)
//	//
//	//	rows, err = conn.Query("select fio from client where id=?", m["CLIENTID"])
//	//	rows.Next()
//	//	var fio string
//	//	rows.Scan(&fio)
//	//	fmt.Printf("%s", fio)
//	//	//row.Next()
//	//	//var id int
//	//	//row.Scan(&id)
//	//	//fmt.Printf("%d\n", id)
//	//	//var id int
//	//	//var name string
//	//	//for row.Next() {
//	//	//	row.Scan(&id, &name)
//	//	//	fmt.Printf("%d %s\n", id, name)
//	//	//}
//	//
//	//	//cols, _ := row.Columns()
//	//	//fmt.Printf("%v", cols)
//	//}
//	//
//	//func main() {
//	//	use  "select * from visit where id=[recordid]"
//	//	list "select * from appointment where visitid=?", visit["ID"]
//	//	list "select * from appointmentmanipulation where appointmentid=?", appointment["ID"]
//	//	use  "select * from manipulation where id=?", appointmentmanipulation["MANIPULATIONID"]
//	//	text manipulation["name"]
//	//
//	//	use  "select * from visit where id=[recordid]"
//	//	use  "select * from client where id=?", visit["CLIENTID"]
//	//	text client["fio"]
//}

func postFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	w, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(w, fh)
	if err != nil {
		return err
	}

	bodyWriter.Close()

	req, err := http.NewRequest("POST", targetUrl, bodyBuf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	req.Header.Add("User-Agent", "sycret handler user")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resp_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(res.Status)
	fmt.Println(string(resp_body))
	return nil
}

// sample usage
func main() {
	target_url := "https://online.sycretreg.ru/gendoc/result/"
	filename := "./check.doc"
	postFile(filename, target_url)
}
