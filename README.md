gorevel
=======

1. Go语言Revel框架学习—[Revel中文社区](http://gorevel.cn)源码，本站使用Revel、Qbs构建。

2. 配置文件在 src/gorevel/conf 目录中，主配置app.conf，自定义配置my.conf (数据库、邮件)。

3. 默认的数据库是mysql，数据库名gorevel，表结构不需要创建，程序启动时由Qbs自动创建。

###Requirements

- Go1.1+
- github.com/robfig/revel
- github.com/robfig/revel/revel
- github.com/coocood/qbs
- github.com/coocood/mysql
- github.com/disintegration/imaging
- code.google.com/p/go-uuid/uuid
- github.com/cbonello/revel-csrf

###Install

    $ git clone git://github.com/goofcc/gorevel.git
    $ cd gorevel

Linux/Unix/OS X:

    $ ./install.sh
    $ ./run.sh

Windows:

    > install.bat
    > run.bat
    
打开浏览器访问 [http://localhost:9000](http://localhost:9000)

