# Monitor Trade

[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://hub.docker.com/r/ddhdocker/monitor-trade)
[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/)
[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://telegram.org/)

一个基于 Go 语言开发的加密货币交易监控系统，支持实时价格监控、Telegram 机器人通知、与 Freqtrade 集成等功能。

## ✨ 功能特性

- 🚀 **实时价格监控**: 通过 Binance WebSocket 获取实时价格数据
- 🤖 **Telegram 机器人**: 支持多种命令进行交易操作和查询
- 📊 **交易集成**: 与 Freqtrade 交易机器人深度集成
- 💾 **数据持久化**: 使用 Redis 进行高性能数据存储
- 🌐 **HTTP API**: 提供 RESTful API 接口
- 🐳 **Docker 支持**: 完整的容器化解决方案
- 📈 **价格预警**: 自定义价格阈值监控
- 🔄 **自动同步**: Redis 与本地数据自动同步

## 🛠️ 技术栈

- **后端**: Go 1.19+
- **数据库**: Redis
- **消息通知**: Telegram Bot API
- **交易接口**: Freqtrade API, Binance API
- **容器化**: Docker & Docker Compose
- **前端**: React (Web 管理界面)

## 🚀 快速开始

### 使用 Docker Compose (推荐)

1. 克隆项目
```bash
git clone https://github.com/your-username/monitor-trade.git
cd monitor-trade
```

2. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，填入必要的配置
```

3. 启动服务
```bash
docker-compose up -d
```

### 手动安装

1. 安装依赖
```bash
# 安装 Go 1.19+
# 安装 Redis
```

2. 构建项目
```bash
make build
```

3. 配置环境变量
```bash
export TELEGRAM_TOKEN="your_telegram_bot_token"
export TELEGRAM_ID="your_telegram_user_id"
export REDIS_ADDR="localhost:6379"
export BOT_BASE_URL="http://localhost:8080"
export BOT_USER_NAME="your_freqtrade_username"
export BOT_PASSWD="your_freqtrade_password"
```

4. 运行程序
```bash
./bin/monitor-trade
```

## ⚙️ 配置说明

### 环境变量

| 变量名 | 描述 | 默认值 | 必填 |
|--------|------|--------|------|
| `TELEGRAM_TOKEN` | Telegram Bot Token | - | ✅ |
| `TELEGRAM_ID` | Telegram 用户 ID | - | ✅ |
| `REDIS_ADDR` | Redis 服务器地址 | `redis:6379` | ❌ |
| `REDIS_PASSWORD` | Redis 密码 | - | ❌ |
| `REDIS_DB` | Redis 数据库编号 | `0` | ❌ |
| `KEY_EXPIRE` | Redis 键过期时间(秒) | `2592000` | ❌ |
| `FUNDING_RATE` | 资金费率阈值 | `-0.1` | ❌ |
| `BOT_BASE_URL` | Freqtrade API 地址 | `http://127.0.0.1:8080` | ❌ |
| `BOT_USER_NAME` | Freqtrade 用户名 | - | ❌ |
| `BOT_PASSWD` | Freqtrade 密码 | - | ❌ |

### Telegram Bot 配置

1. 通过 [@BotFather](https://t.me/BotFather) 创建 Telegram Bot
2. 获取 Bot Token
3. 获取您的 Telegram User ID (可通过 [@userinfobot](https://t.me/userinfobot) 获取)

## 🤖 Telegram 命令

| 命令 | 参数 | 描述 | 示例 |
|------|------|------|------|
| `/s` | `[pair] [price]` | 做空监控 | `/s BTCUSDT 50000` |
| `/l` | `[pair] [price]` | 做多监控 | `/l ETHUSDT 3000` |
| `/c` | `[pair] [direction]` | 取消监控 | `/c BTCUSDT short` |
| `/show` | `[pair]` | 显示监控状态 | `/show BTCUSDT` |
| `/adjust` | - | 显示持仓信息 | `/adjust` |
| `/ad` | `[pair] [amount] [price]` | 添加仓位 | `/ad BTCUSDT 100 50000` |
| `/pc` | `[pair] [amount]` | 部分平仓 | `/pc BTCUSDT 50` |
| `/whitelist` | - | 查看白名单 | `/whitelist` |

## 🌐 HTTP API

### 监控管理

```bash
# 获取所有监控数据
GET /api/monitor

# 添加监控
POST /api/monitor
{
  "pair": "BTCUSDT",
  "price": 50000,
  "direction": "long"
}

# 删除监控
DELETE /api/monitor/{pair}/{direction}
```

### 交易操作

```bash
# 强制卖出
POST /api/v1/forcesell
{
  "tradeid": "36",
  "ordertype": "limit",
  "amount": "20"
}
```

## 🐳 Docker 部署

### 构建镜像

```bash
make docker
```

### 推送镜像

```bash
make push
```

### Docker Compose

```yaml
version: '3.8'

services:
  monitor-trade:
    image: ddhdocker/monitor-trade:latest
    environment:
      - TELEGRAM_TOKEN=${TELEGRAM_TOKEN}
      - TELEGRAM_ID=${TELEGRAM_ID}
      - REDIS_ADDR=redis:6379
      - BOT_BASE_URL=${BOT_BASE_URL}
      - BOT_USER_NAME=${BOT_USER_NAME}
      - BOT_PASSWD=${BOT_PASSWD}
    depends_on:
      - redis
    ports:
      - "8080:8080"

  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  webapp:
    build: ./public/web-app
    ports:
      - "80:80"
    depends_on:
      - monitor-trade

volumes:
  redis_data:
```

## 🏗️ 项目结构

```
monitor-trade/
├── bin/                    # 编译输出
├── config/                 # 配置管理
├── controller/             # 控制器层
│   ├── binance/           # Binance API 集成
│   ├── freqtrade/         # Freqtrade API 集成
│   ├── http/              # HTTP 服务器
│   ├── redis/             # Redis 操作
│   └── tg/                # Telegram Bot
├── model/                  # 数据模型
├── public/web-app/         # Web 前端
├── Dockerfile             # Docker 配置
├── docker-compose.yml     # Docker Compose 配置
├── Makefile              # 构建脚本
└── main.go               # 程序入口
```

## 🔧 开发指南

### 环境要求

- Go 1.19+
- Redis 6.0+
- Node.js 16+ (前端开发)

### 本地开发

```bash
# 安装依赖
go mod download

# 运行测试
go test ./...

# 本地运行
go run main.go

# 前端开发
cd public/web-app
npm install
npm start
```

### 代码规范

- 使用 `gofmt` 格式化代码
- 遵循 Go 官方编码规范
- 添加必要的注释和文档
- 编写单元测试

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

### 提交规范

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

## 📝 许可证

本项目基于 MIT 许可证开源 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Freqtrade](https://github.com/freqtrade/freqtrade) - 优秀的交易机器人框架
- [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) - Telegram Bot API Go 客户端
- [go-redis](https://github.com/go-redis/redis) - Redis Go 客户端

## 📞 联系方式

- 项目主页: [https://github.com/riven-blade/monitor-trade](https://github.com/riven-blade/monitor-trade)
- 问题反馈: [Issues](https://github.com/riven-blade/monitor-trade/issues)

---

⭐ 如果这个项目对您有帮助，请给我们一个 Star！ 