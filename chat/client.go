package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// client struct represents a single chatting user
type client struct {
	socket  *websocket.Conn   //websocket for this particular client
	send    chan *message     //channel which used to send message from this user
	room    *room             //a place this user use to chat in
	userData map[string]interface{}  //menyimpan data user
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		//read message
		//_, msg, err := c.socket.ReadMessage()
		//if  err != nil {
		//	return
		//}
		////send message to forward channel in room object
		//c.room.forward <- msg

		var msg *message
		//mmbaca objek json dari interface message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now().Format("Mon Jan _2 2006 15:04:05 ")
		msg.Name = c.userData["name"].(string)
		//cek apakah avatar user ada atau tidak
		//msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
		if avatarUrl, ok := c.userData["avatar_url"]; ok {
			msg.AvatarURL = avatarUrl.(string)
		}

		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		//send message to websocket
		//if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
		//	return
		//}
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}