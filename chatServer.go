package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	cons map[*websocket.Conn]bool
}

func NewSever() *Server {

	return &Server{cons: make(map[*websocket.Conn]bool)}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	s.cons[ws] = true
	s.readLoop(ws)
}

func (s *Server) handleBookFeed(ws *websocket.Conn){
	feed := fmt.Sprintf("Here are you books -> %d\n",time.Now().UnixNano())
	for{
		ws.Write([]byte(feed))
		time.Sleep(time.Second * 2)
	}
}
func (s *Server) readLoop(ws *websocket.Conn) {
	buff := make([]byte, 1024)
	for {
		n, err := ws.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("An error has occured", err)
			continue
		}
		msg := buff[:n]
		fmt.Println(string(msg))
		ws.Write([]byte("Thank you for the message"))
		s.broadCast(msg)
	}

}

func (s *Server) broadCast(buff []byte) {

	for ws := range s.cons {
		if _, err := ws.Write(buff); err != nil {
			fmt.Println("An error has occured", err)
		}
	}

}
func main() {
	server := NewSever()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/bookFeed",websocket.Handler(server.handleBookFeed))
	http.ListenAndServe(":3000", nil)

}
