package messengerserviceapi

import (
	messengermanager "MessengerService/mesermanager"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/zishang520/socket.io/socket"
)

// hasUserMiddleware checks if user is loged in
func hasUserMiddleware[T any](conn *socket.Socket, next func(*socket.Socket, T) error, args ...any) (err error) {

	context := conn.Data()
	contextMap, ok := context.(gin.H)
	if ok {
		arg, ok := args[0].(T)

		if ok && contextMap["user"] != nil {
			token, ok := conn.Handshake().Query.Get("token")
			if ok {
				mm, err := messengermanager.NewMessengerManager(nil)
				if err == nil {
					// Reset Tiker
					err = mm.ResetUserTime(token)

					if err == nil {
						return next(conn, arg)
					}
				}
			}

		} else {
			return errors.New("incorrect object sent")
		}

	}
	return errors.New("the connection should be bound to a token")
}

// HasUserMiddleware checks if user is loged in
func hasUserMiddlewareNoParam(conn *socket.Socket, next func(*socket.Socket) error) (err error) {
	context := conn.Data()
	contextMap, ok := context.(gin.H)
	if ok && contextMap["user"] != nil {
		token, ok := conn.Handshake().Query.Get("token")
		if ok {
			mm, err := messengermanager.NewMessengerManager(nil)
			if err == nil {
				// Reset Tiker
				err = mm.ResetUserTime(token)

				if err == nil {
					return next(conn)
				}
			}
			return err
		}
	}

	return errors.New("connection should be bound to a token")
}
