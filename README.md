# gerty-goblog

一个go语言实现的blog应用，支持多用户管理,多数据库支持

### TODOS

- [x]  基本文章的增删改查
- [ ]  多用户支持
- [ ]  后台管理
- [ ]  插件工具支持


### 数据库
创建数据库
```sql
CREATE DATABASE goblog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

创建articles表
```sql
CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
);
```
