package server

import (
	"DynamoDBLocal/pkg/types"
	"bytes"
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

//map[Item:map[title:map[S:EdmundMartin] year:map[N:1991]] TableName:Hello]

func (dl *DynamoLocal) PutItem(w http.ResponseWriter, req *http.Request) {
	var putItemReq types.PutItemRequest
	if err := json.NewDecoder(req.Body).Decode(&putItemReq); err != nil {
		JSONResponse(w, types.ErrorBadRequest{Message: "bad request"}, http.StatusBadRequest)
		return
	}
	fmt.Println(putItemReq)

	var tableDef types.TableDescriptionResponse

	err := dl.db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(tableInfo))

		contents := bucket.Get([]byte(putItemReq.TableName))

		return json.NewDecoder(bytes.NewBuffer(contents)).Decode(&tableDef)
	})

	if err != nil {
		// TODO check actual Dynamo behaviour
		JSONResponse(w, types.ErrorBadRequest{Message: "no such table"}, http.StatusBadRequest)
		return
	}

	err = types.ValidatePutItem(&tableDef, &putItemReq)
	if err != nil {
		// TODO check actual Dynamo behaviour
		JSONResponse(w, types.ErrorBadRequest{Message: "bad schema"}, http.StatusBadRequest)
		return
	}
	fmt.Println("Is valid req ", err)

	putKey := types.GetPutKey(&tableDef, &putItemReq, []byte("#"))

	buffer := bytes.Buffer{}
	err = json.NewEncoder(&buffer).Encode(putItemReq.Item)
	if err != nil {
		// TODO check actual Dynamo behaviour in case of internal error
		JSONResponse(w, types.ErrorBadRequest{Message: "bad request"}, http.StatusBadRequest)
		return
	}

	err = dl.db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(putItemReq.TableName))
		return bucket.Put(putKey, buffer.Bytes())
	})
	if err != nil {
		// TODO check actual Dynamo behaviour in case of internal error
		JSONResponse(w, types.ErrorBadRequest{Message: "bad request"}, http.StatusBadRequest)
		return
	}

	JSONResponse(w, types.PutItemResponse{ConsumedCapacity: types.NewDefaultWriteConsumedCapacity(putItemReq.TableName)}, http.StatusOK)
}
