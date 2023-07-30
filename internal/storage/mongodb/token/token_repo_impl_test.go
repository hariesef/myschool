package token_test

import (
	"context"
	"myschool/internal/storage/mongodb/token"
	"myschool/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestToken(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test create", func(mt *mtest.T) {

		tokenRepoImpl := token.NewRepo(mt.DB)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		model, err := tokenRepoImpl.Create(context.TODO(), model.TokenCreationParam{
			Token:  "123",
			UserID: "idharies",
			Email:  "haries@banget.net",
			Expiry: 10002000,
		})
		assert.Nil(t, err)
		assert.Equal(t, len(model.GetID()), 24)
	})

	mt.Run("test failed create", func(mt *mtest.T) {

		tokenRepoImpl := token.NewRepo(mt.DB)
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: -1},
			},
		)
		_, err := tokenRepoImpl.Create(context.TODO(), model.TokenCreationParam{
			Token:  "123",
			UserID: "idharies",
			Email:  "haries@banget.net",
			Expiry: 10002000,
		})
		assert.NotNil(t, err)
	})

	mt.Run("test find a token", func(mt *mtest.T) {

		tokenRepoImpl := token.NewRepo(mt.DB)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: "abcdefg"},
				{Key: "token", Value: "12345"},
				{Key: "userId", Value: "idharies"},
				{Key: "email", Value: "haries@banget.net"},
				{Key: "expiry", Value: 10002000},
			}))

		model, err := tokenRepoImpl.Find(context.TODO(), "12345")
		assert.Nil(t, err)
		assert.Equal(t, model.GetID(), "abcdefg")
		assert.Equal(t, model.GetToken(), "12345")
		assert.Equal(t, model.GetEmail(), "haries@banget.net")
		assert.Equal(t, model.GetUserID(), "idharies")
		assert.Equal(t, model.GetExpiry(), 10002000)
	})

	mt.Run("test find a token not found", func(mt *mtest.T) {

		tokenRepoImpl := token.NewRepo(mt.DB)
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch)) //notice the 0 and no bson data in last argument

		_, err := tokenRepoImpl.Find(context.TODO(), "12345")
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "mongo: no documents in result")
	})

	mt.Run("test delete a token", func(mt *mtest.T) {

		tokenRepoImpl := token.NewRepo(mt.DB)
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})
		err := tokenRepoImpl.Delete(context.TODO(), "12345")
		assert.Nil(t, err)
	})

	mt.Run("test delete a token but not found", func(mt *mtest.T) {

		tokenRepoImpl := token.NewRepo(mt.DB)
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 0}})
		err := tokenRepoImpl.Delete(context.TODO(), "12345")
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "the token is not found")
	})

	mt.Run("test delete a token but having other error", func(mt *mtest.T) {

		tokenRepoImpl := token.NewRepo(mt.DB)
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: -1},
			})
		err := tokenRepoImpl.Delete(context.TODO(), "12345")
		assert.NotNil(t, err)
	})

}

/*
For future reference

InsertMany
Similarl with InsertOne, only require a success response.

Find needs cursor response with one or multiple batches and an end of the cursor.
	first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{// our data})
	getMore := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{// our data})
	lastCursor := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
	mt.AddMockResponses(first, getMore, lastCursor)


FindOneAndUpdate
mt.AddMockResponses(bson.D{
   {"ok", 1},
   {"value", bson.D{// our data }},
})

Upsert
Same as FindOneAndUpdate.

FindOneAndDelete
Same as FindOneAndUpdate.



*/
