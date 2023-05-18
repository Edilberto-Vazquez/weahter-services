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

func (m MongoDBDriver) createLinealPipeline(query models.FindRecords) []primitive.M {

	var pipeline []primitive.M

	fieldsOperations := map[string]bson.M{
		"lightning": bson.M{
			"$sum": bson.M{
				"$cond": bson.A{bson.M{"$eq": []interface{}{"$lightning", true}}, 1, 0},
			},
		},
	}

	// create filter for dates range
	matchStage := bson.M{
		"datetime": bson.M{"$gte": query.DateStart, "$lt": query.DateEnd},
	}
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	// group by dates
	groupTimeStage := bson.M{
		"_id": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%dT%H:00:00.000Z", "date": "$datetime"}},
	}
	for _, field := range query.Fields {
		operation, ok := fieldsOperations[field]
		if ok {
			groupTimeStage[field] = operation
		} else {
			groupTimeStage[field] = bson.M{"$avg": fmt.Sprintf("$%s", field)}
		}
	}
	pipeline = append(pipeline, bson.M{"$group": groupTimeStage})

	// sort data by _id
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": 1}})

	// group by fields
	groupFieldStage := bson.M{"_id": nil}
	for _, field := range query.Fields {
		groupFieldStage[field] = bson.M{"$push": bson.M{"x": "$_id", "y": fmt.Sprintf("$%s", field)}}
	}
	pipeline = append(pipeline, bson.M{"$group": groupFieldStage})

	// se crea el stage final
	dataProject := make([]bson.M, 0)
	for _, field := range query.Fields {
		dataProject = append(
			dataProject,
			bson.M{
				"name": field,
				"data": fmt.Sprintf("$%s", field),
			},
		)
	}
	pipeline = append(pipeline, bson.M{"$project": bson.M{
		"_id":  0,
		"data": dataProject,
	}})

	return pipeline

}

func (m *MongoDBDriver) createRadialPipeline(query models.FindRecords) []primitive.M {
	var pipeline []primitive.M

	ranges := map[string]map[string]float64{
		"range1": {"min": 0.0, "max": 9.5},
		"range2": {"min": 9.5, "max": 19},
		"range3": {"min": 18, "max": 27.5},
		"range4": {"min": 27.5, "max": 38},
	}

	// create filter for dates range
	matchStage := bson.M{
		"datetime": bson.M{
			"$gte": query.DateStart,
			"$lt":  query.DateEnd,
		},
	}
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	// create ranges
	groupRangeStage := bson.M{"_id": nil}
	for k, v := range ranges {
		groupRangeStage[k] = bson.M{
			"$sum": bson.M{
				"$cond": []interface{}{
					bson.M{
						"$and": []interface{}{
							bson.M{"$gte": []interface{}{"$distance", v["min"]}},
							bson.M{"$lte": []interface{}{"$distance", v["max"]}},
						},
					},
					1,
					0,
				},
			},
		}
	}
	pipeline = append(pipeline, bson.M{"$group": groupRangeStage})

	return pipeline
}

func (m *MongoDBDriver) GetRecords(query models.FindRecords, ctx context.Context) ([]map[string]interface{}, error) {

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
	defer cursor.Close(ctx)

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *MongoDBDriver) GetLineChart(query models.FindRecords, ctx context.Context) (*models.LineChart, error) {

	var result *models.LineChart

	coll := m.client.Database(query.DB).Collection(query.Collection)

	var pipeline []primitive.M

	if len(query.Fields) > 0 {
		pipeline = m.createLinealPipeline(query)
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil

}

func (m *MongoDBDriver) GetRadialChart(query models.FindRecords, ctx context.Context) (*models.RadialRanges, error) {

	var result *models.RadialRanges

	coll := m.client.Database(query.DB).Collection(query.Collection)

	var pipeline []primitive.M

	if len(query.Fields) > 0 {
		pipeline = m.createRadialPipeline(query)
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil

}
