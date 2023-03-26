cd ../..
rem go mod tidy
go mod vendor
docker build -f cmd/monolith/Dockerfile -t dmibod/kanban-monolith .
