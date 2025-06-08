package tcp

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)

func HandleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Client disconnected:", conn.RemoteAddr())
            break
        }

        msg = strings.TrimSpace(msg)
        fmt.Printf("[%v] %s\n", conn.RemoteAddr(), msg)

        // Echo message back (replace with custom logic)
        conn.Write([]byte("Echo: " + msg + "\n"))
    }
}
