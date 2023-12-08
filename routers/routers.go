package routers

import (
	"cmdb/handdler/host"
	"cmdb/handdler/netdevice"
	"cmdb/handdler/tag"
	"cmdb/handdler/user"
	"cmdb/handdler/vm"
	"cmdb/handdler/webssh"
	"cmdb/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func Routers() *gin.Engine {
	routers := gin.Default()
	routers.Use(gin.Recovery())
	routers.Use(Cors())
	routers.Use(middleware.Jwt())

	userRouter := routers.Group("/api/v1")
	{
		userRouter.POST("/login", user.Login)
		userRouter.POST("/user/add", user.AddUser)
		userRouter.POST("/user/update", user.UpdateUser)
		userRouter.DELETE("/user/delete", user.DeleteUser)
		userRouter.GET("/user/list", user.GetUserList)
	}
	hostRouter := routers.Group("/api/v1/resource/host")
	{
		hostRouter.POST("/add", host.AddHost)
		hostRouter.GET("/list", host.GetHostList)
		hostRouter.GET("/detail", host.GetHostByID)
		hostRouter.POST("/update", host.UpdateHost)
		hostRouter.DELETE("/delete", host.DeleteHost)

	}
	vmRouter := routers.Group("/api/v1/resource/vm")
	{
		vmRouter.POST("/add", vm.AddVm)
		vmRouter.GET("/list", vm.GetVmList)
		vmRouter.GET("/detail", vm.GetVmByID)
		vmRouter.POST("/update", vm.UpdateVm)
		vmRouter.DELETE("/delete", vm.DeleteVm)
	}
	netDeviceRouter := routers.Group("/api/v1/resource/netdevice")
	{
		netDeviceRouter.POST("/add", netdevice.AddNetDevice)
		netDeviceRouter.GET("/list", netdevice.GetNetDeviceList)
		netDeviceRouter.GET("/detail", netdevice.GetNetDeviceByID)
		netDeviceRouter.POST("/update", netdevice.UpdateNetDevice)
		netDeviceRouter.DELETE("/delete", netdevice.DeleteNetDevice)
		netDeviceRouter.POST("/interfacetopology/add", netdevice.AddInterfaceTopology)
		netDeviceRouter.POST("/interfacetopology/update", netdevice.UpdateInterfaceTopology)
		netDeviceRouter.DELETE("/interfacetopology/delete", netdevice.DeleteInterfaceTopology)
		netDeviceRouter.DELETE("/interfacetopology/srcdevice/delete", netdevice.DeleteInterfaceTopologyBySrcDeviceId)
		netDeviceRouter.GET("/interfacetopology/detail", netdevice.GetInterfaceTopologyBySrcDeviceId)
	}
	tagRouter := routers.Group("/api/v1/pub/tag")
	{
		tagRouter.POST("/add", tag.AddTag)
		tagRouter.GET("/list", tag.GetTagList)
		tagRouter.POST("/update", tag.UpdateTag)
		tagRouter.DELETE("/delete", tag.DeleteTag)
	}
	websshRouter := routers.Group("/api/v1/webssh")
	{
		websshRouter.GET("/ssh", webssh.WsSsh)
	}
	return routers
}
