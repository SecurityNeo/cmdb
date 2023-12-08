package host

import (
	"cmdb/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary 根据ID查询主机信息
// @Tags 主机模块
// @version v1
// @Param id path int true "主机ID"
// @Success 200 object models.Response 返回值
// @Router /api/v1/resource/host/detail [get]

func GetHostByID(c *gin.Context) {
	var host models.Host
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
	data, err := host.GetById(id)
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

func AddHost(c *gin.Context) {
	var host models.Host
	if err := c.BindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	err := host.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
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

func GetHostList(c *gin.Context) {
	var host models.Host
	var pageStr string
	var pageSizeStr string
	TagIdStr := c.Query("tag")
	pageStr = c.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response := models.HostsRsp{
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
		response := models.HostsRsp{
			Code: http.StatusBadRequest,
			Msg:  "Wrong request parameter: size",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if hosts, total, err := host.GetList(page, pageSize, TagIdStr); err != nil {
		response := models.HostsRsp{
			Code: http.StatusInternalServerError,
			Msg:  "Internal error",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	} else {
		response := models.HostsRsp{
			Code:     http.StatusOK,
			Msg:      "",
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Hosts:    hosts,
		}
		c.JSON(http.StatusOK, response)
		return
	}
}

func UpdateHost(c *gin.Context) {
	var host models.Host
	if err := c.BindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	id := host.Id
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Missing request data: id.",
			"data":    nil,
		})
		return
	}

	if err := host.Update(); err != nil {
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

func DeleteHost(c *gin.Context) {
	var host models.Host
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
	if err := host.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
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
