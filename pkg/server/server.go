package server

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"net/http"
	"sync"
)

const (
	tableInfo = "DynamoLocalTableInfo"
)

type OperationHandler func(w http.ResponseWriter, req *http.Request)

type DynamoLocal struct {
	handlers map[string]OperationHandler
	db       *bolt.DB
	mu       sync.RWMutex
	keySep   []byte
}

func (dl *DynamoLocal) initTables() {
	dl.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucket([]byte(tableInfo)); err != nil {
			return err
		}
		return nil
	})
}

func NewDynamoLocal(filepath string, keySep []byte) (*DynamoLocal, error) {

	db, err := bolt.Open(filepath, 0600, nil)
	if err != nil {
		return nil, err
	}

	dl := DynamoLocal{
		handlers: map[string]OperationHandler{},
		db:       db,
		mu:       sync.RWMutex{},
	}
	dl.handlers["DynamoDB_20120810.CreateTable"] = dl.CreateTable
	dl.handlers["DynamoDB_20120810.ListTables"] = dl.ListTables
	dl.handlers["DynamoDB_20120810.PutItem"] = dl.PutItem
	// DynamoDB_20120810.ListTables - ListTables
	dl.initTables()
	return &dl, nil
}

func (dl *DynamoLocal) RunDynamo(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Host)
	fmt.Println(req.Header)
	fmt.Println(req.URL)
	fmt.Println(req.Method)

	method := req.Header.Get("X-Amz-Target")

	handler, ok := dl.handlers[method]
	if !ok {
		var result interface{}
		err := json.NewDecoder(req.Body).Decode(&result)

		/*
			var result types.CreateTablePayload
			err := json.NewDecoder(req.Body).Decode(&result)
		*/
		fmt.Println(err)
		fmt.Println(result)
		defer req.Body.Close()
		return
	}
	handler(w, req)
}
