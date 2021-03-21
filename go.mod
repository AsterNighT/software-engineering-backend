module github.com/AsterNighT/software-engineering-backend

go 1.16

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/labstack/echo/v4 v4.2.1
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/swaggo/echo-swagger v1.1.0
	github.com/swaggo/swag v1.7.0
	golang.org/x/mod v0.4.0 // indirect
	golang.org/x/tools v0.1.0 // indirect
)

replace (
    github.com/AsterNighT/software-engineering-backend/docs => ./docs
    github.com/AsterNighT/software-engineering-backend/pkg/router => ./pkg/router
)