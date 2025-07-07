# 贡献指南

感谢您对 Monitor Trade 项目的关注和贡献！我们欢迎所有形式的贡献，包括但不限于：

- 🐛 报告 Bug
- 💡 提出新功能建议
- 📝 改进文档
- 💻 提交代码修复或新功能
- 🧪 编写测试用例

## 🚀 开始贡献

### 1. Fork 项目

点击页面右上角的 "Fork" 按钮，将项目 fork 到您的 GitHub 账户。

### 2. 克隆代码

```bash
git clone https://github.com/your-username/monitor-trade.git
cd monitor-trade
```

### 3. 创建分支

为您的更改创建一个新分支：

```bash
git checkout -b feature/your-feature-name
# 或者
git checkout -b fix/your-bug-fix
```

### 4. 配置开发环境

```bash
# 安装 Go 依赖
go mod download

# 启动开发环境
docker-compose up -d redis

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，填入您的测试配置
```

## 📝 代码规范

### Go 代码规范

- 使用 `gofmt` 格式化代码
- 遵循 [Effective Go](https://golang.org/doc/effective_go.html) 指南
- 为公共函数和类型添加注释
- 变量和函数命名使用驼峰命名法
- 包名使用小写字母

### 提交消息规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**类型 (type):**
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 仅文档更改
- `style`: 代码格式调整（不影响功能）
- `refactor`: 代码重构
- `perf`: 性能优化
- `test`: 添加测试
- `chore`: 构建过程或辅助工具的变动

**示例:**
```
feat(tg): 添加新的 /balance 命令

添加查询账户余额的 Telegram 命令，支持显示：
- USDT 余额
- 未实现盈亏
- 可用保证金

Closes #123
```

## 🧪 测试

在提交代码前，请确保所有测试通过：

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./controller/tg

# 运行测试并显示覆盖率
go test -cover ./...
```

## 📋 提交 Pull Request

1. **确保代码质量**
   - 代码格式化：`gofmt -w .`
   - 运行测试：`go test ./...`
   - 检查构建：`make build`

2. **推送更改**
   ```bash
   git add .
   git commit -m "feat: 添加新功能"
   git push origin feature/your-feature-name
   ```

3. **创建 Pull Request**
   - 在 GitHub 上创建 Pull Request
   - 使用清晰的标题和描述
   - 链接相关的 Issue
   - 添加截图或演示（如适用）

### Pull Request 模板

```markdown
## 变更类型
- [ ] 新功能
- [ ] Bug 修复
- [ ] 文档更新
- [ ] 代码重构
- [ ] 性能优化

## 变更描述
简要描述您的更改...

## 相关 Issue
Closes #(issue number)

## 测试
- [ ] 已运行现有测试
- [ ] 已添加新测试
- [ ] 手动测试通过

## 截图
如果适用，请添加截图来说明您的更改。
```

## 🐛 报告 Bug

使用 [GitHub Issues](https://github.com/your-username/monitor-trade/issues) 报告 bug，请包含：

- **环境信息**: OS、Go 版本、Docker 版本等
- **重现步骤**: 详细的重现步骤
- **期望行为**: 您期望发生什么
- **实际行为**: 实际发生了什么
- **日志**: 相关的错误日志
- **截图**: 如果适用

## 💡 功能建议

我们欢迎新功能建议！请先创建一个 Issue 来讨论您的想法：

- 描述功能的用途和好处
- 提供使用场景
- 考虑可能的实现方式
- 讨论对现有功能的影响

## 📖 文档贡献

文档改进同样重要：

- 修正拼写错误
- 改进说明的清晰度
- 添加使用示例
- 翻译文档到其他语言

## 🤝 代码审查

所有的 Pull Request 都需要经过代码审查：

- 至少需要一名维护者的批准
- 确保 CI 检查通过
- 解决所有审查意见
- 保持提交历史的整洁

## 📞 获取帮助

如果您在贡献过程中遇到问题：

- 查看已有的 [Issues](https://github.com/your-username/monitor-trade/issues)
- 创建新的 Issue 提问
- 加入我们的社区讨论

## 🙏 致谢

感谢所有为项目做出贡献的开发者！您的每一个贡献都让项目变得更好。

---

再次感谢您的贡献！ 🎉 