package repositories

import (
	"context"
	"fmt"
	"myschool/pkg/model"
	"os"
	"sync"
	"time"

	mongoConnection "myschool/internal/storage/mongodb"
	authTokenRepo "myschool/internal/storage/mongodb/token"
	userRepo "myschool/internal/storage/mongodb/user"
	sqLiteConnection "myschool/internal/storage/sqlite"
	studentRepo "myschool/internal/storage/sqlite/student"

	"github.com/dewanggasurya/logger/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	wg            *sync.WaitGroup
	StudentRepo   model.StudentRepo
	UserRepo      model.UserRepo
	AuthTokenRepo model.AuthTokenRepo
	mongoClient   *mongo.Client
}

func Setup(wg *sync.WaitGroup) *Repositories {

	//SQLite setup
	sqlDB, err := sqLiteConnection.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//MongoDB setup
	mdb, client, err := mongoConnection.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &Repositories{
		wg:            wg,
		StudentRepo:   studentRepo.NewRepo(sqlDB),
		UserRepo:      userRepo.NewRepo(mdb),
		AuthTokenRepo: authTokenRepo.NewRepo(mdb),
		mongoClient:   client,
	}
}

func (r *Repositories) Disconnect() {
	err := r.mongoClient.Disconnect(context.TODO())
	if err != nil {
		log.Debugf("Error while disconnecting mongoDB: ", err)
	}
	//test lagging disconnect:
	time.Sleep(3 * time.Second)
	r.wg.Done()
}
