cd ..
go mod tidy
go mod vendor
docker build -f command/Dockerfile -t dmibod/kanban-command .