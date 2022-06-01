package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/valerianomacuri/task-manager/common"
	"github.com/valerianomacuri/task-manager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Handler for HTTP Post - "/users/register"
// Add a new User document
func Register(w http.ResponseWriter, r *http.Request) {
	var dataResource UserResource
	// Decode the incoming User json
	err := json.NewDecoder(r.Body).Decode(&dataResource)
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"Invalid User data",
			500,
		)
		return
	}
	user := &dataResource.Data
	ctx := NewContext()
	defer ctx.Close()
	collection := ctx.DBCollection("users")
	// save to database
	obj_id := primitive.NewObjectID()
	user.Id = obj_id
	hpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.HashPassword = hpass
	//clear the incoming text password
	user.Password = ""
	_, err = collection.InsertOne(context.TODO(), user)

	// Clean-up the hashpassword to eliminate it from response
	user.HashPassword = nil
	if j, err := json.Marshal(UserResource{Data: *user}); err != nil {
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

// Handler for HTTP Post - "/users/login"
// Authenticate with username and apssword
func Login(w http.ResponseWriter, r *http.Request) {
	var dataResource LoginResource
	var token string
	// Decode the incoming Login json
	err := json.NewDecoder(r.Body).Decode(&dataResource)
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"Invalid Login data",
			500,
		)
		return
	}
	loginModel := dataResource.Data
	loginUser := models.User{
		Email:    loginModel.Email,
		Password: loginModel.Password,
	}
	ctx := NewContext()
	defer ctx.Close()
	c := ctx.DBCollection("users")
	cur := c.FindOne(context.TODO(), bson.M{
		"email": loginUser.Email,
	})
	var user models.User
	err = cur.Decode(&user)

	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"Invalid login credentials",
			401,
		)
		return
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword(user.HashPassword, []byte(loginUser.Password))
	if err != nil {
		user = models.User{}
		common.DisplayAppError(
			w,
			err,
			"Invalid login credentials",
			401,
		)
		return
	}

	///if login is successful

	// Generate JWT token
	token, err = common.GenerateJWT(user.Email, "member")
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"Eror while generating the access token",
			500,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Clean-up the hashpassword to eliminate it from response JSON
	user.HashPassword = nil
	authUser := AuthUserModel{
		User:  user,
		Token: token,
	}
	j, err := json.Marshal(AuthUserResource{Data: authUser})
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
