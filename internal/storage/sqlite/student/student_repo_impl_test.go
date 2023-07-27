package student_test

import (
	"context"
	"myschool/pkg/model"
	"testing"

	sqLiteStudent "myschool/internal/storage/sqlite/student"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db              *gorm.DB
	studentRepoImpl model.StudentRepo
	ctrl            *gomock.Controller
)

func TestStudent(t *testing.T) {
	RegisterFailHandler(Fail)

	var err error

	BeforeSuite(func() {
		db, err = gorm.Open(sqliteDriver.Open("localtest.db"), &gorm.Config{})
		if err != nil {
			errString := "Cannot access localtest.db: " + err.Error()
			panic(errString)
		}
		studentRepoImpl = sqLiteStudent.NewRepo(db)
		db.Migrator().DropTable(&sqLiteStudent.Student{})
		db.Migrator().CreateTable(&sqLiteStudent.Student{})

		ctrl = gomock.NewController(t)

	})

	AfterSuite(func() {
		ctrl.Finish()
	})

	RunSpecs(t, "Student Suite")
}

var _ = Describe("Model Implementation", func() {

	var stu *sqLiteStudent.Student

	BeforeEach(func() {
		stu = &sqLiteStudent.Student{
			UID:       123,
			CreatedAt: 456,
			UpdatedAt: 789,
			Name:      "haries",
			Gender:    "M",
		}
	})

	Describe("Setter Getter", func() {
		Context("Getting values after setting up the object", func() {
			It("should match initial values", func() {

				Expect(stu.GetName()).To(Equal("haries"))
				Expect(stu.GetCreatedAt()).To(Equal(456))
			})
		})
	}) //end of describe

})

var _ = Describe("Student Repo Implementation", func() {

	BeforeEach(func() {

	})

	AfterEach(func() {

	})

	//serial mode due to I/O testing with DB
	Describe("Testing all student repo functions using actual SQLite DB", Serial, func() {
		Context("User Creation", func() {
			It("tests creating first user", func() {

				student, err := studentRepoImpl.Create(context.TODO(), model.StudentCreationParam{Name: "Steve Haries", Gender: "M"})
				Expect(err).To((BeNil()))
				Expect(student.GetUID()).To(Equal(uint(1)))
				Expect(student.GetName()).To(Equal("Steve Haries"))
				Expect(student.GetGender()).To(Equal("M"))
				Expect(student.GetCreatedAt()).Should(BeNumerically(">", 1600000000))
				Expect(student.GetUpdatedAt()).Should(BeNumerically(">", 1600000000))
			})

			It("tests creating 2nd user", func() {

				student, err := studentRepoImpl.Create(context.TODO(), model.StudentCreationParam{Name: "John Efrika", Gender: "F"})
				Expect(err).To((BeNil()))
				Expect(student.GetUID()).To(Equal(uint(2)))
				Expect(student.GetName()).To(Equal("John Efrika"))
			})

			It("tests creating 3rd user", func() {

				student, err := studentRepoImpl.Create(context.TODO(), model.StudentCreationParam{Name: "Steve Efrika", Gender: "F"})
				Expect(err).To((BeNil()))
				Expect(student.GetUID()).To(Equal(uint(3)))
				Expect(student.GetName()).To(Equal("Steve Efrika"))
			})

			It("tests reading the 2nd user", func() {

				student, err := studentRepoImpl.Read(context.TODO(), 2)
				Expect(err).To((BeNil()))
				Expect(student.GetName()).To(Equal("John Efrika"))
				Expect(student.GetGender()).To(Equal("F"))
				Expect(student.GetCreatedAt()).Should(BeNumerically(">", 1600000000))
				Expect(student.GetUpdatedAt()).Should(BeNumerically(">", 1600000000))
			})

			It("tests reading multiple users based on name contains teve", func() {

				students, err := studentRepoImpl.FindByName(context.TODO(), "teve")
				Expect(err).To((BeNil()))
				Expect(students[0].GetName()).To(Equal("Steve Haries"))
				Expect(students[1].GetName()).To(Equal("Steve Efrika"))
				Expect(len(students)).Should(BeNumerically("==", 2))
			})

			It("tests deleting the 2nd user", func() {

				student, err := studentRepoImpl.Delete(context.TODO(), 2)
				Expect(err).To((BeNil()))
				Expect(student.GetName()).To(Equal("John Efrika"))
				Expect(student.GetGender()).To(Equal("F"))
			})

			It("tests reading the 2nd user after deletion", func() {

				_, err := studentRepoImpl.Read(context.TODO(), 2)
				Expect(err.Error()).To(Equal("user not found"))
			})

			It("tests deleting again the 2nd user that had been deleted", func() {

				_, err := studentRepoImpl.Delete(context.TODO(), 2)
				Expect(err.Error()).To(Equal("user not found"))
			})

		})
	}) //end of describe

})
