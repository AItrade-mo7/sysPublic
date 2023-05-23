package dbUser

import (
	"fmt"
	"strings"

	"UserCenter.net/server/global"
	"UserCenter.net/sysPublic/dbType"
	"UserCenter.net/sysPublic/taskPush"
	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mTime"
)

type RegisterOpt struct {
	Email          string
	EntrapmentCode string
}

func (dbObj *AccountType) Register(opt RegisterOpt) (resErr error) {
	resErr = nil

	// 检查数据库连接状态
	db := dbObj.DB
	err := db.Ping()
	if err != nil {
		resErr = fmt.Errorf("注册用户,数据库连接错误 %+v", err)
		global.LogErr(resErr)
		return
	}

	if len(dbObj.Data.UserID) > 0 {
		resErr = fmt.Errorf("该账号已注册，请直接登录")
		return
	}

	newPwd := mEncrypt.RandStr(8) // 生成密码
	UserEmail := []string{}
	UserEmail = append(UserEmail, opt.Email)
	str_arr := strings.Split(opt.Email, `@`) // 获取邮箱前缀
	NickName := "AItrade用户"
	if len(str_arr) > 0 {
		NickName = str_arr[0]
	}

	var Body dbType.UserTable
	Body.UserID = mEncrypt.GetUUID()                 // 生成 UserID
	Body.Email = opt.Email                           // 插入邮箱
	Body.UserEmail = UserEmail                       // 插入邮箱
	Body.Avatar = "//file.mo7.cc/AItrade/avatar.png" // 生成默认头像
	Body.NickName = NickName                         // 生成昵称,昵称应该为邮箱前缀
	Body.CreateTime = mTime.GetTime().TimeUnix       // 生成创建时间
	Body.UpdateTime = mTime.GetTime().TimeUnix       // 生成更新时间
	Body.EntrapmentCode = opt.EntrapmentCode         // 防伪标识符
	Body.Password = mEncrypt.MD5(newPwd)             // 密码加密存储

	// 2. 插入数据库
	_, err = db.Table.InsertOne(db.Ctx, Body)
	if err != nil {
		resErr = fmt.Errorf("注册,插入数据失败 %+v", err)
		global.LogErr(resErr)
		return
	}

	// 1. 发送邮件告知密码
	taskPush.RegisterEmail(taskPush.RegisterEmailOpt{
		To:             Body.Email,
		Password:       newPwd,
		EntrapmentCode: Body.EntrapmentCode,
	})

	dbObj.UserID = Body.UserID
	dbObj.Update()

	return
}
