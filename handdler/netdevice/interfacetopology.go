package netdevice

import (
	"cmdb/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func isInArray(arr []uint64, target uint64) (int, bool) {
	for k, v := range arr {
		if v == target {
			return k, true
		}
	}
	return 0, false
}

func AddInterfaceTopology(c *gin.Context) {
	var interfaceTopologys []models.InterfaceTopology
	if err := c.BindJSON(&interfaceTopologys); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	var errData []int
	for k, v := range interfaceTopologys {
		err := v.Save()
		if err != nil {
			errData = append(errData, k)
		}
	}
	if len(errData) != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Add interface topology error.",
			"data":    errData,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "",
			"data":    nil,
		})
		return
	}
}

func UpdateInterfaceTopology(c *gin.Context) {
	var interfaceTopologys []models.InterfaceTopology
	var interfaceTopology models.InterfaceTopology
	if err := c.BindJSON(&interfaceTopologys); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if len(interfaceTopologys) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Update interface topology error",
			"data":    nil,
		})
		return
	}

	srcDeviceId := interfaceTopologys[0].SrcNetDeviceId

	oldToposOwnedBySrcDevice, err := interfaceTopology.GetList(srcDeviceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Update interface topology error,query old data error",
			"data":    nil,
		})
		return
	}

	var oldTopoIndexArr []uint64
	for _, v := range oldToposOwnedBySrcDevice {
		oldTopoIndexArr = append(oldTopoIndexArr, v.Id)
	}

	var errData []uint64
	for _, v := range interfaceTopologys {
		index, exist := isInArray(oldTopoIndexArr, v.Id)
		if exist {
			// 如果存在,就更新,并且删掉oldTopoIndexArr中指定元素
			if err := v.Update(); err != nil {
				oldTopoIndexArr = append(oldTopoIndexArr[:index], oldTopoIndexArr[index+1:]...)
				errData = append(errData, v.Id)
			} else {
				oldTopoIndexArr = append(oldTopoIndexArr[:index], oldTopoIndexArr[index+1:]...)
			}
		} else {
			// 如果不存在，就新增
			if err := v.Save(); err != nil {
				errData = append(errData, v.Id)
			}
		}
	}
	//
	for _, v := range oldTopoIndexArr {
		err := interfaceTopology.Delete(v)
		if err != nil {
			errData = append(errData, v)
		}
	}

	if len(errData) != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Update interface topology error.",
			"data":    errData,
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

func DeleteInterfaceTopology(c *gin.Context) {
	var interfaceTopology models.InterfaceTopology
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
	if err := interfaceTopology.Delete(id); err != nil {
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

func DeleteInterfaceTopologyBySrcDeviceId(c *gin.Context) {
	var interfaceTopology models.InterfaceTopology
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
	var errData []uint64
	oldToposOwnedBySrcDevice, err := interfaceTopology.GetList(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	for _, v := range oldToposOwnedBySrcDevice {
		if err := interfaceTopology.Delete(v.Id); err != nil {
			errData = append(errData, v.Id)
		}
	}
	if len(errData) != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Delete interface topology failed.",
			"data":    errData,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    nil,
	})
}

func GetInterfaceTopologyBySrcDeviceId(c *gin.Context) {
	var interfaceTopology models.InterfaceTopology
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
	data, err := interfaceTopology.GetList(id)
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
