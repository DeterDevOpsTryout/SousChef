package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var db *sql.DB = nil

func checkQuantity(table string) int {

	q := fmt.Sprintf("SELECT count FROM %s", table)
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
		fmt.Println("the bistrofridge is empty!")
		fmt.Println(err)
	}

	return count

}

func restock(table string, count int) {
	fmt.Printf("bistro %s supplies low, restocking\n", table)
	cmd := fmt.Sprintf("UPDATE %s SET count = %d", table, count+200)
	_, err := db.Query(cmd)
	if err != nil {
		fmt.Println("cant put stuff in the bistrofridge")
		fmt.Println(err)
	}
}

func main() {

	fmt.Println("I am the souschef har de har har")

	var err error = nil

	for {

		db, err = sql.Open("mysql", "root:pizza@tcp(bistrofridge.deterlab.net:3306)/bistrofridge")
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

		if q := checkQuantity("ingredient_packs"); q < 100 {
			restock("ingredient_packs", q)
		}

		if q := checkQuantity("doughballs"); q < 100 {
			restock("doughballs", q)
		}

		time.Sleep(2 * time.Second)

	}

}
