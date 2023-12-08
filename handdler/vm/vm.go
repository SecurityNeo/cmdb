package vm

import (
	"cmdb/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetVmByID(c *gin.Context) {
	var vm models.Vm
	//var host models.Host
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
	data, err := vm.GetById(id)
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

func AddVm(c *gin.Context) {
	var vm models.Vm
	if err := c.BindJSON(&vm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	err := vm.Save()
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

func GetVmList(c *gin.Context) {
	var vm models.Vm
	var pageStr string
	var pageSizeStr string
	pageStr = c.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response := models.VmsRsp{
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
		response := models.VmsRsp{
			Code: http.StatusBadRequest,
			Msg:  "Wrong request parameter: size",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if vms, total, err := vm.GetList(page, pageSize); err != nil {
		response := models.VmsRsp{
			Code: http.StatusInternalServerError,
			Msg:  "Internal error",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	} else {
		response := models.VmsRsp{
			Code:     http.StatusOK,
			Msg:      "",
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Vms:      vms,
		}
		c.JSON(http.StatusOK, response)
		return
	}
}

func UpdateVm(c *gin.Context) {
	var vm models.Vm
	if err := c.BindJSON(&vm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	id := vm.Id
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Missing request data: id.",
			"data":    nil,
		})
		return
	}

	if err := vm.Update(); err != nil {
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

func DeleteVm(c *gin.Context) {
	var vm models.Vm
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
	if err := vm.Delete(id); err != nil {
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
