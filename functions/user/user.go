package user

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"time"

	"github.com/gofrs/uuid"
)

// User defines the User model
type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	IsActive bool      `json:"is_active"`
	Email    string    `json:"email"`

	Address     Address       `json:"address"`
	Enrollments []Enrollments `json:"enrollments"`
}

// Address defines the Address model
type Address struct {
	Street string `json:"street"`
	Zip    string `json:"zip"`
	City   string `json:"city"`
}

// Enrollments defines the Enrollments model
type Enrollments struct {
	CourseID  uuid.UUID `json:"course_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// Put extracts the User from JSON and writes it to DynamoDB
func Put(body string) (User, error) {
	var sess *session.Session
	var tableName string

	local, err := strconv.ParseBool(os.Getenv("AWS_SAM_LOCAL"))
	if err != nil {
		return User{}, err
	}
	// Create dynamo client object locally if running SAM CLI
	if local {
		sess = session.Must(session.NewSession(&aws.Config{
			Endpoint: aws.String("http://dynamodb:8000"),
		}))
		tableName = "users"
	} else {
		sess = session.Must(session.NewSession())
		tableName = os.Getenv("USER_TABLE_NAME")
	}
	svc := dynamodb.New(sess)

	// Marshall the requrest body
	var user User
	json.Unmarshal([]byte(body), &user)

	// Generate new UUID to store User in case user doesn't have one
	if user.ID == uuid.Nil {
		id, _ := uuid.NewV4()
		user.ID = id
	}

	// Marshall the Item into a Map DynamoDB can deal with
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		fmt.Println("Got error marshalling map:", err.Error())
		return user, err
	}

	// Create Item in table and return
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	return user, err
}

// Read gets the User from DynamoDB
func Read(id string) (User, error) {
	var sess *session.Session
	var tableName string

	local, err := strconv.ParseBool(os.Getenv("AWS_SAM_LOCAL"))
	if err != nil {
		return User{}, err
	}
	// Create dynamo client object locally if running SAM CLI
	if local {
		sess = session.Must(session.NewSession(&aws.Config{
			Endpoint: aws.String("http://dynamodb:8000"),
		}))
		tableName = "users"
	} else {
		sess = session.Must(session.NewSession())
		tableName = os.Getenv("USER_TABLE_NAME")
	}
	svc := dynamodb.New(sess)
	user := User{}

	// Perform the query
	fmt.Println("Trying to read from table: ", "users")
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				B: []byte(aws.StringValue(aws.String(id))),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}

	// Unmarshall the result in to an Item
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}

	return user, nil
}

// Delete erases the User from DynamoDB
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
		tableName = "users"
	} else {
		sess = session.Must(session.NewSession())
		tableName = os.Getenv("USER_TABLE_NAME")
	}
	svc := dynamodb.New(sess)

	// Perform the delete
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				B: []byte(aws.StringValue(aws.String(id))),
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

// List returns the Users from DynamoDB
func List() ([]User, error) {
	var sess *session.Session
	var tableName string

	local, err := strconv.ParseBool(os.Getenv("AWS_SAM_LOCAL"))
	if err != nil {
		return []User{}, err
	}
	// Create dynamo client object locally if running SAM CLI
	if local {
		sess = session.Must(session.NewSession(&aws.Config{
			Endpoint: aws.String("http://dynamodb:8000"),
		}))
		tableName = "users"
	} else {
		sess = session.Must(session.NewSession())
		tableName = os.Getenv("USER_TABLE_NAME")
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

	var users []User
	dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return users, nil
}
