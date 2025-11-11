# API接口文档

本文档描述了Moment项目的API接口规范，包括API网关和CodeShare服务的接口定义。

## 项目架构概述

项目采用微服务架构，主要包含以下组件：
- **API网关**：统一入口，负责请求路由和负载均衡
- **CodeShare服务**：代码片段分享服务，支持代码上传和获取

## 1. API网关

### 网关地址
- 地址：`http://localhost:8080`
- 端口可在配置文件`internal/config/config.yaml`中修改

### 网关功能
- 请求路由：根据路径前缀将请求转发到对应的后端服务
- 负载均衡：支持权重配置的负载均衡策略

## 2. CodeShare服务接口

CodeShare服务提供代码片段的上传和获取功能，支持直接访问或通过API网关访问。

### 2.1 上传代码片段

**直接访问**：
- 请求方法：`POST`
- 请求路径：`http://localhost:8081/upload`（或其他端口实例）

**通过API网关访问**：
- 请求方法：`POST`
- 请求路径：`http://localhost:8080/CodeShare/upload`

**请求头**：
- `Content-Type: application/json`

**请求体**：
```json
{
    "author": "作者名称",
    "language": "代码语言",
    "content": "代码内容",
    "ttl": 过期时间（秒，可选，默认3600）
}
```

**请求参数说明**：
| 参数名 | 类型 | 是否必选 | 描述 |
|-------|------|---------|------|
| author | string | 是 | 代码作者名称 |
| language | string | 是 | 代码语言类型（如python、javascript、java等） |
| content | string | 是 | 代码内容，最大支持100000字符 |
| ttl | int64 | 否 | 代码片段有效期（秒），默认为3600秒（1小时） |

**响应格式**：
```json
{
    "message": "success",
    "hash": "代码哈希值",
    "url": "/code/哈希值"
}
```

**响应参数说明**：
| 参数名 | 类型 | 描述 |
|-------|------|------|
| message | string | 操作结果信息，成功时为"success" |
| hash | string | 代码片段的唯一标识哈希值 |
| url | string | 访问该代码片段的URL路径 |

**状态码**：
- 200：上传成功
- 400：请求体格式错误
- 405：请求方法不正确

### 2.2 获取代码片段

**直接访问**：
- 请求方法：`GET`
- 请求路径：`http://localhost:8081/code/{hash}`（或其他端口实例）

**通过API网关访问**：
- 请求方法：`GET`
- 请求路径：`http://localhost:8080/CodeShare/code/{hash}`

**请求参数**：
- `hash`：URL路径参数，表示代码片段的唯一标识哈希值

**响应格式**：
```json
{
    "author": "作者名称",
    "language": "代码语言",
    "content": "代码内容",
    "hash": "代码哈希值",
    "destroy_time": 销毁时间戳
}
```

**响应参数说明**：
| 参数名 | 类型 | 描述 |
|-------|------|------|
| author | string | 代码作者名称 |
| language | string | 代码语言类型 |
| content | string | 代码内容 |
| hash | string | 代码片段的唯一标识哈希值 |
| destroy_time | int64 | 代码片段的销毁时间戳（Unix时间） |

**状态码**：
- 200：获取成功
- 404：代码片段不存在或已过期

## 3. 后端服务配置

根据`internal/config/config.yaml`文件，当前系统配置了3个CodeShare服务实例，通过API网关进行负载均衡：

| 路径前缀 | 目标地址 | 权重 |
|---------|---------|------|
| /CodeShare | http://localhost:8081/CodeShare | 3 |
| /CodeShare | http://localhost:8082/CodeShare | 2 |
| /CodeShare | http://localhost:8083/CodeShare | 1 |

## 4. 数据存储

CodeShare服务支持两种存储模式：
- **内存存储**：默认模式，适合开发和测试环境
- **Redis存储**：通过环境变量配置，适合生产环境

**Redis配置环境变量**：
- `REDIS_ADDR`：Redis服务器地址（如"127.0.0.1:6379"）
- `REDIS_PASSWORD`：Redis密码（可选）
- `REDIS_DB`：Redis数据库索引（默认0）

## 5. 代码结构参考

### 主要模块
- `internal/Gateway/`：API网关实现
- `internal/CodeShare/`：代码分享服务实现
- `internal/Server/`：HTTP服务器封装
- `internal/config/`：配置管理

### 核心接口定义
- **Gateway结构体**：<mcfile name="gateway.go" path="/home/leon/GoCode/Moment/backend/internal/Gateway/gateway.go"></mcfile>
- **CodeShare结构体**：<mcfile name="main.go" path="/home/leon/GoCode/Moment/backend/internal/CodeShare/main.go"></mcfile>
- **处理器函数**：<mcfile name="handler.go" path="/home/leon/GoCode/Moment/backend/internal/CodeShare/handler.go"></mcfile>
        