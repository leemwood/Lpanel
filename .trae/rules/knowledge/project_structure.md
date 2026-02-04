---
alwaysApply: false
description: Lpanel 项目目录结构说明
---
# Lpanel 项目结构

## 目录布局
- `backend/`: Golang 后端源代码。
  - `main.go`: 程序入口，负责路由分发与服务启动。
  - `models/`: 数据库模型定义。
  - `handlers/`: API 接口业务逻辑。
  - `config/`: 数据库与系统配置。
  - `utils/`: 通用工具类（如 Nginx 配置生成）。
- `frontend/`: 前端静态资源。
  - `index.html`: 主页面。
  - `css/style.css`: 样式表。
  - `js/app.js`: 前端逻辑处理。
- `data/`: 持久化数据存储。
  - `lpanel.db`: SQLite 数据库文件。
  - `nginx/conf.d/`: 自动生成的 Nginx 站点配置文件。
- `bin/`: 存放二进制文件及可选的内置 Nginx 运行环境。

## 核心逻辑
1. **后端**：使用 Gin 提供 RESTful API，GORM 操作 SQLite。
2. **前端**：通过 Fetch API 与后端通信，实现站点管理的 CRUD 操作。
3. **Nginx 协作**：后端生成符合 Nginx 语法的配置文件，并存储在指定目录，方便 Nginx 引用。
