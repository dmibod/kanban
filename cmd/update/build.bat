cd ../..
rem go mod tidy
go mod vendor
docker build -f cmd/update/Dockerfile -t dmibod/kanban-update .