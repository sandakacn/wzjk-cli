# wzjk-cli - 网站监控系统 CLI 工具

`wzjk-cli` 是网站监控系统的命令行工具，用于管理域名监控、查看 SSL 证书状态等。

## 特性

- 🔐 API Key 认证
- 🌐 域名管理（添加、删除、查看）
- 🔒 SSL 证书检查
- 📊 监控状态概览
- 🎨 彩色终端输出
- 📱 跨平台支持（Linux/macOS/Windows）

## 安装

### Homebrew（推荐）

**macOS / Linux：**

```bash
brew tap sandakacn/wzjk-cli
brew install wzjk-cli
```

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/sandakacn/wzjk-cli.git
cd wzjk-cli

# 构建
make build

# 安装到 GOPATH/bin
make install
```

### 预编译二进制文件

从 [GitHub Releases](https://github.com/sandakacn/wzjk-cli/releases) 下载对应平台的二进制文件。

## 快速开始

### 1. 生成 API Key

在网页版登录后，访问个人资料页面 (`/profile`)，在 "CLI 工具 (API Keys)" 部分：

1. 点击"生成新 API Key"
2. 输入名称（可选，用于标识）
3. 选择有效期（默认 30 天）
4. 点击"生成"
5. **立即复制 API Key**（关闭弹窗后将无法再次查看）

### 2. 登录 CLI

```bash
# 交互式登录（会提示输入 API Key）
wzjk-cli login --api-url https://wangzhanjiankong.cn

# 或直接指定 API Key
wzjk-cli login --api-url https://wangzhanjiankong.cn --api-key <your-api-key>
```

API Key 将保存在 `~/.wzjk-cli/config.yaml` 文件中。

### 3. 使用 CLI

```bash
# 查看域名列表
wzjk-cli domains list

# 添加域名
wzjk-cli domains add example.com

# 检查 SSL 证书
wzjk-cli domains check example.com

# 查看监控状态概览
wzjk-cli status

# 查看用户信息
wzjk-cli profile
```

## 命令参考

### 登录管理

```bash
# 登录（交互式）
wzjk-cli login --api-url <url>

# 直接指定 API Key
wzjk-cli login --api-url <url> --api-key <key>
wzjk-cli login --api-url <url> --token <key>

# 退出登录
wzjk-cli logout [--force]
```

### 域名管理

```bash
# 列出所有域名
wzjk-cli domains list [--format json] [--alerts-only]

# 添加域名
wzjk-cli domains add <domain> [--port <port>] [--alert-days <days>] [--type <type>]

# 删除域名
wzjk-cli domains delete <domain-id> [--force]

# 更新域名设置
wzjk-cli domains update <domain-id> [--alert-days <days>] [--active <true|false>]

# 检查 SSL 证书
wzjk-cli domains check <domain> [--port <port>]
```

### 状态查看

```bash
# 查看监控状态概览
wzjk-cli status

# 查看用户信息
wzjk-cli profile
```

## 配置

配置文件存储在 `~/.wzjk-cli/config.yaml`：

```yaml
api_url: https://wangzhanjiankong.cn
token: <your-api-key>
user:
  id: <user-id>
  name: <user-name>
  email: <user-email>
```

## API Key 管理

- 在网页版个人资料页面可以查看所有 API Key
- 每个 Key 显示创建时间、最后使用时间、过期时间
- 可以随时删除（吊销）API Key
- 过期或删除的 Key 将无法再用于 CLI 登录

## 安全建议

1. **妥善保管 API Key**：就像密码一样，不要分享给他人
2. **定期轮换**：建议定期删除旧 Key，生成新 Key
3. **设置合理的有效期**：根据使用场景设置过期时间
4. **使用 HTTPS**：确保 `--api-url` 使用 HTTPS 协议
5. **及时吊销**：如果怀疑 Key 泄露，立即在网页版删除

## 开发

```bash
# 安装依赖
cd wzjk-cli
go mod download

# 开发模式运行
go run main.go <command>

# 构建
make build

# 运行测试
make test

# 为所有平台构建
make build-all
```

## 许可证

MIT
