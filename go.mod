module github.com/AsterNighT/software-engineering-backend

go 1.16

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/labstack/echo/v4 v4.2.1
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/swaggo/echo-swagger v1.1.0
	github.com/swaggo/swag v1.7.0
	golang.org/x/mod v0.4.0 // indirect
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1 // indirect
	golang.org/x/sys v0.0.0-20210403161142-5e06dd20ab57 // indirect
	golang.org/x/tools v0.1.0 // indirect
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.21.5
)

replace (
	github.com/AsterNighT/software-engineering-backend/docs => ./docs
	github.com/AsterNighT/software-engineering-backend/pkg/router => ./pkg/router
)
