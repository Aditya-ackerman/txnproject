package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func showhandler(c *gin.Context) {
	// Extract the key from the URL path

	for i, v := range slicetestStrings {

		fmt.Println(i, " ", v)

	}

}
