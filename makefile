push: 
	git add . && git commit -m "${m}" && git push origin staging

auth-protoc: 
	protoc --go_out=protoc_gen --go_opt=paths=source_relative --go-grpc_out=protoc_gen --go-grpc_opt=paths=source_relative ./protoc/auth/auth.proto
