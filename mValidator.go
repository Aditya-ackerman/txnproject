package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

func increment() {
	tspent = tspent%360 + 1
}

func Validator(inputTxn map[string]simcon) txn {

	mu.Lock()
	defer mu.Unlock()

	// fmt.Println("hi")

	m := mySim{}

	for id, value := range inputTxn {

		m.Key = id
		m.Val = value.Val
		m.Ver = value.Ver

	}

	key := m.Key
	fmt.Println(key)

	db, err := leveldb.OpenFile("path/to/db", nil)
	defer db.Close()

	if err != nil {

		log.Fatal(err)
	}

	// fmt.Println("I can  open dbfile")

	data, err := db.Get([]byte(key), nil)

	if err != nil {
		// fmt.Println("I can not open dbfile")
		log.Fatal(err)
	}

	var mdata simcon

	err = json.Unmarshal(data, &mdata)

	// fmt.Println("%v", mdata)

	// fmt.Println("h10")

	// fmt.Print(m.Val)
	// fmt.Print(" ")
	// fmt.Print(m.Ver)
	// fmt.Print(" ")
	// fmt.Print(mdata.Val)
	// fmt.Print(" ")
	// fmt.Println(mdata.Ver)

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

	kch <- newtxn
	return newtxn

}
