package dbUser

import (
	"fmt"

	"UserCenter.net/server/global"
	"github.com/EasyGolang/goTools/mTime"
	"go.mongodb.org/mongo-driver/bson"
)

func (dbObj *AccountType) ChangePassword(Password string) (resErr error) {
	resErr = nil 

	db := dbObj.DB

	err := db.Ping()
	if err != nil {
		resErr = fmt.Errorf("修改密码,数据库连接错误 %+v", err)
		global.LogErr(resErr)
		return
	}

	FK := bson.D{{
		Key:   "UserID",
		Value: dbObj.UserID,
	}}

	UK := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "Password",
					Value: Password,
				},
			},
		},
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "UpdateTime",
					Value: mTime.GetUnixInt64(),
				},
			},
		},
	}
	_, err = db.Table.UpdateOne(db.Ctx, FK, UK)
	if err != nil {
		resErr = fmt.Errorf("数据库更新失败 %+v", err)
		global.LogErr(resErr)
		return
	}

	dbObj.Update()

	return
}
