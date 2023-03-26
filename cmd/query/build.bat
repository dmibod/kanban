cd ../..
rem go mod tidy
go mod vendor
docker build -f cmd/query/Dockerfile -t dmibod/kanban-query .