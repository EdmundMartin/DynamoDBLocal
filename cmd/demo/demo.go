package main

import (
	"DynamoDBLocal/pkg/server"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

func main() {

	localDynamo, err := server.NewDynamoLocal("testdatabase.db", []byte("#"))
	if err != nil {
		log.Fatal(err)
	}
	go localDynamo.RunServer(":8080")
	defer localDynamo.Close()

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("localhost"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8080"}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "abcd", SecretAccessKey: "a1b2c3", SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)

	conn := dynamodb.NewFromConfig(cfg)
	res, _ := conn.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: aws.String("Hello"),
		KeySchema: []types.KeySchemaElement{ // key: year + title
			{
				AttributeName: aws.String("year"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("title"),
				KeyType:       types.KeyTypeRange,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("year"),
				AttributeType: types.ScalarAttributeTypeN, // data type descriptor: N == number
			},
			{
				AttributeName: aws.String("title"),
				AttributeType: types.ScalarAttributeTypeS, // data type descriptor: S == string
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})

	fmt.Println(res)
}
