package data

import (
	"context"
	"time"

	"github.com/valerianomacuri/task-manager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(c *mongo.Collection) *TaskRepository {
	return &TaskRepository{collection: c}
}

func (r *TaskRepository) Create(task *models.Task) error {
	obj_id := primitive.NewObjectID()
	task.Id = obj_id
	task.CreatedOn = time.Now()
	task.Status = "Created"
	_, err := r.collection.InsertOne(context.TODO(), *task)
	return err
}

func (r *TaskRepository) GetAll() ([]models.Task, error) {
	ctx := context.TODO()
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return []models.Task{}, err
	}
	defer cur.Close(ctx)
	tasks := make([]models.Task, 0)
	for cur.Next(context.TODO()) {
		var task models.Task
		cur.Decode(&task)
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) GetById(id string) (models.Task, error) {
	ctx := context.TODO()
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := r.collection.FindOne(ctx, bson.M{
		"_id": objectId,
	})
	var task models.Task
	err := cur.Decode(&task)
	if err != nil {
		return task, err
	}
	return task, nil
}
func (r *TaskRepository) GetByUser(user string) ([]models.Task, error) {
	ctx := context.TODO()
	cur, err := r.collection.Find(ctx, bson.M{"createdby": user})
	if err != nil {
		return []models.Task{}, err
	}
	defer cur.Close(ctx)
	tasks := make([]models.Task, 0)
	for cur.Next(ctx) {
		var task models.Task
		cur.Decode(&task)
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) Update(task *models.Task) error {
	ctx := context.TODO()
	// partial update on MogoDB
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"_id": task.Id,
	}, bson.D{{"$set", bson.D{
		{"name", task.Name},
		{"description", task.Description},
		{"due", task.Due},
		{"status", task.Status},
		{"tags", task.Tags},
	}}})

	// alternative update
	// _, err := r.collection.UpdateOne(ctx, bson.M{
	// 	"_id": task.Id,
	// }, bson.D{{"$set", bson.D{
	// 	{Key: "name", Value: task.Name},
	// 	{Key: "description", Value: task.Description},
	// 	{Key: "due", Value: task.Due},
	// 	{Key: "tags", Value: task.Tags},
	// 	{Key: "status", Value: task.Status},
	// }}})
	return err
}

func (r *TaskRepository) Delete(id string) error {
	ctx := context.TODO()
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectId})
	return err
}
