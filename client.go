package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"os"
)

var address = "ws://10.64.221.117"

func main() {
	token := connectToPayserver()
	port := enterMatchmaker(token)
	outcome := enterGame(token, port)

	fmt.Println("Test concluded, game outcome: " + outcome)
}

func connectToPayserver() string {
	fmt.Println("Connecting to payserver")
	conn, _, err := websocket.DefaultDialer.Dial(address+":7000/ws", nil)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get Wallet
	message := map[string]string{}
	conn.ReadJSON(&message)
	fmt.Println(message)

	// Get Token
	message = map[string]string{}
	conn.ReadJSON(&message)
	fmt.Println(message)
	token := message["token"]

	// Wait for status: paid
	message = map[string]string{}
	conn.ReadJSON(&message)
	status := message["status"]

	if status != "paid" {
		fmt.Println("FAIL")
		os.Exit(1)
		return ""
	}
	return token
}

func enterMatchmaker(token string) string {
	fmt.Println("Connecting to matchmaker")
	conn, _, err := websocket.DefaultDialer.Dial(address+":8000/ws", nil)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn.WriteJSON(map[string]string{
		"Token": token,
	})

	for {
		fmt.Println("Waiting for matchmaker port/status")
		// Status or Port
		message := map[string]string{}
		conn.ReadJSON(&message)
		if _, ok := message["Status"]; ok {
			fmt.Println(message)
		}

		if _, ok := message["Port"]; ok {
			fmt.Println(message)
			return message["Port"]
		}
	}
}

func enterGame(token string, port string) string {
	fmt.Println("Connecting to gameserver")
	conn, _, err := websocket.DefaultDialer.Dial(address+":"+port+"/ws", nil)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn.WriteJSON(map[string]string{
		"Token": token,
	})

	frames := 0

	for {
		message := map[string]string{}
		conn.ReadJSON(&message)
		fmt.Println(message)
		if _, ok := message["Leaderboard"]; ok {
			conn.WriteJSON(map[string]interface{}{
				"CurrentZoomLevel": rand.Intn(24),
				"CurrentSprint":    rand.Intn(1) == 0,
				"CurrentDirection": rand.Intn(4),
			})
			frames++
			fmt.Println(frames)
		}
	}

	return "unknown"
}
