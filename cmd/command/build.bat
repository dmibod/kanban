cd ../..
rem go mod tidy
go mod vendor
docker build -f cmd/command/Dockerfile -t dmibod/kanban-command .