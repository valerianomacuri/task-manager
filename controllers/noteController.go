package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/valerianomacuri/task-manager/common"
	"github.com/valerianomacuri/task-manager/data"
	"github.com/valerianomacuri/task-manager/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler for HTTP Post - "/notes"
// Insert a new Note document for a TaskId
func CreateNote(w http.ResponseWriter, r *http.Request) {
	var dataResource NoteResource
	// Decode the incoming Note json
	err := json.NewDecoder(r.Body).Decode(&dataResource)
	if err != nil {
		common.DisplayAppError(w, err, "Invalid Note data", 500)
		return
	}
	noteModel := dataResource.Data
	var objectId primitive.ObjectID
	objectId, err = primitive.ObjectIDFromHex(noteModel.TaskId)
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	}
	note := &models.TaskNote{
		TaskId:      objectId,
		Description: noteModel.Description,
	}
	context := NewContext()
	defer context.Close()
	c := context.DBCollection("notes")
	//Insert a note document
	repo := data.NewNoteRepository(c)
	repo.Create(note)

	if j, err := json.Marshal(note); err != nil {
		common.DisplayAppError(
			w,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(j)
	}

}

// Handler for HTTP Get - "/notes/tasks/{id}
// Returns all Notes documents under a TaskId
func GetNotesByTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	context := NewContext()
	defer context.Close()
	c := context.DBCollection("notes")
	repo := data.NewNoteRepository(c)
	notes, err := repo.GetByTask(id)
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	}
	if j, err := json.Marshal(NotesResource{Data: notes}); err != nil {
		common.DisplayAppError(
			w,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}

}

// Handler for HTTP Get - "/notes"
// Returns all Note documents
func GetNotes(w http.ResponseWriter, r *http.Request) {
	context := NewContext()
	defer context.Close()
	c := context.DBCollection("notes")
	repo := data.NewNoteRepository(c)
	notes, err := repo.GetAll()
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	}
	if j, err := json.Marshal(NotesResource{Data: notes}); err != nil {
		common.DisplayAppError(
			w,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}

}

// Handler for HTTP Get - "/notes/{id}"
// Returns a single Note document by id
func GetNoteById(w http.ResponseWriter, r *http.Request) {
	// Get id from the incoming url
	vars := mux.Vars(r)
	id := vars["id"]
	context := NewContext()
	defer context.Close()
	c := context.DBCollection("notes")
	repo := data.NewNoteRepository(c)
	note, err := repo.GetById(id)
	if err != nil {
		if err == mongo.ErrNilDocument {
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
	}
	if j, err := json.Marshal(note); err != nil {
		common.DisplayAppError(
			w,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}

// Handler for HTTP Put - "/notes/{id}"
// Update an existing Note document
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	// Get id from the incoming url
	vars := mux.Vars(r)
	objectId, _ := primitive.ObjectIDFromHex(vars["id"])
	var dataResource NoteResource
	// Decode the incoming Note json
	err := json.NewDecoder(r.Body).Decode(&dataResource)
	if err != nil {
		common.DisplayAppError(w, err, "Invalid Note data", 500)
		return
	}
	noteModel := dataResource.Data
	note := &models.TaskNote{
		Id:          objectId,
		Description: noteModel.Description,
	}
	context := NewContext()
	defer context.Close()
	c := context.DBCollection("notes")
	repo := data.NewNoteRepository(c)
	//Update note document
	if err := repo.Update(note); err != nil {
		common.DisplayAppError(
			w,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

// Handler for HTTP Delete - "/notes/{id}"
// Delete an existing Note document
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	// Get id from the incoming url
	vars := mux.Vars(r)
	id := vars["id"]
	context := NewContext()
	defer context.Close()
	c := context.DBCollection("notes")
	repo := data.NewNoteRepository(c)
	//Delete a note document
	err := repo.Delete(id)
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
