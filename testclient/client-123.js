const WebSocket = require("ws");

// 连接到 WebSocket 服务器
const ws = new WebSocket("ws://0.0.0.0:3000"); // 替换为你的 WebSocket 服务器地址

let sequence = 0
let localMsg = []

// 当连接打开时
ws.on("open", () => {
    console.log("Connected to server");

    // // 发送登录消息
    const loginMessage = {
        type: "login",
        data: { userId: 123 } // 使用合适的 userId
    };
    ws.send(JSON.stringify(loginMessage));
    console.log("Sent login message:", loginMessage);

    setInterval(() => {
        ws.send(JSON.stringify({ type: "heartbeat", data: {} })); // 客户端发送心跳消息
    }, 10000); // 每5秒发送一次心跳

});

function sendMsg(){
     // // 模拟发送 IM 消息
    // //msgId,from, to, message, type,ts 
    const imMessage = {
            msgId: new Date().getTime()+"-123-456-"+'0-'+Math.floor(Math.random() * 1000000) ,
            chatType:0,     //0=单聊；1=一般群； 2=机器人
            msgType: 1,           // 消息类型； 1=文本；2=图片；3=视频；4=文件；5=通话
            fromId: 123,    // 发送者
            toId: 456,      // 接收者
            message: "Hello!",   // 消息内容
            ts: Date.now()
        }
    fetch('http://0.0.0.0:3100/api/v1/msg/send', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(imMessage)
    })
    .then(response => response.json())
    .then(data => {
        // 处理响应数据
        console.log('Response:', data);
    })
}

// 接收消息
ws.on("message", (data) => {
    const ret = JSON.parse(data);
    if(ret.type == 'loginAck'){ 
        //消息发送成功
        console.log("------  loginAck success --- ",ret);

        sendMsg()
        console.log("------  send msg success --- ",ret);
    }
    else if(ret.type == 'msgSyncNotify'){
        //有新消息
        console.log("------ new msg sync notify --- ",ret.data);

      
    }else{
        console.log("------ other msg --- ",ret);
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
