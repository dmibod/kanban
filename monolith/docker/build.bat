cd ..
copy ..\main.go main.go
go mod tidy
go mod vendor
docker build -f docker/Dockerfile -t dmibod/kanban .
del main.go