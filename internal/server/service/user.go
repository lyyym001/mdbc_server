package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"mdbc_server/internal/config"
	"mdbc_server/internal/helper"
	"mdbc_server/internal/models"
	"mdbc_server/pb"
	"net/http"
)

func UserRegister(c *gin.Context) {
	in := new(UserRegisterRequest)
	//fmt.Println("c = ", c.Request)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}
	if in.Username == "" || in.Password == "" || in.NickName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填信息为空",
		})
		return
	}
	in.Password = helper.GetMd5(in.Password)

	//写入对抗数据
	u := models.UserBasic{
		Username:    in.Username,
		Password:    in.Password,
		NickName:    in.NickName,
		Role:        in.Role,
		AccountType: 0,
		Identify:    in.Username,
	}
	if err := models.DB.Create(&u).Where("username != ", in.Username).Error; err != nil {
		fmt.Println("insert user error")
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "注册失败" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "注册成功",
		})
	}

}

func UserLogin(c *gin.Context) {
	in := new(UserLoginRequest)
	//fmt.Println("c = ", c.Request)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}
	if in.Username == "" || in.Password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填信息为空",
		})
		return
	}
	in.Password = helper.GetMd5(in.Password)

	data := new(models.UserBasic)
	err = models.DB.Where("username = ? AND password = ? ", in.Username, in.Password).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户名或密码错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get UserBasic Error:" + err.Error(),
		})
		return
	}

	token, err := helper.GenerateToken(data.ID, data.Username, data.Identify, data.AccountType)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GenerateToken Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": token,
	})
}

func GetUserConfig(c *gin.Context) {

	//uc := c.MustGet("user_claims").(*helper.UserClaims)
	//if uc == nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": -1,
	//		"msg":  "tokenError",
	//	})
	//	return
	//}
	//
	//data := new(models.UserBasic)
	//err := models.DB.Where("id = ?", uc.Id).First(&data).Error
	//if err != nil {
	//	if err == gorm.ErrRecordNotFound {
	//		c.JSON(http.StatusOK, gin.H{
	//			"code": -1,
	//			"msg":  "id不存在",
	//		})
	//		return
	//	}
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": -1,
	//		"msg":  "Get UserBasic Error:" + err.Error(),
	//	})
	//	return
	//}

	jsonData := pb.UserConfig{
		Blood:     config.YamlConfig.Conf.Blood,
		ReadyTime: config.YamlConfig.Conf.ReadyTime,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(jsonBytes)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": jsonString,
	})
}
