# Changelog

本文档记录了所有 Monitor Trade 项目的重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### 计划中
- 增加更多交易所支持
- Web 界面优化
- 移动端应用
- 高级分析功能

## [1.0.132] - 2024-12-19

### Changed
- 优化 Redis key 过期时间为随机值，避免缓存雪崩
- 过期时间范围：KeyExpire/2 到 KeyExpire 之间

### Fixed
- 修复大量 key 同时过期导致的性能问题

## [1.0.129] - 2024-12-19

### Added
- 完善项目文档，准备开源
- 添加 README.md 详细说明
- 添加 Docker Compose 配置
- 添加 MIT 开源许可证

### Changed
- 清理配置文件，移除不需要的 TelegramActionToken
- 统一 Telegram 控制器架构

## [1.0.128] - 2024-12-19

### Added
- 新增 `/ad` 命令：添加仓位 (原 adjust 命令功能)
- 新增 `/pc` 命令：部分平仓功能
- 集成 Freqtrade ForceSell API

### Changed
- `/adjust` 命令改为显示仓位信息
- 优化命令结构和用户体验

### Fixed
- 修复 Telegram 格式化字符串错误 (%f.2 -> %.2f)
- 修复 "chat not found" 错误

## [1.0.120] - 2024-12-18

### Added
- 重构项目架构，模块化设计
- 分离 controller 包为多个子包：
  - `controller/tg/` - Telegram 功能
  - `controller/redis/` - Redis 操作
  - `controller/freqtrade/` - Freqtrade 集成
  - `controller/binance/` - Binance 集成
  - `controller/http/` - HTTP 服务器

### Changed
- 优化代码组织结构
- 改进错误处理
- 统一命名规范

### Fixed
- 修复包导入路径问题
- 解决编译错误

## [1.0.100] - 2024-12-15

### Added
- 基础项目架构
- Telegram Bot 集成
- Redis 数据存储
- Binance WebSocket 价格监控
- Freqtrade API 集成
- 基本的价格监控功能

### Features
- 实时价格监控
- Telegram 命令支持：
  - `/s` - 做空监控
  - `/l` - 做多监控  
  - `/c` - 取消监控
  - `/show` - 显示状态
  - `/whitelist` - 查看白名单
- HTTP API 接口
- Docker 容器化支持

---

## 版本说明

- **Added** - 新增功能
- **Changed** - 已有功能的变更
- **Deprecated** - 即将移除的功能
- **Removed** - 已移除的功能
- **Fixed** - 修复的 bug
- **Security** - 安全相关的修复

## 贡献

如果您发现任何遗漏的变更或有建议，请创建 [Issue](https://github.com/your-username/monitor-trade/issues) 告诉我们。 