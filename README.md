# gorevel

Go 语言 Revel 框架学习, [Revel 中文社区](http://gorevel.cn) 源码。

## Requirements

* Revel v0.17.1
* Go v1.6+

## Install

``` bash
  cd $GOPATH/src
  git clone https://github.com/goofcc/gorevel.git
```

* 配置文件在 src/gorevel/conf 目录中，app.conf，my.conf。
* 数据库是 mysql，数据库名 gorevel，程序启动时自动创建表结构。
* 默认添加管理员账号 admin 密码 123。

## Run

``` bash
  revel run gorevel
```

* 打开浏览器访问 <http://localhost:9000>
