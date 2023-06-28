package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func getAllHandle(c *gin.Context) {
	// Extract the key from the URL path

	readFile, err := os.Open("ledger.txt")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	count := 1
	for fileScanner.Scan() {
		var tmp Block
		err := json.Unmarshal([]byte(fileScanner.Text()), &tmp)

		if err != nil {
			fmt.Println(err)
		}

		c.JSON(http.StatusOK, tmp)
		fmt.Println(fileScanner.Text())
		count++
	}

	readFile.Close()

}
