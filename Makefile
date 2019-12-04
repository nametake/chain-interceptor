protoc:
	@protoc -I ./protobuf --go_out=plugins=grpc:./pb ./protobuf/ping.proto
