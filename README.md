# term-bestetufan

Beste Tufan - Bitirme Projesi

Çalışma Sırası:
main.go -> app.go -> router.go -> server

1-) Update Swagger Definitions  : go install github.com/swaggo/swag/cmd/swag@latest
                                  swag init -g .\internal\api\router\router.go
2-) Create MYSQL Instance       : docker-compose up -d
3-) go run main.go
4-) teardown                    : docker-compose down