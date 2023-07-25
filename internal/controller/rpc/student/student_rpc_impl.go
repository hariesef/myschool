package student

import (
	"context"

	pb "myschool/pkg/controller/rpc/student"

	"myschool/internal/repositories"
	"myschool/pkg/model"

	"github.com/twitchtv/twirp"
)

type StudentRPCServer struct {
	Repo *repositories.Repositories
}

func (s *StudentRPCServer) Create(ctx context.Context, sp *pb.StudentParam) (*pb.StudentModel, error) {

	studentParam := model.StudentCreationParam{
		Name:   sp.Name,
		Gender: sp.Gender,
	}
	studentModel, err := s.Repo.StudentRepo.Create(ctx, studentParam)
	if err != nil {
		return nil, twirp.InternalError(err.Error())
	}

	return &pb.StudentModel{
		Id:        int32(studentModel.GetUID()),
		CreatedAt: int32(studentModel.GetCreatedAt()),
		UpdatedAt: int32(studentModel.GetUpdatedAt()),
		Name:      studentModel.GetName(),
		Gender:    studentModel.GetGender(),
	}, nil
}

func (s *StudentRPCServer) Read(ctx context.Context, sid *pb.StudentID) (*pb.StudentModel, error) {

	studentModel, err := s.Repo.StudentRepo.Read(ctx, uint(sid.Id))
	if err != nil {
		return nil, twirp.InternalError(err.Error())
	}

	return &pb.StudentModel{
		Id:        int32(studentModel.GetUID()),
		CreatedAt: int32(studentModel.GetCreatedAt()),
		UpdatedAt: int32(studentModel.GetUpdatedAt()),
		Name:      studentModel.GetName(),
		Gender:    studentModel.GetGender(),
	}, nil
}

func (s *StudentRPCServer) Delete(ctx context.Context, sid *pb.StudentID) (*pb.StudentModel, error) {

	studentModel, err := s.Repo.StudentRepo.Delete(ctx, uint(sid.Id))
	if err != nil {
		return nil, twirp.InternalError(err.Error())
	}

	return &pb.StudentModel{
		Id:        int32(studentModel.GetUID()),
		CreatedAt: int32(studentModel.GetCreatedAt()),
		UpdatedAt: int32(studentModel.GetUpdatedAt()),
		Name:      studentModel.GetName(),
		Gender:    studentModel.GetGender(),
	}, nil
}

func (s *StudentRPCServer) FindByName(ctx context.Context, pbName *pb.StudentName) (*pb.StudentModels, error) {

	foundStudents, err := s.Repo.StudentRepo.FindByName(ctx, pbName.Name)
	if err != nil {
		return nil, twirp.InternalError(err.Error())
	}

	models := make([]*pb.StudentModel, len(foundStudents))
	for i, v := range foundStudents {
		models[i] = &pb.StudentModel{
			Id:        int32(v.GetUID()),
			CreatedAt: int32(v.GetCreatedAt()),
			UpdatedAt: int32(v.GetUpdatedAt()),
			Name:      v.GetName(),
			Gender:    v.GetGender(),
		}
	}
	return &pb.StudentModels{
		Students: models,
	}, nil
}
