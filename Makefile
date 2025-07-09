.PHONY: build run tidy clean test

APP_NAME = go-service
BUILD_DIR = bin
CMD_PATH = ./cmd/app

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)

run: build
	./$(BUILD_DIR)/$(APP_NAME)

tidy:
	go mod tidy

clean:
	rm -rf $(BUILD_DIR)
