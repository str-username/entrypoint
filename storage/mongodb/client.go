package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client

// Client structure
type Client struct {
	Mongodb *mongo.Client
}

// NewClient : mongodb client constructor
func NewClient(uri string) (*Client, error) {
	newClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Client{Mongodb: newClient}, nil
}

// Disconnect : implement disconnect method
func (client *Client) Disconnect(ctx context.Context) error {
	return client.Mongodb.Disconnect(ctx)
}

// Requests

// MongoClient interface
type MongoClient interface {
	Disconnect(ctx context.Context) error
	InsertOne(ctx context.Context, db, coll string, doc interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, db, coll string, filter, update interface{}) (*mongo.UpdateResult, error)
	ReplaceOne(ctx context.Context, db, coll string, filter, replacement interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, db, coll string, filter interface{}) (*mongo.DeleteResult, error)
	FindOne(ctx context.Context, db, coll string, filter interface{}) (*bson.M, error)
}

// InsertOne : implement insert one document method
func (client *Client) InsertOne(ctx context.Context, db, coll string, doc interface{}) (*mongo.InsertOneResult, error) {
	collection := client.Mongodb.Database(db).Collection(coll)
	return collection.InsertOne(ctx, doc)
}

// UpdateOne : implement update one document method
func (client *Client) UpdateOne(ctx context.Context, db, coll string, filter, update interface{}) (*mongo.UpdateResult, error) {
	collection := client.Mongodb.Database(db).Collection(coll)
	return collection.UpdateOne(ctx, filter, update)
}

// ReplaceOne : implement replace one document method
func (client *Client) ReplaceOne(ctx context.Context, db, coll string, filter, replacement interface{}) (*mongo.UpdateResult, error) {
	collection := client.Mongodb.Database(db).Collection(coll)
	replaceOptions := options.Replace().SetUpsert(true)
	return collection.ReplaceOne(ctx, filter, replacement, replaceOptions)
}

// DeleteOne : implement delete one document method
func (client *Client) DeleteOne(ctx context.Context, db, coll string, filter interface{}) (*mongo.DeleteResult, error) {
	collection := client.Mongodb.Database(db).Collection(coll)
	return collection.DeleteOne(ctx, filter)
}

// FindOne : implement find one document method
func (client *Client) FindOne(ctx context.Context, db, coll string, filter interface{}) (*bson.M, error) {
	var result bson.M
	findOneOptions := options.FindOne().SetProjection(bson.M{"_id": 0})
	collection := client.Mongodb.Database(db).Collection(coll)
	if err := collection.FindOne(ctx, filter, findOneOptions).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
