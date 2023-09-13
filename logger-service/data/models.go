package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}
type LogEntry struct {
	ID         string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string    `bson:"name" json:"name"`
	Data       string    `bson:"data" json:"data"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdateddAt time.Time `bson:"updated_at"  json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:       entry.Name,
		Data:       entry.Data,
		CreatedAt:  time.Now(),
		UpdateddAt: time.Now(),
	})
	if err != nil {
		log.Println("Inse rt error", err)
		return err
	}
	return nil
}
func (l *LogEntry) All() ([]*LogEntry, error) {
	// dispatch timeout ig get too long
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	options := options.Find()
	options.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, options)
	if err != nil {
		log.Println("error while retrieve collection", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Println("error decodding item into slice", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil

}

func (l *LogEntry) GetOneById(id string) (*LogEntry, error) {
	// dispatch timeout ig get too long
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil

}

func (l *LogEntry) DropCollection() error {
	// dispatch timeout ig get too long
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil

}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	// dispatch timeout ig get too long
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docId},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"update_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil

}
