package mongodb

import (
	"context"

	"github.com/erikrios/my-events/lib/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	Events = "events"
)

type MongoDBLayer struct {
	db *mongo.Database
}

func NewMongoDBLayer(db *mongo.Database) persistence.DatabaseHandler {
	return &MongoDBLayer{
		db: db,
	}
}

func (m *MongoDBLayer) AddEvent(e persistence.Event) ([]byte, error) {
	e.Location.ID = primitive.NewObjectID()
	result, err := m.db.Collection(Events).InsertOne(context.TODO(), &e)
	if err != nil {
		return nil, persistence.ErrDatabase
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, persistence.ErrDatabase
	}

	return []byte(id.Hex()), nil
}

func (m *MongoDBLayer) FindEvent(id []byte) (persistence.Event, error) {
	hexID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return persistence.Event{}, persistence.ErrDatabase
	}

	filter := bson.D{primitive.E{Key: "_id", Value: hexID}}

	var event persistence.Event

	if err := m.db.Collection(Events).FindOne(context.TODO(), filter).Decode(&event); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return persistence.Event{}, persistence.ErrNotFound
		default:
			return persistence.Event{}, persistence.ErrDatabase
		}
	}

	return event, nil
}

func (m *MongoDBLayer) FindEventByName(name string) (persistence.Event, error) {
	filter := bson.D{primitive.E{Key: "name", Value: name}}

	var event persistence.Event

	if err := m.db.Collection(Events).FindOne(context.TODO(), filter).Decode(&event); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return persistence.Event{}, persistence.ErrNotFound
		default:
			return persistence.Event{}, persistence.ErrDatabase
		}
	}

	return event, nil
}

func (m *MongoDBLayer) FindAllAvailableEvents() ([]persistence.Event, error) {
	cursor, err := m.db.Collection(Events).Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, persistence.ErrDatabase
	}

	events := make([]persistence.Event, 0)
	if err := cursor.All(context.TODO(), &events); err != nil {
		return nil, persistence.ErrDatabase
	}

	return events, nil
}
