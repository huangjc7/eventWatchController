rm -f eventwatch
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o eventwatch main.go
scp eventwatch root@192.168.217.10:/root
