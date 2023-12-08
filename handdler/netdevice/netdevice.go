package netdevice

import (
	"cmdb/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetNetDeviceByID(c *gin.Context) {
	var netDevice models.NetDevice
	reqId := c.Query("id")
	if reqId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Missing request parameter: id.",
			"data":    nil,
		})
		return
	}
	id, err := strconv.ParseUint(reqId, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	data, err := netDevice.GetById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Resource not found.",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    data,
	})
}

func AddNetDevice(c *gin.Context) {
	var netDevice models.NetDevice
	if err := c.BindJSON(&netDevice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	err := netDevice.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err,
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "success",
			"data":    nil,
		})
		return
	}
}

func GetNetDeviceList(c *gin.Context) {
	var netDevice models.NetDevice
	var pageStr string
	var pageSizeStr string
	pageStr = c.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response := models.NetDevicesRsp{
			Code: http.StatusBadRequest,
			Msg:  "Wrong request parameter: page",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	pageSizeStr = c.Query("size")
	if pageSizeStr == "" {
		pageSizeStr = "10"
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response := models.NetDevicesRsp{
			Code: http.StatusBadRequest,
			Msg:  "Wrong request parameter: size",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if netDevices, total, err := netDevice.GetList(page, pageSize); err != nil {
		response := models.NetDevicesRsp{
			Code: http.StatusInternalServerError,
			Msg:  "Internal error",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	} else {
		response := models.NetDevicesRsp{
			Code:       http.StatusOK,
			Msg:        "",
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			NetDevices: netDevices,
		}
		c.JSON(http.StatusOK, response)
		return
	}

}

func UpdateNetDevice(c *gin.Context) {
	var netDevice models.NetDevice
	if err := c.BindJSON(&netDevice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err,
			"data":    nil,
		})
		return
	}
	id := netDevice.Id
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Missing request data: id.",
			"data":    nil,
		})
		return
	}

	if err := netDevice.Update(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Resource not found.",
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "success",
			"data":    nil,
		})
		return
	}

}

func DeleteNetDevice(c *gin.Context) {
	var netDevice models.NetDevice
	reqId := c.Query("id")
	if reqId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Missing request parameter: id.",
			"data":    nil,
		})
		return
	}
	id, err := strconv.ParseUint(reqId, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err,
			"data":    nil,
		})
		return
	}
	if err := netDevice.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err,
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    nil,
	})
}
