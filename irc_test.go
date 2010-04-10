package irc

import (
	"bufio"
	"net"
	"testing"
)

const HOST = ":6000"

// A very crude implementation of an IRC echo server for testing purposes
func init() {
	go func() {
		ls, err := net.Listen("tcp", HOST)
		if err == nil {
			for {
				client, _ := ls.Accept()
				go func() {
					in := bufio.NewReader(client)
					out := bufio.NewWriter(client)
					line, _ := in.ReadString('\n')
					_, _ = out.WriteString(line)
					_ = out.Flush()
				}()
			}
		}
	}()
}

func TestIdentity(t *testing.T) {
	msgIn := NewMessage(":johns!john@localhost NICK johns :foo")
	t.Log("Given an IRC message \"NICK johns\"")
	t.Log("and an established IRC connection echoing messages")
	irc, err := Connect(HOST)
	if err != nil {
		t.Fatal("Unable to establish connection: ", err)
	}
	t.Log("When the message is sent")
	err = irc.Write(msgIn)
	if err != nil {
		t.Fatal("Unable to send message: ", err)
	}
	t.Log("Then a message should be received")
	msgOut, err := irc.Read()
	if err != nil {
		t.Fatal("Unable to receive message: ", err)
	}
	t.Log("and the message should be the same as the message sent")
	if msgIn.Command != msgOut.Command {
		t.Errorf("But the command was \"%s\"", msgOut.Command)
	}
	if msgIn.Args[0] != msgOut.Args[0] {
		t.Errorf("But the argument was \"%s\"", msgOut.Args[0])
	}
}
