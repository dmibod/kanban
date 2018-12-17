cd ..
go mod tidy
go mod vendor
docker build -f docker/Dockerfile -t dmibod/kanban-process .