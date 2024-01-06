package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"time"
)

func main() {

	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/fortest")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var insert = func(name string) int64 {
		exec, err := db.Exec("insert into u1 set name=?", name)
		if err != nil {
			panic(err)
		}
		id, _ := exec.LastInsertId()
		return id
	}

	for i := 10000; i < 20000; i++ {
		name := fmt.Sprintf("%d%s", i, getRandStr())
		fmt.Println(insert(name))
	}

}

var seed = []string{"a", "b", "c", "d"}

func getRandStr() string {
	n := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(seed))

	str := ""

	for i := 0; i <= n; i++ {
		str += seed[i]
	}
	return str
}
