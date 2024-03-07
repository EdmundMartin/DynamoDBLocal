package main

import (
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
	if err != nil {
		log.Fatal(err)
	}
	conn := dynamodb.NewFromConfig(cfg)
	fmt.Println(conn)

	res, err := conn.CreateTable(context.Background(), &dynamodb.CreateTableInput{
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
	fmt.Println(res.TableDescription.TableStatus)
	fmt.Println(err)

	tabOut, err := conn.ListTables(context.Background(), &dynamodb.ListTablesInput{
		ExclusiveStartTableName: aws.String(""),
		Limit:                   aws.Int32(10),
	})
	fmt.Println(tabOut.TableNames)
	fmt.Println(err)

	conn.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"year":  &types.AttributeValueMemberN{Value: "1991"},
			"title": &types.AttributeValueMemberS{Value: "EdmundMartin"},
		},
		TableName: aws.String("Hello"),
	})

	conn.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String("Hello"),
		Key: map[string]types.AttributeValue{
			"year":  &types.AttributeValueMemberN{Value: "1991"},
			"title": &types.AttributeValueMemberS{Value: "EdmundMartin"},
		},
	})
}
