module github.com/Medveddo/rocket-science/payment

go 1.24.4

replace github.com/Medveddo/rocket-science/shared => ../shared

require (
	github.com/Medveddo/rocket-science/shared v0.0.0-00010101000000-000000000000
	github.com/brianvoe/gofakeit/v7 v7.3.0
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.73.0
)

require (
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
