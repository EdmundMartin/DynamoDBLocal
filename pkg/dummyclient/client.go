package dummyclient

import (
	"DynamoDBLocal/pkg/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetCreateTableResponse(target string, payload types.CreateTablePayload) {

	cli := http.Client{}

	pay, _ := json.Marshal(payload)

	uri, _ := url.Parse(target)
	req := http.Request{
		Method: "POST",
		URL:    uri,
		Header: map[string][]string{
			"Authorization": {"AWS4-HMAC-SHA256 Credential=abcd/20240301/localhost/dynamodb/aws4_request",
				"SignedHeaders=accept-encoding;amz-sdk-invocation-id;amz-sdk-request;content-length;content-type;host;x-amz-date;x-amz-target",
				"Signature=3a6ed35a64651cbabecfbe30f619a3781a1e57d9fea3ef70dc81f1aa7a1a63b8"},
			"Accept-Encoding": {"identity"},
			"Content-Length":  {"237"},
			"X-Amz-Target":    {"DynamoDB_20120810.CreateTable"},
		},
		Body: ioutil.NopCloser(bytes.NewReader(pay)),
	}
	res, err := cli.Do(&req)
	fmt.Println(err)
	if err != nil {
		return
	}
	fmt.Println(res.Status)
	var contents interface{}

	json.NewDecoder(res.Body).Decode(&contents)
	fmt.Println(contents)
}
