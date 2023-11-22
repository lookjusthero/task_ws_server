package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type Singelton struct {
	Connected map[string]struct{}
	numbers   map[string]struct{}

	MUConnected sync.Mutex
	muNumber    sync.Mutex

	bigIntMax *big.Int
}

func (s *Singelton) GetBigInt() *big.Int {
	s.muNumber.Lock()

	var n *big.Int

	for {
		n, _ = rand.Int(rand.Reader, s.bigIntMax)

		if _, ok := s.numbers[n.Text(62)]; !ok {
			s.numbers[n.Text(62)] = struct{}{}
			break
		}
	}

	s.muNumber.Unlock()

	return n
}

func main() {
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))

	singelton := &Singelton{
		Connected: make(map[string]struct{}, 1000),
		numbers:   make(map[string]struct{}, 1000),
		bigIntMax: max,
	}

	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) { BigIntSender(ws, singelton) }))

	http.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(`

		var webSocket = new WebSocket("ws://localhost:8080/ws", "hi"); 
		webSocket.onmessage = ev => {console.log(ev.data)}; 

		// timeout для установки соединения 
		setTimeout(() => {webSocket.send("test")}, 300);

		webSocket.send("test")
	`))
	}))

	s := &http.Server{
		Addr: "0.0.0.0:8080",
	}

	s.ListenAndServe()
}

func BigIntSender(ws *websocket.Conn, s *Singelton) {
	defer ws.Close()

	clientIP, _, err := net.SplitHostPort(ws.Request().RemoteAddr)
	if err != nil {
		fmt.Printf("Неудалось распарить ip адрес %s", clientIP)
		return
	}

	fmt.Printf("Подключение клиента %s\n", clientIP)

	s.MUConnected.Lock()

	if _, ok := s.Connected[clientIP]; ok {
		s.MUConnected.Unlock()

		fmt.Printf("Разрываем соединение для %s, у клиента уже есть активное соединение\n", clientIP)

		return
	}

	s.Connected[clientIP] = struct{}{}
	s.MUConnected.Unlock()

	var message string

	defer func() {
		s.MUConnected.Lock()
		delete(s.Connected, clientIP)
		s.MUConnected.Unlock()

		fmt.Printf("Разорвано основное соедние у клиента %s\n", clientIP)
	}()

	for {
		if err := websocket.Message.Receive(ws, &message); err != nil {
			return
		}

		msg, _ := json.Marshal(s.GetBigInt())

		ws.Write(msg)
	}
}
