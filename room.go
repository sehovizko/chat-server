package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"github.com/stretchr/objx"
)

// Constant variables
const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// チャットルーム構造体
type room struct {
	forward chan *message    // 他のクライアントに転送するためのメッセージを保持するチャネル。
	join    chan *client     // チャットルームに参加しようとしているクライアントのためのチャネル。
	leave   chan *client     // チャットルームから退室しようとしているクライアントのためのチャネル。
	clients map[*client]bool // 在室しているすべてのクライアントが保持される。
}

// チャットルーム生成関数
func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

// チャットルームの開始
func (r *room) run() {

	log.Info("room.run: チャットルームを開始します。")
	for {
		log.Debug("room.run: ループ処理始端")
		// チャネルに送信された値に応じて処理を分岐させる。
		select {
		case client := <-r.join: // 参加
			r.clients[client] = true
			log.Infof("room.run: 参加される方がいます。client=%s", client.u)
		case client := <-r.leave: // 退室
			delete(r.clients, client)
			close(client.send)
			log.Infof("room.run: 退出される方がいます。 client=%s", client.u)
		case msg := <-r.forward: // すべてのクライアントにメッセージを転送
			log.Info("room.run: メッセージを受信しました。: ", msg.Message)
			for client := range r.clients {
				select {
				case client.send <- msg: // メッセージを送信
					log.Infof("room.run: --送信に成功。 client=%s", client.u)
				default: // 送信に失敗
					delete(r.clients, client)
					close(client.send)
					log.Warnf("room.run: --送信に失敗。 client=%s", client.u)
				}
			}
		}
		log.Debug("room.run: ループ処理の終端")
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// HTTP接続をアップグレードする
	upgrader := websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Error("room.ServeHTTP: ", err)
		return
	}

	log.Debug("room.ServeHTTP: HTTP接続しました!")

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Error("room.ServeHTTP: ", err)
		return
	}
	// クライアントを生成して現在のチャットルームのjoinチャネルに渡す。
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		u:        uuid.Must(uuid.NewV4()),
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client

	// defer文で、クライアントの終了時に退室の処理を行うように指定する。(ユーザーがいなくなった際のクリーンアップ)
	defer func() { r.leave <- client }()

	// goroutineとしてwriteメソッドを実行する。
	go client.write()
	// メインスレッドではreadメソッドを呼び出し接続は保持。
	client.read() // 終了を指示されるまで他の処理はブロックされる。
}
