cd ..
go mod tidy
go mod vendor
docker build -f query/Dockerfile -t dmibod/kanban-query .