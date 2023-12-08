package tag

import (
	"cmdb/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AddTag(c *gin.Context) {
	var tag models.Tag
	if err := c.BindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	err := tag.Save()
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

func GetTagList(c *gin.Context) {
	var tag models.Tag
	var pageStr string
	var pageSizeStr string
	pageStr = c.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response := models.TagRsp{
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
		response := models.TagRsp{
			Code: http.StatusBadRequest,
			Msg:  "Wrong request parameter: size",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if tags, total, err := tag.GetList(page, pageSize); err != nil {
		response := models.TagRsp{
			Code: http.StatusInternalServerError,
			Msg:  "Internal error",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	} else {
		response := models.TagRsp{
			Code:     http.StatusOK,
			Msg:      "",
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Tags:     tags,
		}
		c.JSON(http.StatusOK, response)
		return
	}
}

func UpdateTag(c *gin.Context) {
	var tag models.Tag
	if err := c.BindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	id := tag.Id
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Missing request data: id.",
			"data":    nil,
		})
		return
	}

	if err := tag.Update(); err != nil {
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

func DeleteTag(c *gin.Context) {
	var tag models.Tag
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
	if err := tag.Delete(id); err != nil {
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
