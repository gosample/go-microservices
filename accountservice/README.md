export GOOS=linux
export GOARCH=amd64
go build -o accountservice-linux-amd64


go run *.go

go build
./accountservice