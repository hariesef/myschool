go-gen:
	rm -rf internal/mocks/*
	go generate ./...

api-gen:
	protoc --twirp_out=. --go_out=. pkg/controller/rpc/student/student.proto
	protoc --twirp_out=. --go_out=. pkg/controller/rpc/account/account.proto


	