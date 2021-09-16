一、依赖：
mysql 5.5+
redis 4.0.9+
ubuntu 14.04+  其它系统未测试

1、安装mysql，apt install mysql-server；设置密码等按照提示进行
   启动mysql-server: service mysql restart启动mysql-server
   mysql -hxx -uxxx -pxxx -Pxxx 检查安装是否成功；

2、安装redis，apt install redis-server redis-cli
   启动redis-server(使用默认配置):  redis-server&
   redis-cli 检查redis是否安装成功；


二、编译
1、下载本代码；

2、包管理需要，在go/src目录下面建立软连接以便项目之间文件可以互相引用；执行：
    cd $GOPATH/go/src;
    ln -s yourdirectory/goprofile  entry_task

3、修改mysql\redis链接参数，
   打开goprofile/mysqlconn/mysqlconn.go，修改第17行 const dbconfig string 为系统mysql登入参数；
   打开goprofile/redisconn/redisconn.go，修改第14行 const rediscfg string 为系统redis登入参数；
   进入mysql，创建数据库和表：
        create database my_test;
        CREATE TABLE IF NOT EXISTS `profile`(
            `username` VARCHAR(64) NOT NULL,
            `nickname` VARCHAR(64),
            `password` VARCHAR(64),
            `picture` varchar(127),
            PRIMARY KEY ( `username` )
        )ENGINE=InnoDB DEFAULT CHARSET=utf8;

4、编译页面端，
    cd goprofile/www;
    go build page.go

5、编译后端TcpServer，
    cd goprofile/RPCServer；
    go build serverbase.go；
   


