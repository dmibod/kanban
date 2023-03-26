cd ../..
rem go mod tidy
go mod vendor
docker build -f cmd/notify/Dockerfile -t dmibod/kanban-notify .