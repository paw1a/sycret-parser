package main

import (
	"fmt"
	_ "github.com/nakagami/firebirdsql"
	"github.com/paw1a/sycret-parser/internal/db"
	"github.com/paw1a/sycret-parser/internal/doc"
	"io"
	"os"
)

//func main() {
//	file, err := os.Open("example.xml")
//	if err != nil {
//		return
//	}
//
//	docData, err := io.ReadAll(file)
//	if err != nil {
//		return
//	}
//
//	_, err = doc.Generate(docData, "54511712")
//	//resultData, err := doc.Generate(docData, "2000", conn)
//	if err != nil {
//		fmt.Printf("%v", err)
//		os.Exit(10)
//	}
//}

func main() {
	file, err := os.Open("checkpretty.xml")
	if err != nil {
		return
	}

	docData, err := io.ReadAll(file)
	if err != nil {
		return
	}

	conn, err := db.NewDB(db.DBConnection{
		DBServer: "45.132.17.160",
		DBPath:   "C:/DemoServer/MOB_APP/beautyVialange.fdb",
	})
	defer conn.Close()

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	//resultData, err := doc.GenerateDoc(docData, "54511712", conn)
	resultData, err := doc.GenerateDoc(docData, "387", conn)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(10)
	}

	resultFile, err := os.Create("result.doc")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	_, err = resultFile.Write(resultData)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	resultFile, err = os.Create("result.xml")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	_, err = resultFile.Write(resultData)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
}
