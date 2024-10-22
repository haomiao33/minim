const WebSocket = require("ws");

// 设置 WebSocket 服务器地址
const serverUrl = "ws://0.0.0.0:3000"; // 替换为你的 WebSocket 服务器地址

// 定义一个函数来创建 WebSocket 客户端并发送消息
async function createWebSocketClientTob(clientId) {
    return new Promise((resolve, reject) => {
        const ws = new WebSocket(serverUrl);

        // 当连接打开时
        ws.on("open", () => {
            console.log(`Client ${clientId} connected to server`);

            // 发送登录消息
            const loginMessage = {
                type: "login",
                data: { userId: clientId } // 使用合适的 userId
            };
            ws.send(JSON.stringify(loginMessage));
            // console.log(`Client ${clientId} sent login message:`, loginMessage);

            // 发送 IM 消息
            const imMessage = {
                type: "msg",
                data: {
                    msgId: new Date().getTime() + `-123-456-` + '0-' + Math.floor(Math.random() * 1000000) + `${clientId}`,
                    chatType: 0,     //0=单聊；1=一般群； 2=机器人
                    msgType: 1,           // 消息类型； 1=文本；2=图片；3=视频；4=文件；5=通话
                    fromId: clientId,    // 发送者
                    toId: 456,      // 接收者
                    message: "Hello!-456-to-123-" + clientId,   // 消息内容
                    ts: Date.now()
                }
            };
            ws.send(JSON.stringify(imMessage));
            // console.log(`Client ${clientId} sent IM message:`, imMessage);

            // 接收消息
            ws.on("message", (data) => {
                // console.log(`Client ${clientId} received message from server:`, JSON.parse(data));
            });

            // 处理关闭
            ws.on("close", () => {
                console.error(`Client ${clientId} connection closed`);
                resolve(); // 关闭时 resolve promise
            });

            // 处理错误
            ws.on("error", (error) => {
                console.error(`Client ${clientId} WebSocket error:`, error);
                reject(error); // 发生错误时 reject promise
            });
        });

        // 处理连接错误
        ws.on("error", (error) => {
            console.error('---------------------errir', error)
            reject(error); // 连接错误时 reject promise
        });
    });
}



// 定义一个函数来创建 WebSocket 客户端并发送消息
async function createWebSocketClientToa(clientId) {
    return new Promise((resolve, reject) => {
        const ws = new WebSocket(serverUrl);

        // 当连接打开时
        ws.on("open", () => {
            console.log(`Client ${clientId} connected to server`);

            // 发送登录消息
            const loginMessage = {
                type: "login",
                data: { userId: clientId } // 使用合适的 userId
            };
            ws.send(JSON.stringify(loginMessage));
            // console.log(`Client ${clientId} sent login message:`, loginMessage);


            // console.log(`Client ${clientId} sent IM message:`, imMessage);

            // 接收消息
            ws.on("message", (data) => {
                const { type, code, msg } = JSON.parse(data);
                if (type == 'loginAck') {
                    console.log('loginAck', msg)
                    // 发送 IM 消息
                    const imMessage = {
                        type: "msg",
                        data: {
                            msgId: new Date().getTime() + `-456-123-` + '0-' + Math.floor(Math.random() * 1000000) + `${clientId}`,
                            chatType: 0,     //0=单聊；1=一般群； 2=机器人
                            msgType: 1,           // 消息类型； 1=文本；2=图片；3=视频；4=文件；5=通话
                            fromId: clientId,    // 发送者
                            toId: 50001,      // 接收者
                            message: "Hello!-456-to-123-" + clientId,   // 消息内容
                            ts: Date.now()
                        }
                    };
                    ws.send(JSON.stringify(imMessage));
                }

                // console.log(`Client ${clientId} received message from server:`, JSON.parse(data));
            });

            // 处理关闭
            ws.on("close", () => {
                console.error(`Client ${clientId} connection closed`);
                resolve(); // 关闭时 resolve promise
            });

            // 处理错误
            ws.on("error", (error) => {
                console.error(`Client ${clientId} WebSocket error:`, error);
                reject(error); // 发生错误时 reject promise
            });
        });

        // 处理连接错误
        ws.on("error", (error) => {
            console.error('---------------------errir', error)
            reject(error); // 连接错误时 reject promise
        });

        return new Promise(resolve => setTimeout(resolve, 10000));
    });
}


// 创建并发客户端的主函数
async function main() {
    const clientCount = 2000; // 客户端数量
    const promises = [];

    for (let i = 1; i <= 3000; i++) {
        promises.push(createWebSocketClientToa(i)); // 为每个客户端创建 promise
    }



    // 等待所有客户端完成
    try {
        await Promise.all(promises);
        console.log("All clients have completed their actions.");
    } catch (error) {
        console.error("An error occurred while handling WebSocket clients:", error);
    }

    //wait
    await new Promise(resolve => setTimeout(resolve, 130000));
}

// 启动客户端
main();
