package server

import (
	"DynamoDBLocal/pkg/types"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func (dl *DynamoLocal) ListTables(w http.ResponseWriter, req *http.Request) {

	var tableListReq types.ListTableRequest
	if err := json.NewDecoder(req.Body).Decode(&tableListReq); err != nil {
		JSONResponse(w, types.ErrorBadRequest{Message: "bad request"}, http.StatusBadRequest)
		return
	}
	dl.mu.RLock()
	defer dl.mu.RUnlock()

	var tables []string
	var lastTable string
	err := dl.db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(tableInfo))

		if tableListReq.ExclusiveTableName != "" {
			cursor := bucket.Cursor()
			for k, _ := cursor.Seek([]byte(tableListReq.ExclusiveTableName)); k != nil; k, _ = cursor.Next() {
				tableName := string(k)
				tables = append(tables, tableName)
				lastTable = tableName
			}
		} else {
			if err := bucket.ForEach(func(k []byte, v []byte) error {
				tableName := string(k)
				tables = append(tables, tableName)
				lastTable = tableName
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		// TODO - Properly model internal server errors
		JSONResponse(w, types.ErrorBadRequest{Message: "bad request"}, http.StatusBadRequest)
	}

	JSONResponse(w, types.ListTableResponse{
		TableNames:             tables,
		LastEvaluatedTableName: lastTable,
	}, http.StatusOK)
}
