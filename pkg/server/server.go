package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	done     chan os.Signal
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
		done:     make(chan os.Signal, 1),
	}
	dl.handlers["DynamoDB_20120810.CreateTable"] = dl.CreateTable
	dl.handlers["DynamoDB_20120810.ListTables"] = dl.ListTables
	dl.handlers["DynamoDB_20120810.PutItem"] = dl.PutItem
	// dl.handlers["DynamoDB_20120810.GetItem"]
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

func (dl *DynamoLocal) RunServer(addr string) error {
	router := mux.NewRouter()
	router.HandleFunc("/", dl.RunDynamo)

	signal.Notify(dl.done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := http.ListenAndServe(addr, router); err != nil {
			log.Fatal(err)
		}
	}()
	<-dl.done
	return nil
}

func (dl *DynamoLocal) Close() error {
	dl.done <- syscall.SIGTERM
	return nil
}
