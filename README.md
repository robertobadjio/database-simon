# :zap: Simon Database

<div align="center">

Blazing fast database in Go.

![Simon logo](assets/logo.png.webp)

</div>

### Start server master
```
go run cmd/server/server.go -config=./config.yml
```

### Start replica server
```
go run cmd/server/server.go -config=./config_replica.yml
```

### Connect to DB
```
go run cmd/cli/cli.go
```