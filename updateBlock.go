package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func createHandler(c *gin.Context) {
	// Extract the key from the URL path

	var sce []map[string]simcon

	if err := c.ShouldBindJSON(&sce); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wg sync.WaitGroup

	// wg.Add(1)
	// 2

	for i := 0; i < len(sce); i++ {

		// mu.Lock()
		// fmt.Println("hi")
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			Validator(sce[i])
			// putinblock(mtxn)
			fmt.Println(i)
		}(i)

		// defer mu.Unlock()
		// time.Sleep(time.Millisecond * time.Duration(1000))
	}
	wg.Wait()

}
