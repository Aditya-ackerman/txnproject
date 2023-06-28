package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

type Users struct {
	Users []mySim `json:"users"`
}

type Barrs struct {
	Barrs []Block `json:"barrs"`
}

type BlockChain struct {
	blocks []*Block
}

type Block struct {
	Hash        []byte
	Blocknumber int    `json:"blocknumber"`
	Txns        []txn  `json:"trnx"`
	BlockStatus string `json:"blockStatus"`
	PrevHash    []byte
}

type txn struct {
	ID     string  `json:"id"`
	Val    int32   `json:"val"`
	Ver    float64 `json:"ver"`
	MyHash []byte  `json: "mhash"`
	Valid  bool    `json:"valid"`
}

type mySim struct {
	Key string  `json:"key"`
	Val int32   `json:"val"`
	Ver float64 `json:"ver"`
}

type simcon struct {
	Val int32   `json:"val"`
	Ver float64 `json:"ver"`
}

var marr [3]txn
var sz int
var curBlock int

var mgch chan txn
var mblockChan chan Block
var prevHash []byte
var wg sync.WaitGroup
var mu sync.Mutex
var mschan chan mySim
var kch chan txn

var test int
var test1 int
var tspent int
var ttime int

var slicetestStrings []time.Duration

var interval = time.Duration(1) * time.Second

func pushBlock(inputTxn map[string]simcon) {

	mu.Lock()
	defer mu.Unlock()
	m := mySim{}

	for id, value := range inputTxn {

		m.Key = id
		m.Val = value.Val
		m.Ver = value.Ver

	}

	key := m.Key
	// fmt.Println(key)

	db, err := leveldb.OpenFile("path/to/db", nil)

	if err != nil {
		fmt.Println("I can not open dbfile")
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("I can  open dbfile")

	data, err := db.Get([]byte(key), nil)

	if err != nil {
		fmt.Println("I can not open dbfile")
		log.Fatal(err)
	}

	var mdata simcon

	err = json.Unmarshal(data, &mdata)

	

	flag := false

	if mdata.Ver == m.Ver {
		mdata.Ver = m.Ver + 1.0
		mdata.Val = m.Val
		flag = true
	}

	info := bytes.Join([][]byte{[]byte(key)}, []byte{})
	hash := sha256.Sum256(info)
	// b.Hash = hash[:]

	newtxn := txn{key, mdata.Val, mdata.Ver, hash[:], flag}
	// fmt.Println(newtxn.myHash)
	marr[sz] = newtxn

	bmc, err := json.Marshal(mdata)

	err = db.Put([]byte(key), []byte(bmc), nil)

	// mgch <- newtxn

	// fmt.Println("h112")

	select {
	case mgch <- newtxn:
	default:
		// mu.Lock()

		// fmt.Println("hi mycont is full")
		// CreateBlock(prevHash)
		// mu.Unlock()

	}

	if len(mgch) == cap(mgch) {
		// fmt.Println("buff is full")
		CreateBlock()
	}

}

func appBlock() {
	// fmt.Println("hi")
	// mu.Lock()
	// defer mu.Unlock()

	dur := time.Since(now)

	slicetestStrings = append(slicetestStrings, dur)

	block := <-mblockChan
	// fmt.Print(block)
	data, _ := json.Marshal(block)

	file, err := os.OpenFile("ledger.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	_, _ = datawriter.WriteString(string(data) + "\n")

	datawriter.Flush()
	file.Close()

	
}

func CreateBlock() {

	fmt.Println("myhi createBlock")

	// mu.Lock()

	// fmt.Println(prevHash)
	// fmt.Println(curBlock)

	msice := []txn{}

	now = time.Now()

	// for result := range mgch {
	// 	fmt.Println((result))
	// 	msice = append(msice, result)
	// 	fmt.Println("err is here at part")
	// }

	for len(mgch) > 0 {
		result := <-mgch
		// fmt.Println(result)
		msice = append(msice, result)
		// fmt.Println("err is here at part")
	}

	// mbdata := blockdata{int32(curBlock), []byte{}, marr, "commit"}
	curBlock++
	// fmt.Println("err is here at part ", curBlock)

	block := Block{[]byte{}, int(curBlock), []txn(msice), "commit", prevHash}
	// fmt.Println(block)
	em, err := json.Marshal(block)

	if err != nil {
		log.Fatal(err)
	}

	info := bytes.Join([][]byte{em, block.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	block.Hash = hash[:]
	prevHash = hash[:]
	// fmt.Println(block)
	// block.DeriveHash()

	select {
	case mblockChan <- block:
	default:
		// fmt.Println("hi block is full")
	}

	if len(mblockChan) == cap(mblockChan) {
		// fmt.Println("buff is full")
		appBlock()
	}

	// mu.Unlock()

}



var now = time.Now()

func myfunc() {
	for {
		select {
		case v := <-kch:
			mgch <- v

			if len(mgch) == cap(mgch) {
				// fmt.Println("buff is full")
				CreateBlock()
				tspent = 0
				// timer2.Reset(delay1)
			}
		case <-time.After(5 * time.Second):
			if len(mgch) > 0 {
				CreateBlock()
			}

		}
	}
}

func main() {

	test1 = 3

	

	mgch = make(chan txn, 3)
	kch = make(chan txn, 1001)
	mblockChan = make(chan Block, 1)
	mschan = make(chan mySim, 1)

	

	jsonFile, err := os.Open("data.json")

	if err != nil {
		fmt.Println(err)
	}

	
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	
	var users Users

	
	json.Unmarshal(byteValue, &users)

	

	sz = 0
	curBlock = 0
	tspent = 0
	ttime = 3

	

	

	router := gin.Default()

	go myfunc()

	// POST /create/{key}
	router.POST("/localdata", pushlocalBlock)
	router.POST("/updatedata", createHandler)
	router.GET("/getexe", showhandler)

	// GET /get/{key}
	router.GET("/getu/:key", getHandler)

	router.GET("/get", getAllHandle)

	//router.GET("/delete/:key", deleteHandler)
	// Start the server on port 8080
	log.Println("Server listening on port 18080...")
	router.Run(":18080")

}
