# Regman (Registry Manager) 🐳

`regman` 是一个用 Go 编写的轻量级 CLI 工具，专门用于管理 Docker V2 私有镜像仓库。它支持列出仓库、查询标签以及彻底删除镜像（包括清理后端存储索引）。

## ✨ 特性

- **多维配置**：支持命令行 Flag、环境变量 (`REGMAN_*`) 和配置文件 (`~/.regman.yaml`)。
- **智能认证**：
  - 支持 Basic Auth (用户名/密码)。
  - **自动集成 Docker 凭据**：如果未提供账号，自动读取 `~/.docker/config.json` 中的登录信息。
- **合规删除**：严格遵循 Docker Registry V2 API，先获取 Manifest Digest 再执行删除。
- **自动化维护**：附带用于服务器端垃圾回收（GC）和空索引清理的通用 Ansible Playbook。

## 🚀 安装

确保你已安装 Go 1.21+ 环境：

```bash
git clone <your-repo-url>
cd regman
go build -o regman
```

## ⚙️ 配置

你可以通过以下三种方式配置 `regman`（优先级从高到低）：

1. **命令行参数**：`--registry`, `--user`, `--pass`, `--insecure`
2. **环境变量**：
   ```bash
   export REGMAN_REGISTRY="my-registry.com"
   export REGMAN_USER="admin"
   export REGMAN_PASS="password"
   ```
3. **配置文件** (`~/.regman.yaml`)：
   ```yaml
   registry: "my-registry.com"
   user: "admin"
   pass: "password"
   insecure: true
   ```

## 📖 使用指南

### 列出所有仓库
```bash
./regman ls
```

### 查看镜像的所有标签
```bash
./regman tags my-app
```

### 删除指定镜像/标签
注意：删除操作需要服务器端开启 `DELETE` 功能。
```bash
./regman rm my-app:v1.0.1
```

---

## 🛠️ 服务器端维护 (重要)

### 1. 开启删除功能
默认情况下，Docker Registry 禁用删除操作。请在服务器的 `compose.yml` 中添加：

```yaml
services:
  registry:
    environment:
      REGISTRY_STORAGE_DELETE_ENABLED: "true"
```

### 2. 垃圾回收与索引清理
执行 `rm` 命令仅删除了逻辑引用。要释放磁盘空间并从 `ls` 列表中移除已删空的仓库名，需要运行维护脚本。

我们在 `ansible/` 目录下提供了一个通用自动化方案：

**部署维护任务：**
```bash
ansible-playbook ansible/setup-registry-gc.yml \
  -e "target_hosts=your-registry-server" \
  -e "registry_project_path=/path/to/registry/compose" \
  -e "registry_storage_path=/path/to/registry/data"
```

**变量说明：**
- `target_hosts`: 你主机清单中的主机组 (默认: `all`)。
- `registry_project_path`: 包含 `compose.yml` 的远程目录 (默认: `/root/registry`)。
- `registry_storage_path`: 宿主机上 Registry 存储的根目录 (默认: `/data/registry`)。

## 🏗️ 项目结构

```text
/regman
  /cmd          # Cobra 命令实现 (ls, tags, rm, root)
  /ansible      # 自动化维护工具 (Playbook)
  main.go       # 程序入口
  design.md     # 设计规格文档
  go.mod        # 依赖管理
```

## 📄 开源协议

MIT License
