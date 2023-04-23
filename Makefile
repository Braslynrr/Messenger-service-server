BINARY_NAME=MessengerService.exe
MICRO_NAME=MicroBlob.exe
## build: builds all binaries


build_MicroBlob:
	@go build -o ./Micro/MicroBlob/${MICRO_NAME} ./Micro/MicroBlob
	@echo Micro blob built!


build_Micro: build_MicroBlob


build_front:
	@cd ServerFiles\messenger-ui && yarn build
	@cd ..\.. &
	@echo front end built

build_server:
	@go build -o ./Endpoint/main/${BINARY_NAME} ./Endpoint/main
	@echo back end built!

run_micro: build_Micro
	@cd Micro\MicroBlob && start /min /b "" ${MICRO_NAME}
	@cd ..\..\ &
	@echo Micro Blob started!

run_Back: build_front build_server
	@echo Starting...
	@cd Endpoint\main && start /min /b "" ${BINARY_NAME}
	@cd ..\..\ &
	@echo back end started!

run: run_Back run_micro

clean:
	@echo Cleaning...
	@DEL ${BINARY_NAME}
	@go clean
	@echo Cleaned!

start: run


stop_micro:
	@echo "Stopping..."
	@taskkill /IM ${MICRO_NAME} /F
	@echo Stopped back end


stop_back:
	@echo "Stopping..."
	@taskkill /IM ${BINARY_NAME} /F
	@echo Stopped back end


stop: stop_back stop_micro 


restart: stop start

test:
	@echo "Testing..."
	go test -v MessengerService/...