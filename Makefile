PROJECT_PATH=$(shell pwd)
MODULE_NAME=live-chat-server

BUILD_NUM_FILE=build_num.txt
BUILD_NUM=$$(cat ./build_num.txt)
APP_VERSION=0.0
TARGET_VERSION=$(APP_VERSION).$(BUILD_NUM)
TARGET_DIR=bin
OUTPUT=$(PROJECT_PATH)/$(TARGET_DIR)/$(MODULE_NAME)

# app 구성
REDIS_SINGLE_MAIN=/cmd/redis_single/main.go
MYSQL_KAFKA_CONTROLLER=/cmd/mysql_kafka/controller/main.go
MYSQL_KAFKA_WORKER=/cmd/mysql_kafka/worker/main.go

LDFLAGS=-X main.BUILD_TIME=`date -u '+%Y-%m-%d_%H:%M:%S'`
LDFLAGS+=-X main.APP_VERSION=$(TARGET_VERSION)
LDFLAGS+=-X main.GIT_HASH=`git rev-parse HEAD`
LDFLAGS+=-s -w

redis-single: config test redis_single-build
mk_controller: config test mk_controller-build
mk_worker: config test mk_worker-build

config:
	@if [ ! -d $(TARGET_DIR) ]; then mkdir $(TARGET_DIR); fi

redis_single-build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS)" -o $(OUTPUT) $(PROJECT_PATH)$(REDIS_SINGLE_MAIN)
	cp $(OUTPUT) ./$(MODULE_NAME)

mk_controller-build:
	GOOS=linux go build -ldflags "$(LDFLAGS)" -o $(OUTPUT) $(PROJECT_PATH)$(MYSQL_KAFKA_CONTROLLER)
	cp $(OUTPUT) ./mk_controller

mk_worker-build:
	GOOS=linux go build -ldflags "$(LDFLAGS)" -o $(OUTPUT) $(PROJECT_PATH)$(MYSQL_KAFKA_WORKER)
	cp $(OUTPUT) ./mk_worker

mkc_local-build:
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(PROJECT_PATH)/$(MODULE_NAME) $(PROJECT_PATH)$(MYSQL_KAFKA_CONTROLLER)

mkw_local-build:
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" $(GOFLAGS) -o $(PROJECT_PATH)/$(MODULE_NAME) $(PROJECT_PATH)$(MYSQL_KAFKA_WORKER)

target-version:
	@echo "========================================"
	@echo "APP_VERSION    : $(APP_VERSION)"
	@echo "BUILD_NUM      : $(BUILD_NUM)"
	@echo "TARGET_VERSION : $(TARGET_VERSION)"
	@echo "========================================"

build_num:
	@echo $$(($$(cat $(BUILD_NUM_FILE)) + 1 )) > $(BUILD_NUM_FILE)
	@echo "BUILD_NUM      : $(BUILD_NUM)"

test:
	@go test ./...
