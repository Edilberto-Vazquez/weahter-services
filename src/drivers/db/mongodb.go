package db

import (
	"context"
	"fmt"
	"log"

	"github.com/Edilberto-Vazquez/weahter-services/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBDriver struct {
	client *mongo.Client
}

func NewMongoDBConnection(dbConfig models.DBConfig) (*MongoDBDriver, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbConfig.URI))
	if err != nil {
		log.Fatal("Could not connect to mongoDB")
		return nil, err
	}
	return &MongoDBDriver{client: client}, nil
}

func parseProjectionFields(fields []string) primitive.M {
	projection := bson.M{}
	for _, field := range fields {
		projection[field] = 1
	}
	projection["_id"] = 0
	return projection
}

func (m *MongoDBDriver) GetRecords(query models.FindRecords, ctx context.Context) ([]map[string]interface{}, error) {
	fmt.Print(query)

	var results []map[string]interface{}

	coll := m.client.Database(query.DB).Collection(query.Collection)
	filter := bson.M{
		"datetime": bson.M{
			"$gte": query.DateStart,
			"$lte": query.DateEnd,
		},
	}

	var opts *options.FindOptions

	if len(query.Fields) > 0 {
		opts = options.Find().SetProjection(parseProjectionFields(query.Fields))
	}

	cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
