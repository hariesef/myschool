package student_test

import (
	"context"
	"regexp"
	"testing"

	sqLiteStudent "myschool/internal/storage/sqlite/student"

	"myschool/internal/repositories"

	internalRPC "myschool/internal/controller/rpc/student"
	pkgRPC "myschool/pkg/controller/rpc/student"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ctrl *gomock.Controller
)

func TestStudent(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		ctrl = gomock.NewController(t)
	})

	AfterSuite(func() {
		ctrl.Finish()
	})

	RunSpecs(t, "Student Suite")
}

var _ = Describe("Student Repo Implementation", func() {

	var mock sqlmock.Sqlmock
	var repo *repositories.Repositories
	var server *internalRPC.StudentRPCServer

	BeforeEach(func() {
		mockDb, _mock, _ := sqlmock.New()
		mock = _mock
		dialector := postgres.New(postgres.Config{
			Conn:       mockDb,
			DriverName: "postgres",
		})
		db, err := gorm.Open(dialector, &gorm.Config{})
		Expect(err).ShouldNot(HaveOccurred())

		studentRepoImpl := sqLiteStudent.NewRepo(db)
		repo = &repositories.Repositories{
			StudentRepo: studentRepoImpl,
		}

		server = &internalRPC.StudentRPCServer{Repo: repo}

	})

	AfterEach(func() {

	})

	Describe("Testing RPC functions using sqlMock", func() {
		Context("User Creation", func() {
			It("tests creating a user", func() {

				mock.ExpectBegin()

				query := regexp.QuoteMeta("INSERT INTO \"students\" (\"created_at\",\"updated_at\",\"deleted_at\",\"name\",\"gender\")")
				rows := sqlmock.
					NewRows([]string{"uid"}).
					AddRow(777)

				mock.ExpectQuery(query).WillReturnRows(rows)
				mock.ExpectCommit()

				student, err := server.Create(context.TODO(), &pkgRPC.StudentParam{Name: "Mama Ishar", Gender: "F"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(student.GetId()).To(Equal(int32(777)))
				Expect(student.GetName()).To(Equal("Mama Ishar"))
				Expect(student.GetGender()).To(Equal("F"))
				Expect(student.GetCreatedAt()).Should(BeNumerically(">", 1600000000))
			})

		})
	}) //end of describe

})
