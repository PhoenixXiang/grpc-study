# 官方示例：helloworld

## try

在当前目录使用protoc生成helloworld.pb.go
`protoc -I helloworld/ helloworld.proto --go_out=plugins=grpc:helloworld`

启动服务器
``` 
$ go run .\greeter_server\main.go
```
启动客户端
``` 
$ go run .\greeter_client\main.go
```
看到`Greeting: Hello world`，表示成功！


## protoc命令使用

> 参照 https://studygolang.com/articles/9386

使用protoc命令编译.proto文件,不同语言支持需要指定输出参数，如：

`protoc --proto_path=IMPORT_PATH --cpp_out=DST_DIR --java_out=DST_DIR --python_out=DST_DIR --go_out=DST_DIR --ruby_out=DST_DIR --javanano_out=DST_DIR --objc_out=DST_DIR --csharp_out=DST_DIR path/to/file.proto`
这里详细介绍golang的编译姿势:

+ -I 参数：指定import路径，可以指定多个-I参数，编译时按顺序查找，不指定时默认查找当前目录

+ --go_out ：golang编译支持，支持以下参数

    + plugins=plugin1+plugin2 - 指定插件，目前只支持grpc，即：plugins=grpc

    + M 参数 - 指定导入的.proto文件路径编译后对应的golang包名(不指定本参数默认就是.proto文件中import语句的路径)

    + import_prefix=xxx - 为所有import路径添加前缀，主要用于编译子目录内的多个proto文件，这个参数按理说很有用，尤其适用替代一些情况时的M参数，但是实际使用时有个蛋疼的问题导致并不能达到我们预想的效果，自己尝试看看吧

    + import_path=foo/bar - 用于指定未声明package或go_package的文件的包名，最右面的斜线前的字符会被忽略

    + 末尾 :编译文件路径 .proto文件路径(支持通配符)

完整示例：

`protoc -I . --go_out=plugins=grpc,Mfoo/bar.proto=bar,import_prefix=foo/,import_path=foo/bar:. ./*.proto`