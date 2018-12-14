cd ..\..
go mod tidy
go mod vendor
docker build -f monolith/docker/Dockerfile -t dmibod/kanban .
