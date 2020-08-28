package mondb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// DefaultMongoURL defines a local mongodb server address
const DefaultMongoURL = "mongodb://localhost:27017"

// Mongo is a implementation of Interface,
// using mongodb as database
type Mongo struct {
	// Name represents the database name
	Name string
	// Collection is the document collection used to
	// store the queries in the database database
	Collection string
	// QueryTimeout set context timeout for each query operation
	QueryTimeout int
	client       *mongo.Client
}

// New creates a Mongo instance.
// With query timeout up to 5 seconds
func New(dbName, collection string) *Mongo {
	return &Mongo{
		Name:         dbName,
		Collection:   collection,
		QueryTimeout: 5,
	}
}

// NewWithQueryTimeout creates a Mongo instance.
// queryTimeout is represented as seconds
func NewWithQueryTimeout(dbName, collection string, queryTimeout int) *Mongo {
	return &Mongo{
		Name:         dbName,
		Collection:   collection,
		QueryTimeout: queryTimeout,
	}
}

// Conn creates a connection to database
func (d *Mongo) Conn(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.QueryTimeout*2)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return err
	}
	d.client = client
	ctxPing, cancelPing := context.WithTimeout(context.Background(), time.Duration(d.QueryTimeout*2)*time.Second)
	defer cancelPing()
	err = client.Ping(ctxPing, readpref.Primary())
	return err
}

// Discn closes the connection to database
func (d *Mongo) Discn() {
	if d.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.QueryTimeout)*time.Second)
		d.client.Disconnect(ctx)
		cancel()
	}
}

// primitiveToBuiltin changes mongoDB types to go natives types
func primitiveToBuiltin(value map[string]interface{}) {
	for key, v := range value {
		switch castV := v.(type) {
		case primitive.A:
			value[key] = []interface{}(castV)
		default:
		}
	}
}

// FindOne using a filter find a document that matches the filter.
// If filter has _id param, then _id is converted to mongo Object ID,
// it must be in hexdecimal format.
func (d *Mongo) FindOne(filter map[string]interface{}) (map[string]interface{}, error) {
	if d.client == nil {
		return nil, ErrDBNotConnected
	}
	collection := d.client.Database(d.Name).Collection(d.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.QueryTimeout)*time.Second)
	defer cancel()
	var res map[string]interface{}
	if id, ok := filter["_id"]; ok {
		objID, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			return nil, ErrInvalidID
		}
		filter["_id"] = objID
	}
	if err := collection.FindOne(ctx, filter).Decode(&res); err == mongo.ErrNoDocuments {
		return nil, ErrDocumentNotFound
	} else if err != nil {
		return nil, err
	}
	res["_id"] = res["_id"].(primitive.ObjectID).Hex()
	primitiveToBuiltin(res)
	return res, nil
}

// FindMany using a filter finds all matching documents
func (d *Mongo) FindMany(filter map[string]interface{}) ([]map[string]interface{}, error) {
	if d.client == nil {
		return nil, ErrDBNotConnected
	}
	collection := d.client.Database(d.Name).Collection(d.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.QueryTimeout)*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	ctx, cancelCursorAll := context.WithTimeout(context.Background(), time.Duration(d.QueryTimeout)*time.Second)
	defer cancelCursorAll()
	var res []map[string]interface{}
	if err := cursor.All(ctx, &res); err == mongo.ErrNoDocuments || res == nil {
		return nil, ErrDocumentNotFound
	} else if err != nil {
		return nil, err
	}
	// Convert the ObjectIDs (_id) object to string
	for i := range res {
		res[i]["_id"] = res[i]["_id"].(primitive.ObjectID).Hex()
		primitiveToBuiltin(res[i])
	}
	return res, nil
}

// InsertOne inserts a new document to the database
// If obj has _id param, then _id is converted to mongo Object ID,
// it must be in hexdecimal format.
func (d *Mongo) InsertOne(obj map[string]interface{}) error {
	if d.client == nil {
		return ErrDBNotConnected
	}
	if obj == nil {
		return ErrEmptyObject
	}
	collection := d.client.Database(d.Name).Collection(d.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.QueryTimeout)*time.Second)
	defer cancel()
	if id, ok := obj["_id"]; ok {
		objID, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			return ErrInvalidID
		}
		obj["_id"] = objID
	}
	if _, err := collection.InsertOne(ctx, obj); err != nil {
		return err
	}
	return nil
}

// UpdateOne updates an existing document in the database
// If obj has _id param, then _id is converted to mongo Object ID,
// it must be in hexdecimal format.
func (d *Mongo) UpdateOne(filter map[string]interface{}, obj map[string]interface{}) (bool, error) {
	if d.client == nil {
		return false, ErrDBNotConnected
	}
	if obj == nil {
		return false, ErrEmptyObject
	}
	if id, ok := filter["_id"]; ok {
		objID, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			return false, ErrInvalidID
		}
		filter["_id"] = objID
	}
	if id, ok := obj["_id"]; ok {
		objID, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			return false, ErrInvalidID
		}
		obj["_id"] = objID
	}
	collection := d.client.Database(d.Name).Collection(d.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.QueryTimeout)*time.Second)
	defer cancel()
	result, err := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			bson.E{Key: "$set", Value: obj},
		},
	)
	return result.ModifiedCount > 0, err
}

// DeleteOne deletes a document
func (d *Mongo) DeleteOne(filter map[string]interface{}) (bool, error) {
	if d.client == nil {
		return false, ErrDBNotConnected
	}
	if id, ok := filter["_id"]; ok {
		objID, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			return false, ErrInvalidID
		}
		filter["_id"] = objID
	}
	collection := d.client.Database(d.Name).Collection(d.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.QueryTimeout)*time.Second)
	defer cancel()
	result, err := collection.DeleteOne(ctx, filter)
	return result.DeletedCount > 0, err
}
