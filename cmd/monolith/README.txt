- create .env file here with the following contents:

MONGO_USER=******
MONGO_PASSWORD=******
NATS_USER=******
NATS_PASSWORD=******

- replace ****** with real values

====================================================

work with mongo:

- open terminal
- mongosh mongodb://user:password@localhost:27017/?authSource=admin or mongodb+srv://user:password@localhost:27017/?authSource=admin
- use kanban
- db.help()
- db.getName()
- db.getCollectionNames()
- db.createCollection("boards")
- ...

work with nats:

set BUS_URL nats://user:password@localhost:4222
export BUS_URL="nats://user:password@localhost:4222"
export BUS_CLIENT=""