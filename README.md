# gorevel

1. Go 语言 Revel 框架学习— [Revel中文社区](http://gorevel.cn) 源码，本站使用 Revel、xorm 构建。

1. 配置文件在 src/gorevel/conf 目录中，主配置 app.conf，自定义配置 my.conf (数据库、邮件)。

1. 默认的数据库是 mysql，数据库名 gorevel，表结构不需要创建，程序启动时由 xorm 自动创建。

1. 默认添加管理员账号 admin 密码 123。

## Install

`
    $ go get github.com/goofcc/gorevel
    $ revel run gorevel
`

打开浏览器访问 [http://localhost:9000](http://localhost:9000)
