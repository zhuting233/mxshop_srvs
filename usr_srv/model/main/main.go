//package main
//
//import (
//	"crypto/md5"
//	"encoding/hex"
//	"fmt"
//	"io"
//)
//
//func genMd5(code string) string {
//	Md5 := md5.New()
//	_, _ = io.WriteString(Md5, code)
//	return hex.EncodeToString(Md5.Sum(nil))
//}
//
//func main() {
//	//
//	//newLogger := logger.New(
//	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
//	//	logger.Config{
//	//		SlowThreshold:        time.Second, // Slow SQL threshold
//	//		LogLevel:             logger.Info, // Log level
//	//		ParameterizedQueries: false,
//	//		Colorful:             true, // color
//	//	},
//	//)
//	//
//	//// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
//	//dsn := "root:123456@tcp(192.168.57.128:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
//	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
//	//	NamingStrategy: schema.NamingStrategy{
//	//		SingularTable: true,
//	//	},
//	//	Logger: newLogger,
//	//})
//	//
//	//if err != nil {
//	//	panic("连接数据库失败\n" + err.Error())
//	//}
//	//
//	//db.AutoMigrate(&model.User{})
//	fmt.Println(genMd5("zhuting2000"))
//}

package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"strings"
)

func main() {
	// Using the default options
	salt, encodedPwd := password.Encode("generic password", nil)
	check := password.Verify("generic password", salt, encodedPwd, nil)
	fmt.Println(check) // true

	// Using custom options
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd = password.Encode("generic password", options)
	newPassword := fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(newPassword)

	passwordInfo := strings.Split(newPassword, "$")

	check = password.Verify("generic password", passwordInfo[1], passwordInfo[2], options)
	fmt.Println(check) // true
}
