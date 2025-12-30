# Proton RDS SDK Go

一个高效的多数据库 RDS SDK，支持多种数据库类型并提供读写分离功能。

## 项目简介

Proton RDS SDK Go 是一个专为 RDS（关系型数据库服务）设计的 Go 语言 SDK，提供统一的数据库访问接口，支持多种数据库类型，并内置读写分离功能，帮助开发者更高效地使用 RDS 服务。

## 主要特性

### 多数据库支持
- **MySQL/MariaDB**: 支持标准的 MySQL 和 MariaDB 数据库
- **GoldenDB**: 支持 GoldenDB 分布式数据库
- **DM8**: 支持达梦数据库 DM8
- **TiDB**: 支持 TiDB 分布式数据库
- **KingBase**: 支持人大金仓数据库 KDB9

### 读写分离
- **sqlx 包**: 提供读写分离功能，自动将读操作路由到从库，写操作路由到主库
- **连接池优化**: 智能连接池管理，支持连接复用和生命周期管理
- **负载均衡**: 支持多个读节点的负载均衡

### 高级功能
- **连接池管理**: 优化的连接池配置，避免 TIME_WAIT 问题
- **类型安全**: 提供类型安全的 SQL 操作接口
- **事务支持**: 完整的事务支持，包括分布式事务
- **监控和日志**: 内置连接监控和 SQL 日志功能

## 快速开始

### 安装

```bash
go get github.com/kweaver-ai/proton-rds-sdk-go
```

### 基本使用

#### 1. 使用 Driver 包连接数据库

```go
package main

import (
    "database/sql"
    "fmt"
    "os"
    
    _ "github.com/kweaver-ai/proton-rds-sdk-go/driver"
)

func main() {
    // 设置数据库类型
    os.Setenv("DB_TYPE", "mysql") // 支持: mysql, mariadb, goldendb, dm8, tidb, kdb9
    
    // 连接数据库
    db, err := sql.Open("proton-rds", "user:password@tcp(host:port)/database")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    
    // 执行 SQL
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INT, name VARCHAR(50))")
    if err != nil {
        panic(err)
    }
    
    fmt.Println("数据库连接成功！")
}
```

#### 2. 使用 sqlx 包实现读写分离

```go
package main

import (
    "fmt"
    
    "github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
)

func main() {
    // 配置数据库连接信息
    connInfo := sqlx.DBConfig{
        User:      "username",
        Password:  "password",
        Host:      "master-host",     // 主库地址
        Port:      3306,              // 主库端口
        HostRead:  "slave-host",      // 从库地址
        PortRead:  3306,              // 从库端口
        Database:  "mydb",
        ParseTime: "true",
        Loc:       "Local",
    }
    
    // 创建数据库连接
    db, err := sqlx.NewDB(&connInfo)
    if err != nil {
        panic(err)
    }
    defer db.Close()
    
    // 写操作 - 自动路由到主库
    _, err = db.Exec("INSERT INTO users (name) VALUES (?)", "张三")
    if err != nil {
        panic(err)
    }
    
    // 读操作 - 自动路由到从库
    rows, err := db.Query("SELECT id, name FROM users")
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var id int
        var name string
        err = rows.Scan(&id, &name)
        if err != nil {
            panic(err)
        }
        fmt.Printf("ID: %d, Name: %s\n", id, name)
    }
}
```

## 环境变量配置

| 环境变量 | 说明 | 可选值 |
|---------|------|--------|
| DB_TYPE | 数据库类型 | mysql, mariadb, goldendb, dm8, tidb, kdb9, default |

## 数据库特定配置

### MySQL/MariaDB
```go
dsn := "user:password@tcp(host:port)/database?timeout=10s&charset=utf8mb4&parseTime=true"
```

### DM8 (达梦数据库)
```go
dsn := "user:password@tcp(host:port)/database?timeout=10s&client_encoding=utf-8"
```

### TiDB
```go
dsn := "user:password@tcp(host:port)/database?timeout=10s&charset=utf8mb4"
```

### KingBase (人大金仓)
```go
dsn := "user:password@tcp(host:port)/database?timeout=10s&sslmode=disable"
```

## 连接池配置建议

基于 Go 数据库连接池的最佳实践，建议采用以下配置：

```go
db.SetMaxOpenConns(100)      // 最大连接数
db.SetMaxIdleConns(100)      // 最大空闲连接数，建议与最大连接数相等
db.SetConnMaxIdleTime(120 * time.Second)  // 空闲连接最大存活时间
db.SetConnMaxLifetime(0)     // 连接最大生命周期，0 表示不限制
```

## 项目结构

```
proton-rds-sdk-go/
├── driver/           # 数据库驱动包
│   ├── common/      # 通用工具
│   ├── dmdb/        # 达梦数据库驱动
│   ├── goldendb/    # GoldenDB 驱动
│   ├── kingbase/    # 人大金仓驱动
│   ├── mysql/       # MySQL 驱动
│   └── tidb/        # TiDB 驱动
├── sqlx/            # 读写分离和连接池管理
├── example/         # 使用示例
│   ├── driver/      # 驱动使用示例
│   └── rw-split/    # 读写分离示例
├── go.mod           # Go 模块定义
└── VERSION          # 版本信息
```

## 开发计划

- [x] 多数据库支持
- [x] 读写分离功能
- [x] 连接池优化
- [ ] 分布式事务支持
- [ ] 数据库分片功能
- [ ] 性能监控面板
- [ ] SQL 审计功能

## 贡献指南

欢迎提交 Issue 和 Pull Request 来改进这个项目。

### 开发环境要求
- Go 1.24.0 或更高版本
- 相应的数据库驱动依赖

### 测试
```bash
go test ./...
```

## 许可证

本项目采用内部许可证，详情请咨询项目维护者。

## 联系方式

如有问题或建议，请联系项目维护者。

---

**版本**: 1.4.0  
**最后更新**: 2025年12月