package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var	client *mongo.Client
var	collections map[string]map[string]string

var databases map[string]*mongo.Database

const defaultDb = "default"

func Connect(host string, user string, password string, data map[string]interface{}) error {
	var err error
	ctx := context.TODO()
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/", user, password, host))
	client, err = mongo.Connect(ctx, clientOptions)

	collections = make(map[string]map[string]string)
	databases = make(map[string]*mongo.Database)
	for dbName, dbNameCollections := range data {
		collectionArray := make(map[string]string)
		databases[dbName] = client.Database(dbName)
		for action, collection := range dbNameCollections.(map[string]interface{}) {
			collectionArray[action] = collection.(string)
		}
		collections[dbName] = collectionArray
	}
	return err
}

func Disconnect() error {
	var err error

	if err = client.Disconnect(context.TODO()); err != nil {
		return err
	}

	return nil
}

func InsertValue(body []byte, pid string, createdTime int64, action string, level string) error {
	var message interface{}
	ctx := context.TODO()
	findDb := ""
	parseErr := json.Unmarshal(body, &message)
	if parseErr != nil {
		return parseErr
	}

	log := bson.D{
		{"pid", pid},
		{"createdTime", createdTime},
		{"level", level},
	}

	for key, data := range message.(map[string]interface{}) {
		log = append(log, primitive.E{Key: key, Value: data})
	}

	for db, collection := range collections {
		if _, ok := collection[action]; ok {
			findDb = db
			break
		}
	}

	if findDb == "" {findDb = defaultDb}

	if mapErr := checkMapKeys(findDb, action); mapErr != nil {
		return mapErr
	}

	currentCollection := databases[findDb].Collection(
		fmt.Sprintf("%s_%s",collections[findDb][action],time.Unix(0, createdTime).Format("200601")))

	_, err := currentCollection.InsertOne(ctx, log)

	if err != nil {
		return err
	}

	return nil
}

func checkMapKeys(findDb string, action string) error {
	if _, ok := databases[findDb]; !ok {
		return errors.New("log_consumer: undefined database")
	}

	if _, ok := collections[findDb]; !ok {
		return errors.New("log_consumer: undefined collection")
	}

	if _, ok := collections[findDb][action]; !ok {
		return errors.New("log_consumer: undefined action")
	}

	return nil
}
