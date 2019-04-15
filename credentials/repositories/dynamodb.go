package repositories

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/mitchell/selfpass/credentials/types"
)

// NewDynamoTable TODO
func NewDynamoTable(name string) DynamoTable {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic(err.Error())
	}
	return DynamoTable{
		name: name,
		svc:  dynamodb.New(cfg),
	}
}

// DynamoTable TODO
type DynamoTable struct {
	name string
	svc  *dynamodb.DynamoDB
}

// GetAllMetadata TODO
func (t DynamoTable) GetAllMetadata(ctx context.Context, sourceService string, errch chan<- error) (output <-chan types.Metadata) {
	mdch := make(chan types.Metadata, 1)
	in := &dynamodb.ScanInput{TableName: &t.name}

	if sourceService != "" {
		filterExpr := "SourceHost = :sh"
		in.FilterExpression = &filterExpr
		in.ExpressionAttributeValues = map[string]dynamodb.AttributeValue{
			":sh": {S: &sourceService},
		}
	}

	req := t.svc.ScanRequest(in)

	go func() {
		defer close(mdch)

		pgr := req.Paginate()
		for pgr.Next() {
			mds := []types.Metadata{}
			out := pgr.CurrentPage()
			if err := dynamodbattribute.UnmarshalListOfMaps(out.Items, &mds); err != nil {
				errch <- err
				return
			}

			for _, md := range mds {
				mdch <- md
			}
		}

		if err := pgr.Err(); err != nil {
			errch <- err
			return
		}
	}()

	return mdch
}

// Get TODO
func (t DynamoTable) Get(ctx context.Context, id string) (output types.Credential, err error) {
	req := t.svc.GetItemRequest(&dynamodb.GetItemInput{
		TableName: &t.name,
		Key: map[string]dynamodb.AttributeValue{
			"ID": {S: &id},
		},
	})

	out, err := req.Send()
	if err != nil {
		return output, err
	}

	err = dynamodbattribute.UnmarshalMap(out.Item, &output)
	return output, err
}

// Put TODO
func (t DynamoTable) Put(ctx context.Context, c types.Credential) (err error) {
	item, err := dynamodbattribute.MarshalMap(c)
	if err != nil {
		return err
	}

	req := t.svc.PutItemRequest(&dynamodb.PutItemInput{
		TableName: &t.name,
		Item:      item,
	})
	req.SetContext(ctx)

	_, err = req.Send()

	return err
}

// Delete TODO
func (t DynamoTable) Delete(ctx context.Context, id string) (err error) {
	req := t.svc.DeleteItemRequest(&dynamodb.DeleteItemInput{
		TableName: &t.name,
		Key: map[string]dynamodb.AttributeValue{
			"ID": {S: &id},
		},
	})

	_, err = req.Send()
	return err
}
