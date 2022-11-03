package middleware

import (
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
)

func Session(c *gin.Context) {
	store := ginsession.FromContext(c)
	storeSessionId := store.SessionID()
	cookieSessionId, isExists := store.Get("id")
	if isExists {
		if cookieSessionId != storeSessionId {
			c.AbortWithStatus(401)
			return
		}
	} else {
		c.AbortWithStatus(401)
		return
	}

	c.Next()
}
