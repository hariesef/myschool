package user_test

import (
	"context"
	"fmt"
	"myschool/internal/storage/mongodb/user"
	"myschool/pkg/helper"
	"myschool/pkg/model"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	localMongoClient *mongo.Client
	userRepoImpl     model.UserRepo
	ctrl             *gomock.Controller
)

func TestUsers(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {

		var err error
		opt := options.Client().ApplyURI("mongodb://localhost:27017")
		localMongoClient, err = mongo.Connect(context.Background(), opt)
		Expect(err).To((BeNil()))

		err = localMongoClient.Ping(context.Background(), readpref.Primary())
		Expect(err).To((BeNil()))

		db := localMongoClient.Database("unit_test_db")

		//initialize DB
		db.Drop(context.Background())

		//make email as unique index
		indexModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true)}
		_, err = db.Collection(user.UsersCollectionName).Indexes().CreateOne(context.TODO(), indexModel)
		Expect(err).To((BeNil()))

		userRepoImpl = user.NewRepo(db)
		ctrl = gomock.NewController(t)
	})

	AfterSuite(func() {
		err := localMongoClient.Disconnect(context.Background())
		Expect(err).To((BeNil()))
		ctrl.Finish()
	})

	RunSpecs(t, "User Suite")
}

var _ = Describe("Testing UserRepo Implementation with real local mongoDB", func() {

	BeforeEach(func() {

	})

	AfterEach(func() {

	})

	//serial mode due to I/O testing with DB
	Describe("Testing all user repo functions using local mongoDB", Serial, func() {
		Context("User Repo", func() {

			firstUserId := ""
			randomPassword := helper.RandomStringBytes(5)
			It("tests creating first user", func() {

				user, err := userRepoImpl.Create(context.TODO(),
					model.UserCreationParam{Email: "haries@banget.net", EncryptedPassword: randomPassword})
				firstUserId = user.GetID()
				fmt.Println("first user id: ", firstUserId)
				Expect(err).To((BeNil()))
				Expect(len(firstUserId)).Should(BeNumerically("==", 24))
				Expect(user.IsActive()).To(Equal(true))
			})

			It("tests creating 2nd user with same email, should fail.", func() {

				_, err := userRepoImpl.Create(context.TODO(),
					model.UserCreationParam{Email: "haries@banget.net", EncryptedPassword: "yyyyy"})
				Expect(err.Error()).To(ContainSubstring("duplicate key error collection"))
			})

			It("tests reading back 1st user.", func() {

				user, err := userRepoImpl.Read(context.TODO(), "haries@banget.net")
				Expect(err).To((BeNil()))
				Expect(user.GetEncryptedPassword()).To(Equal(randomPassword))
				Expect(user.GetID()).To(Equal(firstUserId))
			})

			It("tests deactivating the 1st user.", func() {

				user, err := userRepoImpl.Deactivate(context.TODO(), firstUserId)
				Expect(err).To((BeNil()))
				Expect(user.IsActive()).To(Equal(false))
			})

			It("tests reading back 1st user after deactivated.", func() {

				_, err := userRepoImpl.Read(context.TODO(), "haries@banget.net")
				Expect(err).To(Equal(mongo.ErrNoDocuments))
			})

			It("tests finding only active users", func() {

				//create another two users
				_, err := userRepoImpl.Create(context.TODO(),
					model.UserCreationParam{Email: "sulis@banget.net", EncryptedPassword: "yyyy"})
				Expect(err).To((BeNil()))
				_, err = userRepoImpl.Create(context.TODO(),
					model.UserCreationParam{Email: "ayu@banget.net", EncryptedPassword: "xxxx"})
				Expect(err).To((BeNil()))
				_, err = userRepoImpl.Create(context.TODO(),
					model.UserCreationParam{Email: "ryo@banget.net", EncryptedPassword: "zzzz"})
				Expect(err).To((BeNil()))

				users, err := userRepoImpl.FindActive(context.TODO())
				Expect(err).To((BeNil()))
				Expect(len(users)).To(Equal(3))
				Expect(users[0].IsActive()).To(Equal(true))
				Expect(users[0].GetEmail()).To(Equal("ayu@banget.net"))
				Expect(users[1].GetEmail()).To(Equal("ryo@banget.net"))
				Expect(users[2].GetEmail()).To(Equal("sulis@banget.net"))
			})

		})
	}) //end of describe

})
