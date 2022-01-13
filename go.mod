module github.com/dubbogo/dubbo-go-benchmark

go 1.15

require (
	dubbo.apache.org/dubbo-go/v3 v3.0.0-rc4-1
	github.com/apache/dubbo-go-hessian2 v1.10.2
	github.com/dubbogo/tools v0.0.0-00010101000000-000000000000
	github.com/dubbogo/triple v1.1.7
	github.com/golang/protobuf v1.5.2
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
)

replace (
	dubbo.apache.org/dubbo-go/v3 => ../dubbo-go
	github.com/dubbogo/tools => ../tools
)
