package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
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
		rows.Close()
		return -1
	}

	var count int
	err = rows.Scan(&count)
	if err != nil {
		fmt.Println("the bistrofridge is empty!")
		fmt.Println(err)
	}

	rows.Close()

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
	db, err = sql.Open("mysql",
		"root:pizza@tcp(bistrofridge.deterlab.net:3306)/bistrofridge")
	if err != nil {
		fmt.Println("error opening connection to bistrofridge")
		fmt.Println(err)
		os.Exit(1)
	}

	for {

		err = db.Ping()
		if err != nil {
			fmt.Println("unable to ping bistrofridge")
			fmt.Println(err)

			fmt.Println("trying to reconnect")
			db, err = sql.Open("mysql",
				"root:pizza@tcp(bistrofridge.deterlab.net:3306)/bistrofridge")
			if err != nil {
				fmt.Println("error opening connection to bistrofridge")
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if q := checkQuantity("ingredient_packs"); q < 100 {
			restock("ingredient_packs", q)
		}

		if q := checkQuantity("doughballs"); q < 100 {
			resp, err := http.Get("http://mixer.deterlab.net:8085/mix")
			if err != nil {
				fmt.Println("unable to make dough :(")
				fmt.Println(err)
			} else {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("the dough has bugs")
					fmt.Println(err)
				} else if string(body) != "dough" {
					fmt.Println("the dough is not dough")
					fmt.Println(err)
				} else {
					restock("doughballs", q)
				}
			}
		}

		time.Sleep(2 * time.Second)

	}

}
