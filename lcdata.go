package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

func pushlocalBlock(c *gin.Context) {
	db, err := leveldb.OpenFile("path/to/db", nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	for i := 0; i < 1000; i++ {
		var str string
		var str2 string
		str2 = strconv.Itoa(i)
		str = "SIM" + str2

		mycon := simcon{1, 1.0}
		bmc, err := json.Marshal(mycon)
		err = db.Put([]byte(str), []byte(bmc), nil)
		if err != nil {
			fmt.Println((err))
		} else {
			fmt.Println(str)
		}

	}

}
