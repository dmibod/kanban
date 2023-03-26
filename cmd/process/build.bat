cd ../..
rem go mod tidy
go mod vendor
docker build -f cmd/process/Dockerfile -t dmibod/kanban-process .