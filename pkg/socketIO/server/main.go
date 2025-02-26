package main

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"net/http"
)

/*
官方示例代码
*/

func main() {
	// 创建新的 Socket.IO 服务器实例
	server := socketio.NewServer(nil)

	// 处理客户端连接事件
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")                  // 设置连接上下文
		fmt.Println("connected:", s.ID()) // 打印连接成功的客户端 ID
		return nil
	})

	// 处理客户端发送的 “notice” 事件
	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg) //打印客户端发送的消息
		s.Emit("reply", "have"+msg) //发送恢复消息给客户端
	})

	// 处理客户端发送的 “msg” 事件
	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)   //设置连接上下文为收到的消息
		return "recv" + msg //返回接收到的消息
	})

	// 处理客户端发送的 “bye” 事件
	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string) // 获取连接上下文中的消息
		s.Emit("bye", last)          // 发送再见消息给客户端
		s.Close()
		return last // 返回最后收到的消息
	})

	// 处理连接错误事件
	server.OnError("/", func(s socketio.Conn, err error) {
		fmt.Println("meet error:", err)
	})

	// 处理客户端断开连接事件
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	// 使用 redis 广播适配器的示例
	_, err := server.Adapter(&socketio.RedisAdapterOptions{ // 配置服务器的适配器为 Redis
		Addr:   "127.0.0.1:6379",
		Prefix: "socket.io",
		DB:     1,
	})
	if err != nil {
		log.Fatal("error", err)
	}

	// 启动 Socket.IO 服务器
	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset"))) // 指定静态文件目录
	log.Println("serving at localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil)) // 启动 HTTP 服务器并监听 8000 端口
}
