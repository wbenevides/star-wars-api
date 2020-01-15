package db

import (
	"context"
	"log"

	"github.com/wallacebenevides/star-wars-api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseHelper interface {
	Collection(name string) CollectionHelper
	Client() ClientHelper
}

type CollectionHelper interface {
	FindOne(context.Context, interface{}) SingleResultHelper
	InsertOne(context.Context, interface{}) (interface{}, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	Find(ctx context.Context, filter interface{}) (CursorHelper, error)
}

type SingleResultHelper interface {
	Decode(v interface{}) error
}

type ClientHelper interface {
	Database(string) DatabaseHelper
}

type CursorHelper interface {
	All(ctx context.Context, v interface{}) error
	Close(context.Context) error
	Decode(interface{}) error
	Next(context.Context) bool
}

type mongoClient struct {
	cl *mongo.Client
}

type mongoDatabase struct {
	db *mongo.Database
}

type mongoCollection struct {
	coll *mongo.Collection
}

type mongoCursor struct {
	cs *mongo.Cursor
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

func NewClient(cnf *config.Database) (ClientHelper, error) {
	log.Println("initializing a session with db ", cnf.DatabaseName, cnf.Uri)
	clientOptions := options.Client().ApplyURI(cnf.Uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}
	// Check the connection
	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}
	return &mongoClient{cl: client}, err
}

func NewDatabase(cnf *config.Database, client ClientHelper) DatabaseHelper {
	return client.Database(cnf.DatabaseName)
}

func (mc *mongoClient) Database(dbName string) DatabaseHelper {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (md *mongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() ClientHelper {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}) (CursorHelper, error) {
	cursor, err := mc.coll.Find(ctx, filter)
	return cursor, err
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResultHelper {
	singleResult := mc.coll.FindOne(ctx, filter)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.coll.InsertOne(ctx, document)
	return id.InsertedID, err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	deleteResult, err := mc.coll.DeleteOne(ctx, filter)
	return deleteResult, err
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

func (cs *mongoCursor) All(ctx context.Context, v interface{}) error {
	return cs.All(ctx, v)
}

func (cs *mongoCursor) Close(ctx context.Context) error {
	return cs.Close(ctx)
}
func (cs *mongoCursor) Decode(v interface{}) error {
	return cs.Decode(v)
}

func (cs *mongoCursor) Next(ctx context.Context) bool {
	return cs.Next(ctx)
}
