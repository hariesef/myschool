go-gen:
	rm -rf internal/mocks/*
	go generate ./...

api-gen:
	protoc --twirp_out=. --go_out=. pkg/controller/rpc/student/student.proto
	protoc --twirp_out=. --go_out=. pkg/controller/rpc/account/account.proto

go-coverage: 
	@if [ -f coverage.txt ]; then rm coverage.txt; fi;
	@echo ">> running unit test and calculating coverage"
	@ginkgo -r -cover -output-dir=. -coverprofile=coverage.txt -covermode=count -coverpkg= \
		internal/storage/sqlite/student \
		internal/controller/rpc/student \
		internal/storage/mongodb/user \
		internal/services/account
	@go tool cover -func=coverage.txt

	
