# 学习四种rpc调用

在当前目录使用protoc生成route.pb.go
`protoc -I route/ route.proto --go_out=plugins=grpc:route`

启动服务器
``` 
$ go run .\server\server.go
```
启动客户端
``` 
$ go run .\client\client.go
```
