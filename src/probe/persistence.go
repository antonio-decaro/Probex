package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Persistence struct {
	svc *dynamodb.DynamoDB
}

func InitPersistence() (*Persistence, error) {
	ret := new(Persistence)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ret.svc = dynamodb.New(sess)

	return ret, nil
}

func (p *Persistence) PersistProbeData(data ProbeData) error {

	tableName := "Planets"

	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("got error marshaling new item: %s", err))
	}

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(data.Name),
			},
		},
		TableName: aws.String(tableName),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":h": av["Humidity"],
			":w": av["Wind"],
			":t": av["Temperature"],
		},
		UpdateExpression: aws.String("add Humidity = :r, Wind = :w, Temperature = :t"),
	}

	_, err = p.svc.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("got error calling PutItem: %s", err)
	}

	return nil
}
