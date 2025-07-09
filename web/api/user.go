package api

import (
	"VulnFusion/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

// HandleListAllUsers godoc
// @Summary 获取所有用户列表
// @Description 管理员查看所有用户信息
// @Tags 用户管理
// @Security BearerToken
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} gin.H{"error": "获取用户列表失败"}
// @Router /admin/users [get]
func HandleListAllUsers(c *gin.Context) {
	users, err := models.ListAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// HandleDeleteUserByID godoc
// @Summary 删除用户
// @Description 管理员根据 ID 删除用户
// @Tags 用户管理
// @Security BearerToken
// @Param id path int true "用户ID"
// @Produce json
// @Success 200 {object} gin.H{"message": "删除成功"}
// @Failure 400 {object} gin.H{"error": "无效的用户 ID"}
// @Failure 500 {object} gin.H{"error": "删除失败"}
// @Router /admin/users/{id} [delete]
func HandleDeleteUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户 ID"})
		return
	}

	err = models.DeleteUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// HandleUpdateUserByID godoc
// @Summary 更新用户信息
// @Description 管理员根据 ID 更新用户角色或密码
// @Tags 用户管理
// @Security BearerToken
// @Param id path int true "用户ID"
// @Param body body object{role=string,password=string} true "更新字段（可选 role, password）"
// @Produce json
// @Success 200 {object} gin.H{"message": "更新成功"}
// @Failure 400 {object} gin.H{"error": "请求参数错误"}
// @Failure 500 {object} gin.H{"error": "更新失败"}
// @Router /admin/users/{id} [put]
func HandleUpdateUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户 ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体"})
		return
	}

	allowed := map[string]bool{"role": true, "password": true}
	for key := range updates {
		if !allowed[key] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "禁止修改字段: " + key})
			return
		}
	}

	if pwd, ok := updates["password"].(string); ok && pwd != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		updates["password"] = string(hash)
	}

	err = models.UpdateUserByID(uint(id), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// HandleResetPasswordByID godoc
// @Summary 重置用户密码
// @Description 管理员根据用户 ID 重置密码
// @Tags 用户管理
// @Security BearerToken
// @Param id path int true "用户ID"
// @Param body body object{password=string} true "新密码"
// @Produce json
// @Success 200 {object} gin.H{"message": "密码重置成功"}
// @Failure 400 {object} gin.H{"error": "请求错误"}
// @Failure 500 {object} gin.H{"error": "重置失败"}
// @Router /admin/users/{id}/password [put]
func HandleResetPasswordByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户 ID"})
		return
	}

	var body struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码不能为空"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	err = models.UpdateUserByID(uint(id), map[string]interface{}{"password": string(hash)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码重置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}
