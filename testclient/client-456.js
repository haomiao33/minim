const WebSocket = require("ws");

// 连接到 WebSocket 服务器
const ws = new WebSocket("ws://0.0.0.0:3000"); // 替换为你的 WebSocket 服务器地址

// 当连接打开时
ws.on("open", () => {
    console.log("Connected to server");

    // // 发送登录消息
    const loginMessage = {
        type: "login",
        data: { userId: 456 } // 使用合适的 userId
    };
    ws.send(JSON.stringify(loginMessage));
    console.log("Sent login message:", loginMessage);

    setInterval(() => {
        ws.send(JSON.stringify({ type: "heartbeat", data: {} })); // 客户端发送心跳消息
    }, 10000); // 每5秒发送一次心跳

    // // 模拟发送 IM 消息
    // //msgId,from, to, message, type,ts 
    const imMessage = {
        type: "msg",
        data: {
            msgId: new Date().getTime()+"-456-to-123-"+'0-'+Math.floor(Math.random() * 1000000) ,
            chatType:0,     //0=单聊；1=一般群； 2=机器人
            msgType: 1,           // 消息类型； 1=文本；2=图片；3=视频；4=文件；5=通话
            fromId: 456,    // 发送者
            toId: 123,      // 接收者
            message: "Hello!",   // 消息内容
            ts: Date.now()
        }
    };
    ws.send(JSON.stringify(imMessage));
    console.log("Sent IM message:", imMessage);
});

// 接收消息
ws.on("message", (data) => {
    const ret = JSON.parse(data);
    console.log("Message from server:",ret );
    if(ret.code = 'msg'){
        //新消息
        // ws.send(JSON.stringify({ type: "msgack", data:{
        //     id: ret.data.id,
        //     //0=已发送, 1=已送达, 2=已读, 3=已撤回
        //     status: 1
        // } }));
    }
});


// 处理关闭
ws.on("close", () => {
    console.log("Connection closed");
});

// 处理错误
ws.on("error", (error) => {
    console.error("WebSocket error:", error);
});
