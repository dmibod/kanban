cd ..
go mod tidy
go mod vendor
docker build -f update/Dockerfile -t dmibod/kanban-update .