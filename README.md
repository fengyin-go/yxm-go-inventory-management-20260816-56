# 库存管理系统（Inventory Management）

一个纯 Go 标准库实现的库存管理 REST API 服务，采用标准 Go 工程目录结构，内存存储，零第三方依赖。

## 目录结构

```
origin/
├── cmd/server/          # 程序入口
├── internal/
│   ├── app/             # 依赖装配
│   ├── config/          # 配置加载
│   ├── model/           # 领域模型与校验
│   ├── store/           # 数据访问接口 + 内存实现
│   ├── service/         # 业务逻辑层
│   └── handler/         # HTTP 处理器层
└── pkg/
    ├── httpx/           # HTTP 响应工具
    ├── idgen/           # ID / 短码生成
    └── logger/          # 分级日志
```

## 运行

```bash
go run ./cmd/server
# 或指定端口
PORT=8081 go run ./cmd/server
```

默认监听 `:8080`。可用环境变量：`PORT` / `ADDR`、`LOW_STOCK_THRESHOLD`（低库存阈值，默认 10）、`MAX_PAGE_SIZE`、`LOG_LEVEL`。

## 测试

```bash
go test ./...
```

## API 接口

统一响应：`{"code": 0, "message": "ok", "data": ...}`；分页响应：`{"items": [...], "pagination": {"page":1,"size":20,"total":N}}`。

### 商品

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/products` | 创建商品 |
| GET | `/api/products?category=&status=&keyword=&page=&size=` | 商品列表 |
| GET | `/api/products/{id}` | 商品详情 |
| PUT | `/api/products/{id}` | 更新商品 |
| DELETE | `/api/products/{id}` | 删除商品 |
| PATCH | `/api/products/{id}/status` | 上架/下架 `{"status":"active"}` |

### 仓库

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/warehouses` | 创建仓库 |
| GET | `/api/warehouses` | 仓库列表 |
| GET | `/api/warehouses/{id}` | 仓库详情 |
| PUT | `/api/warehouses/{id}` | 更新仓库 |
| DELETE | `/api/warehouses/{id}` | 删除仓库 |

### 库存

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/stocks` | 建立库存记录 |
| GET | `/api/stocks?product_id=&warehouse_id=&page=&size=` | 库存列表 |
| GET | `/api/stocks/{id}` | 库存详情 |
| PATCH | `/api/stocks/{id}/threshold` | 调整预警阈值 |
| DELETE | `/api/stocks/{id}` | 删除库存项 |
| POST | `/api/stocks/in` | 入库 `{"product_id","warehouse_id","quantity","operator","remark"}` |
| POST | `/api/stocks/out` | 出库（校验库存充足） |
| GET | `/api/stocks/low` | 低库存列表 |

### 流水与统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/movements?product_id=&warehouse_id=&type=&page=&size=` | 出入库流水 |
| GET | `/api/movements/{id}` | 流水详情 |
| DELETE | `/api/movements/{id}` | 删除流水 |
| GET | `/api/stats` | 全局统计 |
| GET | `/api/stats/warehouses` | 按仓库汇总 |
