package websocket

import "nhooyr.io/websocket"

type SocketConn struct {
	Conn *websocket.Conn
	Id   string

	CloseFunc func()
}

//func (s *SocketConn) Close() error {
//	if s.CloseFunc != nil {
//		s.CloseFunc()
//	}
//	return s.Conn.Close()
//}
