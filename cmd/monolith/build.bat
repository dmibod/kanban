cd ..
go mod tidy
go mod vendor
docker build -f monolith/Dockerfile -t dmibod/kanban .
