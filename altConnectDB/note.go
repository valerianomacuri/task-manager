package altconnectdb

import (
	"context"
	"log"
	"time"

	"github.com/valerianomacuri/task-manager/common"
	"github.com/valerianomacuri/task-manager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type NoteRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
	ctx        context.Context
}

func NewNoteRepository() *NoteRepository {
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(common.AppConfig.MongoURI))
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	c := client.Database(common.AppConfig.Database).Collection("notes")
	return &NoteRepository{
		client:     client,
		collection: c,
		ctx:        ctx,
	}
}

// Close mongo.Client
func (r *NoteRepository) Close() {
	err := r.client.Disconnect(r.ctx)
	if err != nil {
		log.Println("failed to disconnect")
		return
	}
	log.Println("disconnected from MongoDB")
}

func (r *NoteRepository) Create(note *models.TaskNote) error {
	obj_id := primitive.NewObjectID()
	note.Id = obj_id
	note.CreatedOn = time.Now()
	_, err := r.collection.InsertOne(r.ctx, *note)
	return err
}

func (r *NoteRepository) GetByTask(id string) ([]models.TaskNote, error) {
	notes := make([]models.TaskNote, 0)
	taskid, _ := primitive.ObjectIDFromHex(id)
	cur, err := r.collection.Find(r.ctx, bson.M{
		"taskid": taskid,
	})
	if err != nil {
		return notes, err
	}
	defer cur.Close(r.ctx)
	for cur.Next(r.ctx) {
		var note models.TaskNote
		cur.Decode(&note)
		notes = append(notes, note)
	}
	return notes, nil
}

func (r *NoteRepository) GetAll() ([]models.TaskNote, error) {
	notes := make([]models.TaskNote, 0)
	cur, err := r.collection.Find(r.ctx, bson.M{})
	if err != nil {
		return notes, err
	}
	defer cur.Close(r.ctx)
	for cur.Next(r.ctx) {
		var note models.TaskNote
		cur.Decode(&note)
		notes = append(notes, note)
	}
	return notes, nil
}

func (r *NoteRepository) GetById(id string) (models.TaskNote, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := r.collection.FindOne(r.ctx, bson.M{
		"_id": objectId,
	})
	var note models.TaskNote
	err := cur.Decode(&note)
	if err != nil {
		return note, err
	}
	return note, nil
}

func (r *NoteRepository) Update(note *models.TaskNote) error {
	// partial update on MogoDB
	_, err := r.collection.UpdateOne(
		r.ctx,
		bson.M{
			"_id": note.Id,
		},
		bson.D{{"$set", bson.D{
			{"description", note.Description},
		}}},
	)
	return err
}
func (r *NoteRepository) Delete(id string) error {
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := r.collection.DeleteOne(r.ctx, bson.M{"_id": objectId})
	return err
}
