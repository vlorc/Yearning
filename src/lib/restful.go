package lib

import "github.com/gin-gonic/gin"

type RestfulAPI struct {
	Get    gin.HandlerFunc
	Post   gin.HandlerFunc
	Delete gin.HandlerFunc
	Put    gin.HandlerFunc
}

func (a RestfulAPI) Route(r gin.IRouter, prefix string) {
	if nil != a.Get {
		r.GET(prefix, a.Get)
	}
	if nil != a.Post {
		r.POST(prefix, a.Post)
	}
	if nil != a.Put {
		r.PUT(prefix, a.Put)
	}
	if nil != a.Delete {
		r.DELETE(prefix, a.Delete)
	}
}
