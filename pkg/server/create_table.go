package server

import (
	"DynamoDBLocal/pkg/types"
	"bytes"
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"net/http"
	"time"
)

func (dl *DynamoLocal) CreateTable(w http.ResponseWriter, req *http.Request) {
	var result types.CreateTablePayload

	err := json.NewDecoder(req.Body).Decode(&result)
	defer req.Body.Close()
	if err != nil {
		JSONResponse(w, types.ErrorBadRequest{Message: "bad request"}, http.StatusBadRequest)
		return
	}

	internalRepr := types.TableDescriptionResponse{
		AttributeDefinitions: result.AttributeDefinitions,
		KeySchema:            result.KeySchema,
		CreationDateTime:     time.Now().Unix(),
		ItemCount:            0,
		TableArn:             fmt.Sprintf("TableArn:arn:aws:dynamodb:ddblocal:000000000000:table/%s", result.TableName),
		TableName:            result.TableName,
		TableStatus:          "ACTIVE",
		ProvisionedThroughput: types.ProvisionedThroughPutResponse{
			LastDecreaseDateTime:   0,
			LastIncreaseDateTime:   0,
			NumberOfDecreasesToday: 0,
			ReadCapacityUnits:      result.ProvisionedThroughput.ReadCapacityUnits,
			WriteCapacityUnits:     result.ProvisionedThroughput.WriteCapacityUnits,
		},
	}

	dl.mu.RLock()
	var value []byte
	err = dl.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tableInfo))
		value = bucket.Get([]byte(result.TableName))
		return nil
	})
	dl.mu.RUnlock()

	if value != nil {
		JSONResponse(w, types.ErrorBadRequest{Message: "table already exists"}, http.StatusBadRequest)
		return
	}

	dl.mu.Lock()
	defer dl.mu.Unlock()
	err = dl.db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(tableInfo))
		var contents bytes.Buffer
		if err := json.NewEncoder(&contents).Encode(internalRepr); err != nil {
			return err
		}
		if _, err = tx.CreateBucket([]byte(result.TableName)); err != nil {
			return err
		}

		return bucket.Put([]byte(result.TableName), contents.Bytes())
	})

	if err != nil {
		// TODO - Error with proper internal failure
		JSONResponse(w, types.ErrorBadRequest{Message: "bad request"}, http.StatusBadRequest)
		return
	}

	response := types.CreateTableResponse{TableDescription: internalRepr}

	JSONResponse(w, response, http.StatusOK)
}
