package course

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/gofrs/uuid"
)

// Course defines the Course model
type Course struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
}

// Put extracts the Course from JSON and writes it to DynamoDB
func Put(body string) (Course, error) {
	var sess *session.Session
	var tableName string

	local, err := strconv.ParseBool(os.Getenv("AWS_SAM_LOCAL"))
	if err != nil {
		return Course{}, err
	}
	// Create dynamo client object locally if running SAM CLI
	if local {
		sess = session.Must(session.NewSession(&aws.Config{
			Endpoint: aws.String("http://dynamodb:8000"),
		}))
		tableName = "courses"
	} else {
		sess = session.Must(session.NewSession())
		tableName = os.Getenv("COURSE_TABLE_NAME")
	}
	svc := dynamodb.New(sess)

	// Marshall the requrest body
	var course Course
	json.Unmarshal([]byte(body), &course)

	// Generate new UUID to store Course in case course doesn't have one
	if uuid.Must(uuid.FromString(course.ID)) == uuid.Nil {
		id, _ := uuid.NewV4()
		course.ID = id.String()
	}

	// Marshall the Item into a Map DynamoDB can deal with
	av, err := dynamodbattribute.MarshalMap(course)
	if err != nil {
		fmt.Println("Got error marshalling map:", err.Error())
		return course, err
	}

	// Create Item in table and return
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	return course, err
}

// Read gets the Course from DynamoDB
func Read(id string) (Course, error) {
	var sess *session.Session
	var tableName string

	local, err := strconv.ParseBool(os.Getenv("AWS_SAM_LOCAL"))
	if err != nil {
		return Course{}, err
	}
	// Create dynamo client object locally if running SAM CLI
	if local {
		sess = session.Must(session.NewSession(&aws.Config{
			Endpoint: aws.String("http://dynamodb:8000"),
		}))
		tableName = "courses"
	} else {
		sess = session.Must(session.NewSession())
		tableName = os.Getenv("COURSE_TABLE_NAME")
	}
	svc := dynamodb.New(sess)
	course := Course{}

	// Perform the query
	fmt.Println("Trying to read from table: ", "courses")
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return course, err
	}

	// Unmarshall the result in to an Item
	err = dynamodbattribute.UnmarshalMap(result.Item, &course)
	if err != nil {
		fmt.Println(err.Error())
		return course, err
	}

	return course, nil
}

// Delete erases the Course from DynamoDB
func Delete(id string) error {
	var sess *session.Session
	var tableName string

	local, err := strconv.ParseBool(os.Getenv("AWS_SAM_LOCAL"))
	if err != nil {
		return err
	}
	// Create dynamo client object locally if running SAM CLI
	if local {
		sess = session.Must(session.NewSession(&aws.Config{
			Endpoint: aws.String("http://dynamodb:8000"),
		}))
		tableName = "courses"
	} else {
		sess = session.Must(session.NewSession())
		tableName = os.Getenv("COURSE_TABLE_NAME")
	}
	svc := dynamodb.New(sess)

	// Perform the delete
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err = svc.DeleteItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// List returns the Courses from DynamoDB
func List() ([]Course, error) {
	var sess *session.Session
	var tableName string

	local, err := strconv.ParseBool(os.Getenv("AWS_SAM_LOCAL"))
	if err != nil {
		return []Course{}, err
	}
	// Create dynamo client object locally if running SAM CLI
	if local {
		sess = session.Must(session.NewSession(&aws.Config{
			Endpoint: aws.String("http://dynamodb:8000"),
		}))
		tableName = "courses"
	} else {
		sess = session.Must(session.NewSession())
		tableName = os.Getenv("COURSE_TABLE_NAME")
	}
	svc := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var courses []Course
	dynamodbattribute.UnmarshalListOfMaps(result.Items, &courses)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return courses, nil
}
