2021.09.09
1、学习go语言，变量、流程控制、结构体；了解一些标准库；
2、审题，http服务，tcp DAO层，http->tcp DAO之间的rpc服务；重点是rpc的实现；
3、与廖振良、张浩讨论，rpc需要自己实现；但rpc包含内容庞大，自己取舍；
4、完成go web application教学，完成登入静态页面。
明日计划：
1、完成鉴权页面、展示页面、更新nickname、profile picture页面；
2、了解go语言redis接口、mysql接口，搭建DAO层、TCP server层简单框架；
3、了解session实现，完成鉴权逻辑；
4、了解go语言http上传下载接口，初步实现prfile prcture上传。

2021.09.10
1、测试用的本机docker环境搭建，共享端口、运行展示go静态页面；
2、docker环境redis、mysql安装、启动以及测试；go 接入redis\mysql测试；
3、了解go socket编程，完成简单客户端、服务器socket交互测试case；
4、了解go语言http上传下载接口，完成profile prcture上传与展示。
明日计划：
1、完成页面组合，页面跳转；完成登入；
2、建立数据库表字段，定义好redis字段，搭建tcpServer框架；
3、着手定义httpServer和TcpServer之间的交互协议、协议格式；
4、初步实现部分httpServer->TcpServer的接口调用。

2021.09.13
1、完成了静态页面，以及页面之间的跳转；
2、定义鉴权消息，暂定使用Json作为数据交换格式，使用短链接；
3、数据库表结构、数据库封装完成；
4、学习golang文件组织、包管理;
明日计划：
1、完成鉴权接口rpc调用；
2、开始其它接口的rpc调用；
3、部分测试；

2021.09.14
1、鉴权接口进度70%，完成客户端->服务端消息序列化、反序列化，待完成客户端接受消息响应；
2、完成鉴权接口服务端接受消息->查询redis ->查询mysql基本流程；
3、学习golang文件组织、包管理；重新组织项目文件；
明日计划：
1、完成鉴权流程；
2、完成更新nickname接口；
3、开始更新头像接口；
4、测试

明日计划：
1、完成鉴权接口rpc调用；
2、开始其它接口的rpc调用；
3、部分测试；

2021.09.15
1、学习goroutine、channel、interface；
2、调通整个客户端<->服务器之间鉴权接口请求与响应；
3、完成更新nickname接口，更新profile picture接口；
明日计划：
1、完成页面端与RPCClient对接；
2、开始测试，完善部分文档；
