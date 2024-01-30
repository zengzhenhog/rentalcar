@echo off
setlocal

call :createProto "rental"
call :createProto "auth"

exit /b

:createProto
set PROTO_PATH=.\%~1\api
set GO_OUT_PATH=.\%~1\api\gen\v1
mkdir %GO_OUT_PATH%

protoc -I=%PROTO_PATH% --go_out=paths=source_relative:%GO_OUT_PATH% %~1.proto
protoc -I=%PROTO_PATH% --go-grpc_opt require_unimplemented_servers=false --go-grpc_out=paths=source_relative:%GO_OUT_PATH% %~1.proto
protoc -I=%PROTO_PATH% --grpc-gateway_out=paths=source_relative,grpc_api_configuration=%PROTO_PATH%\%~1.yaml:%GO_OUT_PATH% %~1.proto
exit /b