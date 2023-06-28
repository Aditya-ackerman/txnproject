package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getHandler(c *gin.Context) {
	// Extract the key from the URL path
	key := c.Param("key")

	i, err := strconv.Atoi(key)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(key)
	// fmt.Println(i)

	count := 1

	readFile, err := os.Open("ledger.txt")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {

		// fmt.Println(count)

		if count == i {
			// fmt.Println("hi")

			var tmp Block
			err := json.Unmarshal([]byte(fileScanner.Text()), &tmp)

			if err != nil {
				fmt.Println(err)
			}

			c.JSON(http.StatusOK, tmp)

			fmt.Println(fileScanner.Text())
			// readFile.Close()
			// return
		}

		count++
	}

	readFile.Close()

}
