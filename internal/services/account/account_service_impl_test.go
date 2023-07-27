package account_test

import (
	"context"
	"myschool/internal/repositories"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	authTokenRepo "myschool/internal/storage/mongodb/token"
	userRepo "myschool/internal/storage/mongodb/user"

	accountSvcImpl "myschool/internal/services/account"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	localMongoClient *mongo.Client
	ctrl             *gomock.Controller
	accountService   *accountSvcImpl.AccountService
)

func TestAccount(t *testing.T) {
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

		//make email unique index
		indexModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true)}
		_, err = db.Collection(userRepo.UsersCollectionName).Indexes().CreateOne(context.TODO(), indexModel)
		Expect(err).To((BeNil()))

		repo := &repositories.Repositories{
			UserRepo:      userRepo.NewRepo(db),
			AuthTokenRepo: authTokenRepo.NewRepo(db),
		}
		accountService = &accountSvcImpl.AccountService{Repo: repo}
		ctrl = gomock.NewController(t)
	})

	AfterSuite(func() {
		err := localMongoClient.Disconnect(context.Background())
		Expect(err).To((BeNil()))
		ctrl.Finish()
	})

	RunSpecs(t, "Account Suite")
}

var _ = Describe("Testing Account Service Implementation with real local mongoDB", func() {

	BeforeEach(func() {

	})

	AfterEach(func() {

	})

	//serial mode due to I/O testing with DB
	Describe("Testing account services using local mongoDB", Serial, func() {
		Context("User Repo", func() {

			It("tests creating first user", func() {

				err := accountService.Create(context.TODO(), "haries@banget.net", "plainpassword")
				Expect(err).To((BeNil()))
			})

			token := ""
			It("tests logging in first user", func() {

				tokenInfo, err := accountService.Login(context.TODO(), "haries@banget.net", "plainpassword")
				Expect(err).To((BeNil()))
				Expect(tokenInfo.Expiry).Should(BeNumerically(">", 1600000000))
				token = tokenInfo.Token
				Expect(len(token)).Should(BeNumerically("==", 32))
			})

			It("tests logging OUT first user", func() {

				err := accountService.Logout(context.TODO(), token)
				Expect(err).To((BeNil()))
			})

			It("tests logging OUT first user Again, error", func() {

				err := accountService.Logout(context.TODO(), token)
				Expect(err.Error()).To(Equal("the token is not found"))
			})

		})
	}) //end of describe

})
