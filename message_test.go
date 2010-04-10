package irc

import (
	"fmt"
	"strings"
	"testing"
)

func TestMessageCommand(t *testing.T) {
	cmd := "NICK"
	raw := cmd
	args := make([]string, 0)
	msg := NewMessage(raw)
	t.Logf("Given the IRC message \"%s\"", raw)
	checkCommand(t, msg, cmd)
	checkArgs(t, msg, args)
}

func TestMessageCommandWithArgument(t *testing.T) {
	cmd := "NICK"
	args := make([]string, 1)
	args[0] = "johns"
	raw := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	msg := NewMessage(raw)
	t.Logf("Given the IRC message \"%s\"", raw)
	checkCommand(t, msg, cmd)
	checkArgs(t, msg, args)
}

func TestMessageCommandWithArguments(t *testing.T) {
	raw := "SPOOF one two"
	args := make([]string, 2)
	args[0] = "one"
	args[1] = "two"
	msg := NewMessage(raw)
	t.Logf("Given the IRC message \"%s\"", raw)
	checkCommand(t, msg, "SPOOF")
	checkArgs(t, msg, args)
}

func TestMessagecommandWithTrailingArguments(t *testing.T) {
	raw := "PRIVMSG johns :How are your tests going?"
	args := make([]string, 2)
	args[0] = "johns"
	args[1] = ":How are your tests going?"
	msg := NewMessage(raw)
	t.Logf("Given the IRC message \"%s\"", raw)
	checkCommand(t, msg, "PRIVMSG")
	checkArgs(t, msg, args)
}

func TestMessageCommandWithServerPrefix(t *testing.T) {
	server := "google.com"
	cmd := "SPOOF"
	args := make([]string, 0)
	raw := fmt.Sprintf(":%s %s", server, cmd)
	msg := NewMessage(raw)
	t.Logf("Given the IRC message \"%s\"", raw)
	checkCommand(t, msg, cmd)
	checkServer(t, msg, server)
	checkArgs(t, msg, args)
}

func TestMessageCommandWithNickPrefix(t *testing.T) {
	nick := "johns"
	cmd := "SPOOF"
	args := make([]string, 0)
	raw := fmt.Sprintf(":%s %s", nick, cmd)
	msg := NewMessage(raw)
	t.Logf("Given the IRC message \"%s\"", raw)
	checkCommand(t, msg, cmd)
	checkNick(t, msg, nick, "", "")
	checkArgs(t, msg, args)
}

func TestMessageCommandWithNickAndUserPrefix(t *testing.T) {
	args := make([]string, 0)
	raw := ":johns!John SPOOF"
	msg := NewMessage(raw)
	t.Logf("Given the IRC message \"%s\"", raw)
	checkCommand(t, msg, "SPOOF")
	checkNick(t, msg, "johns", "John", "")
	checkArgs(t, msg, args)
}

func TestMessageCommandWithNickUserAndHostPrefix(t *testing.T) {
	args := make([]string, 0)
	raw := ":johns!John@localhost SPOOF"
	msg := NewMessage(raw)
	t.Logf("Given the IRC message \"%s\"", raw)
	checkCommand(t, msg, "SPOOF")
	checkNick(t, msg, "johns", "John", "localhost")
	checkArgs(t, msg, args)
}

func checkCommand(t *testing.T, msg *Message, cmd string) {
	t.Logf("Then the command should be \"%s\"", cmd)
	if msg.Command != cmd {
		t.Errorf("But the command was \"%s\"", msg.Command)
	}
}

func checkServer(t *testing.T, msg *Message, server string) {
	t.Logf("And the server should be \"%s\"", server)
	if msg.Server != server {
		t.Errorf("But the server was \"%s\"", msg.Server)
	}
}

func makeNick(nick, user, host string) (full string) {
	if nick == "" {
		return
	}
	full = ":" + nick
	if user != "" {
		full += "!" + user
	}
	if host != "" {
		full += "@" + host
	}
	return
}

func checkNick(t *testing.T, msg *Message, nick, user, host string) {
	got := makeNick(msg.Nick, msg.User, msg.Host)
	expected := makeNick(nick, user, host)
	t.Logf("And the prefix should be \"%s\"", expected)
	if got != expected {
		t.Errorf("But the prefix was \"%s\"", got)
	}
}

func checkArgs(t *testing.T, msg *Message, args []string) {
	var nArgs string
	switch n := len(args); n {
	case 0:
		nArgs = "no arguments"
	case 1:
		nArgs = "a single argument"
	default:
		nArgs = fmt.Sprintf("%d arguments", n)
	}
	t.Logf("And there should be %s", nArgs)
	if n := len(msg.Args); n != len(args) {
		t.Fatalf("But there were %d arguments", n)
	}
	if len(args) == 1 {
		t.Logf("And the argument should be \"%s\"", args[0])
		if args[0] != msg.Args[0] {
			t.Errorf("But the argument was \"%s\"", msg.Args[0])
		}
	} else {
		for i, arg := range args {
			t.Logf("And argument %d should be \"%s\"", i+1, arg)
			if arg != msg.Args[i] {
				t.Errorf("But argument %d was \"%s\"", i+1, msg.Args[i])
			}
		}
	}
}
