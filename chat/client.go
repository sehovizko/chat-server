package main

import (
	"github.com/gorilla/websocket"
)

// clientはチャットを行っている1人のユーザーを表す。
type client struct {
	socket *websocket.Conn // このクライアントのためのWebSocket
	send   chan []byte     // メッセージが送られるチャネル
	room   *room           // このクライアントが参加しているチャットルーム
}

// クライアントがWebSocketからReadMessageを使ってデータを読み込むために使用する。
func (c *client) read() {
	log.Infof("client.read: 読み取り開始します。 client=%p", c)
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg // 受け取ったメッセージはすぐにroomのforwardチャネルに送られる。
		} else {
			break // エラーが発生した場合、ループから脱出してWebSocketを閉じる。
		}
	}
	c.socket.Close() // WebSocketを閉じる。
}

// 継続的にsendチャネルからメッセージを受け取り、WebSocketのWriteMessageメソッドを使って内容を書き出す。
func (c *client) write() {
	log.Infof("client.write: 書き込み開始します。 client=%p", c)
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break // WebSocketへの書き込みが失敗した場合、ループから脱出してWebSocketを閉じる。
		}
	}
	c.socket.Close()
}
