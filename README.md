# Stock Monitor (Go + Fiber + MySQL)

一个用于拉取 AlphaVantage 周线数据、计算 52 周指标并展示看板的后端项目。

## 应用信息

### 技术栈
- 后端：Go 1.25 + Fiber
- 数据源：AlphaVantage `TIME_SERIES_WEEKLY_ADJUSTED`
- 数据库：MySQL 8
- 模板渲染：Fiber HTML Template
- 定时任务：`robfig/cron/v3`
- 容器部署：Docker + Docker Compose

### 核心能力
- 手动同步单个股票：拉取外部行情并入库
- 定时同步股票池：按 cron 自动更新
- 看板展示：渲染每个 symbol 的最新记录

### 指标说明
- `last_close`：最近一周 adjusted close
- `high_52w`：最近 52 周最高价
- `low_52w`：最近 52 周最低价
- `percent_diff`：`(last_close - high_52w) / high_52w * 100`

### HTTP 路由
- `GET /health`：健康检查
- `GET /`：看板页面
- `GET /stocks/:id`：按 ID 查询
- `POST /stocks/sync/:symbol`：手动触发同步

## 环境变量

| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| `PORT` | 服务端口 | `8080` |
| `LOG_LEVEL` | 日志级别 | `info` |
| `APP_NAME` | 应用名 | `stock-app` |
| `DB_USER` | MySQL 用户 | `root` |
| `DB_PASSWORD` | MySQL 密码 | 空 |
| `DB_HOST` | MySQL 地址 | `127.0.0.1` |
| `DB_PORT` | MySQL 端口 | `3306` |
| `DB_NAME` | MySQL 库名 | `stock_app` |
| `ALPHA_API_KEY` | AlphaVantage API Key | `test_api_key` |
| `ALPHA_BASE_URL` | AlphaVantage 地址 | `https://www.alphavantage.co/query` |
| `SYNC_SYMBOLS` | 定时同步股票池（逗号分隔） | `MSFT,AAPL,META` |
| `SYNC_CRON` | cron 表达式（5 段） | `*/30 * * * *` |
| `SYNC_RUN_ON_START` | 服务启动时是否先执行一轮 | `true` |

## 数据库

建表 SQL：`scripts/stock.sql`

说明：
- 当前使用唯一键 `(symbol, cur_date)`，同一股票同一天重复同步会更新，不会重复插入。
- `docker-compose` 首次初始化（空数据卷）会自动执行 `scripts/stock.sql`。

## 本地开发启动

1. 准备 `.env`（可参考 `.env.docker.example`）
2. 启动 MySQL 并执行建表 SQL
3. 启动服务：

```bash
go run ./cmd/main.go
```

手动测试：

```bash
curl --noproxy '*' http://127.0.0.1:8080/health
curl --noproxy '*' -X POST http://127.0.0.1:8080/stocks/sync/MSFT
```

## 云服务器 Docker 部署

下面流程适用于 Ubuntu/CentOS 等 Linux 云主机。

### 1. 安装 Docker 与 Compose

Ubuntu 参考：

```bash
sudo apt-get update
sudo apt-get install -y docker.io docker-compose-plugin
sudo systemctl enable --now docker
```

### 2. 拉取代码

```bash
git clone <你的仓库地址> stock-app
cd stock-app
```

### 3. 准备部署环境变量

```bash
cp .env.docker.example .env.docker
```

编辑 `.env.docker`，至少修改：
- `MYSQL_ROOT_PASSWORD`
- `DB_PASSWORD`
- `ALPHA_API_KEY`

### 4. 启动容器

```bash
docker compose --env-file .env.docker up -d --build
```

查看状态：

```bash
docker compose ps
docker compose logs -f app
```

### 5. 验证部署

```bash
curl http://<服务器IP>:8080/health
```

浏览器访问：
- `http://<服务器IP>:8080/` 看板页面

### 6. 云安全组/防火墙

放通：
- `8080/tcp`（应用）
- `3306/tcp`（仅在确实需要外部访问数据库时开放，建议默认不开放）

### 7. 升级发布

```bash
git pull
docker compose --env-file .env.docker up -d --build
```

### 8. 停止与清理

停止服务：

```bash
docker compose down
```

删除容器并清空数据库数据卷（危险操作）：

```bash
docker compose down -v
```

## 常见问题

### 1) `bind: address already in use`
端口占用，先查并杀进程：

```bash
lsof -nP -iTCP:8080 -sTCP:LISTEN
kill <PID>
```

### 2) MySQL 初始化脚本没生效
`docker-entrypoint-initdb.d` 只在空数据目录首次执行。
如果已存在旧数据卷，需要：
- 手动执行 SQL，或
- `docker compose down -v` 后重新启动（会清空数据）。

### 3) AlphaVantage 限流
免费额度较低。可通过以下方式缓解：
- 减少 `SYNC_SYMBOLS`
- 拉长 `SYNC_CRON` 间隔
- 保持当前串行同步策略
