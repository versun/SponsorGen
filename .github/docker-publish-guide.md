# GitHub Actions Docker 发布指南

## 自动 Docker 镜像发布流程

该项目已配置 GitHub Actions 自动构建并将 Docker 镜像推送到 GitHub Container Registry (ghcr.io)。以下是它的工作原理：

### 触发条件

Docker 镜像将在以下情况下自动构建和发布：

1. 当代码推送到 `main` 分支时
2. 当发布新的带有版本标签的版本时（遵循 `v*.*.*` 格式，如 `v1.0.0`）
3. 当创建拉取请求时（仅构建，不推送）

### 镜像标签策略

工作流会自动为镜像生成以下标签：

- 对于分支推送：使用分支名称（例如 `main`）
- 对于标签推送：
  - 完整版本号（例如 `v1.2.3`）
  - 主要和次要版本（例如 `1.2`）
  - 主要版本（例如 `1`）
- 对于所有推送：使用提交 SHA

### 权限设置

GitHub Actions 使用仓库的 `GITHUB_TOKEN` 进行身份验证，这个令牌会自动提供。该令牌必须有足够的权限来推送到 GitHub Packages。

在仓库设置中确保：

1. 在 Settings > Actions > General 下，将工作流权限设置为"读取和写入权限"
2. 在 Settings > Packages 下，确保已启用 GitHub Packages

### 如何手动触发发布

要手动触发镜像构建和发布：

1. 创建一个新的 Git 标签，遵循语义化版本规范：
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. 或者从 GitHub 网页界面创建新的发布版本：
   - 转到仓库的 "Releases" 标签
   - 点击 "Draft a new release"
   - 输入符合 `v*.*.*` 格式的标签
   - 发布后会自动触发工作流

### 如何获取已发布的镜像

发布后，您可以使用以下命令拉取镜像：

```bash
docker pull ghcr.io/USERNAME/REPOSITORY:TAG
```

例如：

```bash
docker pull ghcr.io/your-username/sponsorgen:v1.0.0
```

您可以在仓库的 "Packages" 标签页中找到所有已发布的镜像。

### 故障排除

如果工作流失败：

1. 检查 Actions 标签页中的工作流运行日志
2. 确保仓库设置中配置了适当的权限
3. 验证 Dockerfile 是否存在于仓库的根目录中