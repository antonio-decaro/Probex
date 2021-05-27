package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Persistence struct {
	logger *Logger
	svc    *dynamodb.DynamoDB
}

func InitPersistence(logger *Logger) (*Persistence, error) {
	ret := new(Persistence)

	ret.logger = logger
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
		p.logger.Error("Got error marshaling new item: " + err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = p.svc.PutItem(input)
	if err != nil {
		p.logger.Error("Got error calling PutItem: " + err.Error())
		return err
	}

	return nil
}
