PROJECT_PATH=$(shell pwd)
MODULE_NAME=live-chat-server

BUILD_NUM_FILE=build_num.txt
BUILD_NUM=$$(cat ./build_num.txt)
APP_VERSION=0.0
TARGET_VERSION=$(APP_VERSION).$(BUILD_NUM)
TARGET_DIR=bin
OUTPUT=$(PROJECT_PATH)/$(TARGET_DIR)/$(MODULE_NAME)
MAIN_FILE=/main.go

LDFLAGS=-X main.BUILD_TIME=`date -u '+%Y-%m-%d_%H:%M:%S'`
LDFLAGS+=-X main.APP_VERSION=$(TARGET_VERSION)
LDFLAGS+=-X main.GIT_HASH=`git rev-parse HEAD`
LDFLAGS+=-s -w

all: config test build

config:
	@if [ ! -d $(TARGET_DIR) ]; then mkdir $(TARGET_DIR); fi

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS)" -o $(OUTPUT) $(PROJECT_PATH)$(MAIN_FILE)
	cp $(OUTPUT) ./$(MODULE_NAME)

local-build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(OUTPUT) $(PROJECT_PATH)$(MAIN_FILE)
	mv $(OUTPUT) ./$(MODULE_NAME)

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
