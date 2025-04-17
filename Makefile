PROJECT_NAME = "go-simple-bank-v2"

BINARY_NAME = ${PROJECT_NAME}
BINARIES = "./bin"
MAIN_DIR = "cmd/${BINARY_NAME}"

GIT_LOCAL_NAME = "rodziievskyi-maksym"
GIT_LOCAL_EMAIL = "rodziyevskydev@gmail.com"

GITHUB = "github.com/${GIT_LOCAL_NAME}/${PROJECT_NAME}"
POSTGRES_URL = "postgresql://postgres:postgres@localhost:5434/daily-dose?sslmode=disable"

init:
	@echo "::> Creating a module root..."
	@go mod init ${GITHUB}
	@mkdir "cmd" && mkdir "cmd/"${BINARY_NAME}
	@touch ${MAIN_DIR}/main.go
	@echo "package main\n\nimport \"fmt\"\n\nfunc main(){\n\tfmt.Println(\"${BINARY_NAME}\")\n}" > ${MAIN_DIR}/main.go
	@touch VERSION && echo 0.0.1 > VERSION
	@git add ${MAIN_DIR}/main.go go.mod VERSION
	@echo "::> Finished!"

build:
	@echo "::> Building..."
	@go build -o ${BINARIES}/${BINARY_NAME} ${MAIN_DIR}
	@echo "::> Finished!"

run:
	@go build -o ${BINARIES}/${BINARY_NAME} ${MAIN_DIR}
	@${BINARIES}/${BINARY_NAME}

clean:
	@echo "::> Cleaning..."
	@go clean
	@rm -rf ${BINARIES}
	@go mod tidy
	@echo "::> Finished"

local-git:
	@git config --local user.name ${GIT_LOCAL_NAME}
	@git config --local user.email ${GIT_LOCAL_EMAIL}
	@git config --local --list

git-init:
	@echo "::> Git initialization begin..."
	@git init
	@git config --local user.name ${GIT_LOCAL_NAME}
	@git config --local user.email ${GIT_LOCAL_EMAIL}
	@touch .gitignore
	@echo ".idea" > .gitignore
	@echo "bin" > .gitignore
	@touch README.md
	@git add README.md
	@git commit -m "first commit"
	@git branch -M main
	@git remote add origin https://${GITHUB}
	@git push -u origin main
	@echo "::> Finished"

## Database operations
postgres:
	docker run --name go-simple-bank-db -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15.1-alpine

create-db:
	docker exec -it go-simple-bank-db createdb --username=root --owner=root go-simple-bank

drop-db:
	docker exec -it go-simple-bank-db dropdb go-simple-bank

migrate-up:
	migrate -path migrations -database "postgresql://root:secret@localhost:5433/go-simple-bank?sslmode=disable" -verbose up
migrate-down:
	migrate -path migrations -database "postgresql://root:secret@localhost:5433/go-simple-bank?sslmode=disable" -verbose down

migrate-up-last:
	migrate -path migrations -database "postgresql://root:secret@localhost:5433/go-simple-bank?sslmode=disable" -verbose up 1
migrate-down-last:
	migrate -path migrations -database "postgresql://root:secret@localhost:5433/go-simple-bank?sslmode=disable" -verbose down 1

# Create migration file
cm:
	@migrate create -ext sql -dir migrations -seq $(a)

sqlc:
	@cd "internal/infrastructure/database/"; sqlc generate

test:
	go test -v -cover ./...

mock:
	@mockgen -destination internal/repo/mock/store.go github.com/max-rodziyevsky/go-simple-bank/internal/repo Store

.PNONY: init build run clean local-git git-init postgres create-db drop-db migrate-up migrate-down sqlc test mock migrate-down-last migrate-up-last