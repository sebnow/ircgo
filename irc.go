// vi: noet ft=go:

/* Copyright (c) 2010 Sebastian Nowicki <sebnow@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package irc

import (
	"os"
	"fmt"
	"strings"
	"net"
	"io"
	"bufio"
)

// Abstracts an IRC message. All fields are optional aside from Command.
type Message struct {
	Command string
	Server  string
	Nick    string
	User    string
	Host    string
	Args    []string
	raw     string
}

// Returns a new Message by parsing the raw IRC message.
// An IRC message has the following format:
//
//     [ prefix ' ' ] <Command> [ ' ' <Args> ]
//     prefix :: = ':' <Server> | <Nick> [ '!' <User> ] [ '@' <Host> ]
//
// These fields directly correspond to the fields in Message. These
// fields are separated with no furthur processing. If trailing
// arguments are present (":foo bar" at the end of the message), these
// will be stored verbatim as the last argument in Args.
func NewMessage(raw string) (*Message) {
	var prefix, command, server, nick, user, host, args string
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, ":") {
		pieces := strings.Split(raw[1:], " ", 3)
		if len(pieces) < 2 {
			return nil
		}
		prefix, command = pieces[0], pieces[1]
		if len(pieces) > 2 {
			args = pieces[2]
		}
	} else {
		pieces := strings.Split(raw, " ", 2)
		command = pieces[0]
		if len(pieces) > 1 {
			args = pieces[1]
		}
	}
	// Prefix has a format of:
	// ':' <servername> | <nick> [ '!' <user> ] [ '@' <host> ]
	if prefix != "" {
		// Check for "!user[@host]"
		if strings.Index(prefix, "!") > 0 {
			var userhost string
			pieces := strings.Split(prefix, "!", 2)
			nick, userhost = pieces[0], pieces[1]
			// Check for "@host"
			if strings.Index(userhost, "@") > 0 {
				pieces := strings.Split(userhost, "@", 2)
				user, host = pieces[0], pieces[1]
			} else {
				user = userhost
			}
		} else if strings.Index(prefix, "@") > 0 {
			// Check for "@host"
			pieces := strings.Split(prefix, "@", 2)
			nick, host = pieces[0], pieces[1]
		} else {
			// We only have server or nick
			nick = prefix
			server = prefix
		}
	}

	var splitArgs []string
	if len(args) > 0 {
		trailingPos := strings.Index(args, ":")
		if trailingPos >= 0 {
			nspaces := strings.Count(args[0:trailingPos], " ")
			splitArgs = strings.Split(args, " ", nspaces+1)
		} else {
			splitArgs = strings.Fields(args)
		}
	}

	return &Message{
		raw:     raw,
		Command: command,
		Server:  server,
		Nick:    nick,
		User:    user,
		Host:    host,
		Args:    splitArgs,
	}
}

func (message *Message) String() (raw string) {
	if message.Server != "" {
		raw += ":" + message.Server + " "
	} else if message.Nick != "" {
		raw += ":" + message.Nick
		if message.User != "" {
			raw += "!" + message.User
		}
		if message.Host != "" {
			raw += "@" + message.Host
		}
		raw += " "
	}
	raw += message.Command
	if len(message.Args) > 0 {
		raw += " " + strings.Join(message.Args, " ")
	}
	return
}

// Writer implements writing a Message
type Writer interface {
	// Writes a Message
	Write(*Message) os.Error
}

// Reader implements reading a Message
type Reader interface {
	// Reads and returns a Message
	Read() (*Message, os.Error)
}

// ReadWriter implements reading and writing a Message
type ReadWriter interface {
	Writer
	Reader
}

type Conn struct {
	stream *bufio.ReadWriter
	rwc    io.ReadWriteCloser
}

func newConn(rwc io.ReadWriteCloser) *Conn {
	return &Conn{
		stream: bufio.NewReadWriter(
			bufio.NewReader(rwc),
			bufio.NewWriter(rwc)),
	}
}

// Connect to an IRC server at the specified address.
// Returns a ReadWriter or an error.
func Connect(address string) (conn *Conn, err os.Error) {
	sock, err := net.Dial("tcp", "", address)
	if err == nil {
		conn = newConn(sock)
	}
	return
}

func (conn *Conn) Write(message *Message) (err os.Error) {
	// IRC messages consist of a prefix, command and params
	raw := fmt.Sprintf("%s\r\n", message)
	_, err = conn.stream.WriteString(raw)
	_ = conn.stream.Flush()
	return
}

func (conn *Conn) Read() (message *Message, err os.Error) {
	var raw string
	raw, err = conn.stream.ReadString('\n')
	if err == nil {
		message = NewMessage(raw)
	}
	return
}
