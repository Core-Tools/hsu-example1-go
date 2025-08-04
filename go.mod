module github.com/core-tools/hsu-example1-go

go 1.22.3

replace github.com/core-tools/hsu-core => github.com/core-tools/hsu-core/go v0.0.0-20250804105037-c5ad0603dad5

replace github.com/core-tools/hsu-echo => .

require (
	github.com/core-tools/hsu-core v0.0.0-00010101000000-000000000000
	github.com/core-tools/hsu-echo v0.0.0-00010101000000-000000000000
	github.com/jessevdk/go-flags v1.6.1
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
