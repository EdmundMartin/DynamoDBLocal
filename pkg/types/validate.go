package types

import (
	"bytes"
	"errors"
	"fmt"
)

func ValidatePutItem(tableDescription *TableDescriptionResponse, req *PutItemRequest) error {

	for _, element := range tableDescription.KeySchema {

		val, ok := req.Item[element.AttributeName]
		if !ok {
			// TODO - Work out actual error thrown by Dynamo
			return errors.New("missing key schema element")
		}
		fieldType, fieldVal := flattenItemMap(val)
		// TODO - validate field value
		fmt.Println(fieldVal)
		if err := checkAttributeDefinitions(tableDescription.AttributeDefinitions, element.AttributeName, fieldType); err != nil {
			return err
		}
	}
	return nil
}

func GetPutKey(tableDescription *TableDescriptionResponse, req *PutItemRequest, sep []byte) []byte {

	var hashKey []byte
	var sortKey []byte
	for _, element := range tableDescription.KeySchema {
		if element.KeyType == "HASH" {
			itemMapping := req.Item[element.AttributeName]
			_, fieldVal := flattenItemMap(itemMapping)
			hashKey = []byte(fieldVal)
		} else {
			itemMapping := req.Item[element.AttributeName]
			_, fieldVal := flattenItemMap(itemMapping)
			sortKey = []byte(fieldVal)
		}
	}
	if sortKey == nil {
		return hashKey
	}
	return bytes.Join([][]byte{
		hashKey, sortKey,
	}, sep)
}

func checkAttributeDefinitions(definitions []AttributeDefinition, fieldName, fieldType string) error {

	for _, def := range definitions {
		if def.AttributeName == fieldName {
			if def.AttributeType == fieldType {
				return nil
			}
		}
	}
	return errors.New("mismatch in attribute definition")
}

func flattenItemMap(itemAttr map[string]string) (string, string) {
	var types []string
	var values []string

	for k, v := range itemAttr {
		types = append(types, k)
		values = append(values, v)
	}

	if len(types) == 0 || len(values) == 0 {
		return "", ""
	}

	return types[0], values[0]
}
