package dbUser

import (
	"fmt"

	"UserCenter.net/server/global/config"
	"github.com/EasyGolang/goTools/mEncrypt"
)

func (dbObj *AccountType) CheckPassword(Password string) (resErr error) {
	resErr = nil

	AccountData := dbObj.Data

	if AccountData.Password != mEncrypt.AseDecrypt(Password, config.SecretKey) {
		return fmt.Errorf("密码错误")
	} 

	return
}
