# SponsorGen

一个基于Go的动态赞助者显示生成工具，支持多平台赞助数据聚合与展示。SponsorGen可以生成SVG和JSON格式的赞助者展示，适用于开源项目的README文件或网站，支持GitHub、OpenCollective、Patreon和Afdian等平台。

**请注意：所有代码(包括本README文档)均由 AI 生成，不保证其安全性和可用性，请酌情使用**

## 效果预览

下面是一个示例SVG赞助者展示：

![赞助者展示示例](.github/images/example.svg)

## 功能特性

- 多平台支持：集成GitHub Sponsors、OpenCollective、Patreon和Afdian等赞助平台
- 灵活配置：通过环境变量管理所有设置
- 动态更新：支持基于时间间隔的自动刷新和每日凌晨00:00的定时刷新
- 多格式输出：生成SVG图像和JSON数据
- 自定义样式：支持自定义字体、颜色、尺寸等显示参数
- Docker支持：提供容器化部署方案
- 缓存机制：减少API请求，提高性能

## 快速开始

### 使用Docker

```bash
# 拉取镜像
docker pull ghcr.io/versun/sponsorgen:latest

# 运行容器（使用环境变量配置）
docker run -p 5000:5000 \
  -e GITHUB_TOKEN=your_github_token \
  -e GITHUB_LOGIN=your_github_login \
  -e AFDIAN_USER_ID=your_afdian_user_id \
  -e AFDIAN_TOKEN=your_afdian_token \
  ghcr.io/versun/sponsorgen:latest
```

或者使用docker-compose:

```bash
# 使用docker-compose启动服务
docker-compose up -d
```

### 本地运行

1. 克隆仓库

```bash
git clone https://github.com/versun/sponsorgen.git
cd sponsorgen
```

2. 设置环境变量（示例）

```bash
# 输出设置
export OUTPUT_DIR="./output"
export CACHE_DIR="./cache"
export REFRESH_MINUTES="60"
export DEFAULT_AVATAR="./assets/default_avatar.svg"

# GitHub赞助设置（可选）
export GITHUB_TOKEN="your_github_token"
export GITHUB_LOGIN="your_github_login"

# Afdian赞助设置（可选）
export AFDIAN_USER_ID="your_afdian_user_id"
export AFDIAN_TOKEN="your_afdian_token"

# 渲染设置（可选）
export AVATAR_SIZE="45"
export SVG_WIDTH="800"
export SHOW_AMOUNT="false"
export SHOW_NAME="false"
```

3. 构建并运行

```bash
go build -o sponsorgen
./sponsorgen -port 5000
```

4. 访问服务

- SVG输出: http://localhost:5000/sponsors.svg
- JSON输出: http://localhost:5000/sponsors.json
- 强制刷新: http://localhost:5000/refresh

## 配置选项

以下是可用的环境变量配置选项：

| 环境变量 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| OUTPUT_DIR | string | "./output" | 输出文件目录 |
| CACHE_DIR | string | "./cache" | 缓存文件目录 |
| REFRESH_MINUTES | int | 60 | 自动刷新间隔（分钟） |
| DEFAULT_AVATAR | string | "./assets/default_avatar.svg" | 默认头像路径 |
| GITHUB_TOKEN | string | "" | GitHub Personal Access Token |
| GITHUB_LOGIN | string | "" | GitHub用户名 |
| INCLUDE_PRIVATE | bool | false | 是否包含私人赞助者 |
| GITHUB_ORGS | string | "" | 包含的GitHub组织，用逗号分隔 |
| EXCLUDE_SPONSORS | string | "" | 排除的赞助者，用逗号分隔 |
| INCLUDE_SPONSORS | string | "" | 强制包含的赞助者，用逗号分隔 |
| OPENCOLLECTIVE_SLUG | string | "" | OpenCollective项目标识 |
| OPENCOLLECTIVE_KEY | string | "" | OpenCollective API密钥 |
| PATREON_TOKEN | string | "" | Patreon访问令牌 |
| PATREON_CAMPAIGN_ID | string | "" | Patreon活动ID |
| AFDIAN_USER_ID | string | "" | 爱发电用户ID |
| AFDIAN_TOKEN | string | "" | 爱发电TOKEN |
| AVATAR_SIZE | int | 45 | 头像尺寸（像素） |
| AVATAR_MARGIN | int | 5 | 头像间距（像素） |
| SVG_WIDTH | int | 800 | SVG宽度（像素） |
| FONT_SIZE | int | 14 | 字体大小（像素） |
| FONT_FAMILY | string | "system-ui..." | 字体系列 |
| SHOW_AMOUNT | bool | false | 是否显示赞助金额 |
| SHOW_NAME | bool | false | 是否显示赞助者名称 |
| BACKGROUND_COLOR | string | "transparent" | 背景颜色 |
| PADDING_X | int | 10 | X轴内边距（像素） |
| PADDING_Y | int | 10 | Y轴内边距（像素） |

## 在GitHub README中使用

将以下内容添加到您的README.md文件中：

```markdown
## 赞助者

![赞助者](https://your-sponsorgen-url.com/sponsors.svg)
```

## Docker镜像发布

本项目使用GitHub Actions自动构建并发布Docker镜像。详情请查看 [Docker发布指南](.github/docker-publish-guide.md)。

## API端点

| 端点 | 方法 | 说明 |
|------|------|------|
| / | GET | 首页，显示简单的使用说明 |
| /sponsors.svg | GET | 生成并返回赞助者SVG |
| /sponsors.json | GET | 返回赞助者JSON数据 |
| /refresh | GET | 强制刷新赞助者数据 |
| /static/* | GET | 访问生成的静态文件 |

## 贡献指南

欢迎提交Pull Request或Issue！以下是一些贡献指南：

1. Fork仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的改动 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个Pull Request

## 许可证

此项目采用MIT许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

## 技术栈

- Go 1.19+
- Docker/Docker Compose
- GitHub Actions (CI/CD)