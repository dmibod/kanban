cd ..
go mod tidy
go mod vendor
docker build -f notify/Dockerfile -t dmibod/kanban-process .