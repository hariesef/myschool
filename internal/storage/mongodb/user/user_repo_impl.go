package user

import (
	"context"
	"myschool/pkg/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const UsersCollectionName string = "user"

// implmentation of UserModel
type User struct {
	ID                string `json:"id,omitempty" bson:"_id,omitempty"`
	Email             string `json:"email,omitempty" bson:"email,omitempty"`
	EncryptedPassword string `json:"encryptedPassword,omitempty" bson:"encryptedPassword,omitempty"`
	Active            bool   `json:"active,omitempty" bson:"active,omitempty"`
}

// Safe checker to know if this file already implements the model interface correctly or not
var _ (model.UserModel) = (*User)(nil)

func (s *User) GetID() string {
	return s.ID
}

func (s *User) GetEmail() string {
	return s.Email
}

func (s *User) GetEncryptedPassword() string {
	return s.EncryptedPassword
}

func (s *User) IsActive() bool {
	return s.Active
}

// unexportable
type repoPrivate struct {
	db mongo.Database
}

// Safe checker to know if this file already implements the interface correctly or not
var _ (model.UserRepo) = (*repoPrivate)(nil)

func NewRepo(db *mongo.Database) model.UserRepo {
	return &repoPrivate{db: *db}
}

func (repo *repoPrivate) Create(ctx context.Context, args model.UserCreationParam) (model.UserModel, error) {

	newUser := &User{
		Email:             args.Email,
		EncryptedPassword: args.EncryptedPassword,
		Active:            true,
	}
	res, err := repo.db.Collection(UsersCollectionName).InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}
	newUser.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return newUser, nil
}

func (repo *repoPrivate) Read(ctx context.Context, email string) (model.UserModel, error) {
	// objectId, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	return nil, err
	// }

	result := repo.db.Collection(UsersCollectionName).FindOne(context.Background(), bson.M{"email": email, "active": true})
	var user User
	result.Decode(&user)
	return &user, result.Err()
}

func (repo *repoPrivate) Deactivate(ctx context.Context, id string) (model.UserModel, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}
	update := bson.D{{Key: "$set", Value: bson.M{"active": false}}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := repo.db.Collection(UsersCollectionName).FindOneAndUpdate(context.Background(), filter, update, opts)
	var user User
	result.Decode(&user)
	return &user, result.Err()
}

func (repo *repoPrivate) FindActive(ctx context.Context) ([]model.UserModel, error) {
	var users []*User

	filter := User{Active: true}
	opts := options.Find().SetSort(bson.D{{Key: "email", Value: 1}})

	cur, err := repo.db.Collection(UsersCollectionName).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	models := make([]model.UserModel, len(users))
	for i, v := range users {
		models[i] = model.UserModel(v)
	}
	return models, err
}
