BINARY_NAME=MessengerService.exe
## build: builds all binaries

build_front:
	@cd ServerFiles\messenger-ui && yarn build
	@cd ..\.. &
	@echo front end build

build_server:
	@go build -o ./Endpoint/main/${BINARY_NAME} ./Endpoint/main
	@echo back end built!

run: build_front build_server
	@echo Starting...
	@cd Endpoint\main && start /min /b "" ${BINARY_NAME}
	@cd ..\..\ &
	@echo back end started!

clean:
	@echo Cleaning...
	@DEL ${BINARY_NAME}
	@go clean
	@echo Cleaned!

start: run

stop:
	@echo "Stopping..."
	@taskkill /IM ${BINARY_NAME} /F
	@echo Stopped back end

restart: stop start

test:
	@echo "Testing..."
	go test -v MessengerService/...