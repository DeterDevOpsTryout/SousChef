package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

func main() {

	fmt.Println("I am the souschef har de har har")

	for {

		db, err := sql.Open("mysql", "root:pizza@tcp(bistrofridge.deterlab.net:3306)/bistrofridge")
		if err != nil {
			fmt.Println("error opening connection to bistrofridge")
			fmt.Println(err)
			os.Exit(1)
		}

		err = db.Ping()
		if err != nil {
			fmt.Println("unable to ping bistrofridge")
			fmt.Println(err)
		}

		q := "SELECT count FROM ingredient_packs"
		rows, err := db.Query(q)
		if err != nil {
			fmt.Println("unable to query bistrofridge")
			fmt.Println(err)
		}

		if !rows.Next() {
			fmt.Println("the bistrofridge is broken :(")
		}

		var count int
		err = rows.Scan(&count)
		if err != nil {
			fmt.Println("rotten ingredients in bistrofridge")
			fmt.Println(err)
		}

		if count < 100 {
			fmt.Println("bistro ingredient supplies low, restocking")
			cmd := fmt.Sprintf("UPDATE ingredient_packs SET count = %d", count+200)
			_, err = db.Query(cmd)
			if err != nil {
				fmt.Println("cant put stuff in the bistrofridge")
				fmt.Println(err)
			}
		}

		time.Sleep(2 * time.Second)

	}

}
