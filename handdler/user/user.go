package user

import (
	"cmdb/models"
	"cmdb/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Login(c *gin.Context) {
	var AuthInfo models.User
	if err := c.BindJSON(&AuthInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if AuthInfo.Username == "" || AuthInfo.Password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "username and password must be offered",
			"data":    nil,
		})
		return
	}

	if !utils.JudgeBase64(AuthInfo.Password) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Password is not a valid base64 string",
			"data":    nil,
		})
		return
	}

	encodePwd, err := utils.Encrypt(AuthInfo.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	AuthInfo.Password = encodePwd

	data := make(map[string]interface{})

	isExist, userRole, userId, username := AuthInfo.CheckUser()
	if isExist {
		token, err := utils.GenerateToken(AuthInfo.Username, AuthInfo.Password)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
				"data":    nil,
			})
			return
		} else {
			data["token"] = token
			data["role"] = userRole
			data["userid"] = userId
			data["username"] = username
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": "",
				"data":    data,
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "User not found or invalid username and password",
			"data":    nil,
		})
	}
}

func AddUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "username and password must be offered",
			"data":    nil,
		})
		return
	}
	if user.NameAlias == "" {
		user.NameAlias = user.Username
	}

	if !utils.JudgeBase64(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Password is not a valid base64 string",
			"data":    nil,
		})
		return
	}

	encodePwd, err := utils.Encrypt(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	user.Password = encodePwd

	if err := user.Add(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
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

func UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "The userid must offered: id",
			"data":    nil,
		})
		return
	}
	if user.Id == 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "admin is a system default user,unable to update it",
			"data":    nil,
		})
		return
	}
	if err := user.Update(); err != nil {
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

func DeleteUser(c *gin.Context) {
	var user models.User
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
	if id == 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "admin is a system default user,unable to delete it",
			"data":    nil,
		})
		return
	}
	if err := user.Delete(id); err != nil {
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

func GetUserList(c *gin.Context) {
	var user models.User
	var pageStr string
	var pageSizeStr string
	pageStr = c.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response := models.UserRsp{
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
		response := models.UserRsp{
			Code: http.StatusBadRequest,
			Msg:  "Wrong request parameter: size",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if users, total, err := user.GetList(page, pageSize); err != nil {
		response := models.UserRsp{
			Code: http.StatusInternalServerError,
			Msg:  "Internal error",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	} else {
		response := models.UserRsp{
			Code:     http.StatusOK,
			Msg:      "",
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Users:    users,
		}
		c.JSON(http.StatusOK, response)
		return
	}
}

func RestPwd(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "The userid must offered: id",
			"data":    nil,
		})
		return
	}
	if user.Id == 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "admin is a system default user,unable to update it",
			"data":    nil,
		})
		return
	}

	if !utils.JudgeBase64(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Password is not a valid base64 string",
			"data":    nil,
		})
		return
	}

	encodePwd, err := utils.Encrypt(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	user.Password = encodePwd

	if err := user.UpdatePwd(); err != nil {
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
