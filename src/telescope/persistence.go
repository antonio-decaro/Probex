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

func (p *Persistence) PersistTelescopeData(data TelescopeData) error {

	tableName := "Planets"

	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("got error marshaling new item: %s", err))
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = p.svc.PutItem(input)
	if err != nil {
		return fmt.Errorf("got error calling PutItem: %s", err)
	}

	return nil
}
