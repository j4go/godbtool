package dbops

import (
	"flag"
	"fmt"
	"github.com/akkuman/parseConfig"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"logtools/models"
)

var DB *gorm.DB
var CON parseConfig.Config
var M string
var C string
var L string

func init() {
	// 从命令行读取参数 如环境变量 从文件读取配置
	flag.StringVar(&M, "m", "release", "set gin env debug/release")
	flag.StringVar(&C, "c", "config_release.json", "set config file")
	flag.StringVar(&L, "l", "logs/", "set logs directory")
	flag.Parse()

	CON = parseConfig.New(C)
	mysqlUser := CON.Get("mysqlUser").(string)
	mysqlPwd := CON.Get("mysqlPwd").(string)
	fmt.Println(mysqlUser, mysqlPwd)

	// 连接mysql
	var err error
	DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/log_config?charset=utf8&parseTime=True&loc=Local", mysqlUser, mysqlPwd))
	if err != nil {
		panic(err)
	}
	//defer DB.Close()  不能用在这里??

	DB.SingularTable(true) //全局设置表名不可以为复数形式

	DB.DB().SetMaxIdleConns(20)
	DB.DB().SetMaxOpenConns(200)

	// init db with db model
	DB.AutoMigrate(&models.Person{})
	DB.AutoMigrate(&models.Meter{})

}
