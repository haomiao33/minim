# mini-im

## 1、说明：

### 1.1、项目介绍：
  简单的go写的im服务，流程简单清晰，采用golang编写，分为login、msg、msg-push、online等服务。可以自己扩充其他协议和服务。
    
  作者：目前打算golang编写im，支持单聊、群聊、推送； 然后客户端目前只打算做个flutter im chat 版本的就行。方便大家集成。

  文档放在doc里面了

  websocket底层才用[gnet](https://github.com/panjf2000/gnet) ，后续可以很轻松扩展到tcp和udp
  
  sql里面是数据库，自己创建一个就行

  测试客户端：目前是写到testclient目录里面的，nodejs的

  目前有如下服务（每个服务都可以多实例，方便用户量上来扩展和分布式，目前服务注册到consul里面的）：
    
    login: 用户登录服务，用户和im的长连接都连接到login服务，目前是用的websocket，后续可以扩展到tcp、udp等。不同服务交互使用的grpc，这里面没有写用户认证什么的，大家可以根据自己的业务需求来完成。
    
    msg: 消息服务，处理消息接收
        
    msg-push: 消息推送服务，单聊消息推送
        
    online: 在线状态服务，用户在线状态放在这里，内部使用redis存放。login服务

## 1、2 项目部署：
查看 [部署方法](doc/部署.md) 文件

## 2、登陆
<img width="740" alt="image" src="https://github.com/user-attachments/assets/bd8024fa-f838-43ac-b4be-ee0066ed5a5e">



## 3、单聊消息：流程图

### 3.1、单聊流程图
<img width="1059" alt="image" src="https://github.com/user-attachments/assets/779dc2eb-b814-4131-99c2-935e81601fbf">




