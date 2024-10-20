# gerty-goblog

一个go语言实现的blog应用，支持多用户管理,多数据库支持

### TODOS

- [x]  基本文章的增删改查
- [x]  多用户支持
- [ ]  后台管理
- [ ]  插件工具支持


### 数据库
创建数据库
```sql
CREATE DATABASE goblog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

数据表
- articles
- users
- categories
- comments

数据库支持mysql和postgresql，使用`go run cmd/init_config.go`可以生成`.env`文件的模板，根据实际情况填写数据库信息。数据库交互使用gorm，无需手动创建数据表，只需要在`models/`中定义数据表对应的类即可。
**PS:gorm默认情况下会将表名转换成复数形式，比如自定义的model是article，那么在migrate之后的数据库表就是articles，如果需要自定义表名需要在自定义model中的struct下实现`TabNAme()`方法**

### 项目架构
项目采用MVC架构，`controller`层负责接收处理请求并做出对应的响应，`model`层构建数据表以及和数据库交互，`view`层负责页面的渲染。
同时采用自定义中间件来做访问权限控制。


### 使用到的库

#### Viper
**viper**适用于Go程序的配置解决方案
- 设置默认值
- 读取环境变量
- 从配置文件（JSON,YAML,TOML）读取
- 从远程配置系统重读取
使用`.env`文件来存储配置，由于viper支持默认值，所以可以在不设定对应变量的情况下直接使用。

#### gorm
go专用的orm工具。

#### air
适用于go项目的热加载工具。
