gorevel
=======

1. Go 语言 Revel 框架学习— [Revel中文社区](http://gorevel.cn) 源码，本站使用Revel、xorm构建。

2. 配置文件在 src/gorevel/conf 目录中，主配置app.conf，自定义配置my.conf (数据库、邮件)。

3. 默认的数据库是mysql，数据库名gorevel，表结构不需要创建，程序启动时由xorm自动创建。

4. 默认添加管理员账号admin 密码123。

###Requirements

- Go1.2+

- github.com/robfig/revel
- github.com/robfig/revel/revel
- github.com/go-sql-driver/mysql
- github.com/go-xorm/xorm
- code.google.com/p/go-uuid/uuid
- github.com/disintegration/imaging
- github.com/cbonello/revel-csrf
- github.com/qiniu/api
- github.com/garyburd/redigo/redis
- github.com/robfig/go-cache
- github.com/robfig/gomemcache/memcache

###Install

    $ git clone git://github.com/goofcc/gorevel.git
    $ cd gorevel

Linux/Unix/OS X:

    $ ./install.sh
    $ ./run.sh

Windows:

    > winstall.bat
    > wrun.bat
    
打开浏览器访问 [http://localhost:9000](http://localhost:9000)

