# 项目基础配置
APP_NAME := monitor-trade
VERSION  := v1.0.148
DOCKER_REPO := ddhdocker/$(APP_NAME)
FRONTEND_REPO := ddhdocker/$(APP_NAME)-frontend

# Go 编译器配置
GOOS   := linux
GOARCH := amd64
BIN    := bin/$(APP_NAME)

# 前端配置
FRONTEND_DIR := public/web-app
FRONTEND_BUILD_DIR := $(FRONTEND_DIR)/build
PUBLIC_DIR := public

.PHONY: all build docker push clean frontend build-all frontend-docker frontend-push frontend-app

# 默认任务：构建前端 + 后端 + 镜像
all: build-all docker push

# 构建所有（前端 + 后端）
build-all: frontend build

# 编译 Go 项目
build:
	@echo "==> Building Go binary..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o $(BIN) main.go

# 构建前端项目
frontend:
	@echo "==> Building frontend..."
	cd $(FRONTEND_DIR) && npm install && npm run build
	@echo "==> Copying frontend build files to public directory..."
	cp -r $(FRONTEND_BUILD_DIR)/* $(PUBLIC_DIR)/

# 构建前端 Docker 镜像
frontend-docker:
	@echo "==> Building Frontend Docker image..."
	docker build --platform=linux/amd64 -t $(FRONTEND_REPO):$(VERSION) $(FRONTEND_DIR)

# 推送前端 Docker 镜像
frontend-push:
	@echo "==> Pushing Frontend Docker image to repository..."
	docker push $(FRONTEND_REPO):$(VERSION)

# 构建前端独立应用（构建+镜像+推送）
frontend-app: frontend-docker frontend-push

# 构建 Docker 镜像
docker:
	@echo "==> Building Docker image..."
	docker build --platform=linux/amd64 -t $(DOCKER_REPO):$(VERSION) .

# 推送到 Docker 仓库
push:
	@echo "==> Pushing Docker image to repository..."
	docker push $(DOCKER_REPO):$(VERSION)

# 清理构建产物
clean:
	@echo "==> Cleaning up..."
	rm -rf bin/
	rm -rf $(FRONTEND_BUILD_DIR)

