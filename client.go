package main

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

// clientはチャットを行っている1人のユーザーを表す。
type client struct {
	socket   *websocket.Conn        // このクライアントのためのWebSocket
	send     chan *message          // メッセージが送られるチャネル
	room     *room                  // このクライアントが参加しているチャットルーム
	uuid     uuid.UUID              // UUID(version 4)
	userData map[string]interface{} // ユーザー情報を保持
}

// クライアントがWebSocketからReadMessageを使ってデータを読み込むために使用する。
func (c *client) read() {
	log.Debugf("client.read: 読み取り開始します。 client=%s", c.uuid)
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err != nil {
			log.Error("client.read: ", err)
			break // エラーが発生した場合、ループから脱出してWebSocketを閉じる。
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		if avatarURL, ok := c.userData["avatar_url"]; ok {
			msg.AvatarURL = avatarURL.(string)
		}
		c.room.forward <- msg // 受け取ったメッセージはすぐにroomのforwardチャネルに送られる。
	}
	c.socket.Close() // WebSocketを閉じる。
}

// 継続的にsendチャネルからメッセージを受け取り、WebSocketのWriteMessageメソッドを使って内容を書き出す。
func (c *client) write() {
	log.Debugf("client.write: 書き込み開始します。 client=%s", c.uuid)
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			log.Error("client.write: ", err)
			break // WebSocketへの書き込みが失敗した場合、ループから脱出してWebSocketを閉じる。
		}
	}
	c.socket.Close()
}
