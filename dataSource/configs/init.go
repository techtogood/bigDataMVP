package configs


import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//Config配置
var Config configs

func init() {
	readTimes := 0
readFile:
	f, err := ioutil.ReadFile("./configs/configs.yaml")
	if err != nil {
		if readTimes >= 3 {
			panic("读取配置文件错误，启动失败")
		} else {
			readTimes++
			goto readFile
		}
	}

	yaml.Unmarshal(f, &Config)
}

type env struct {
	Prod bool `yaml:"prod"`
}

type mysqlDB struct {
	IP       string `yaml:"ip"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (db *mysqlDB) Dsn() string {
	return db.Username + ":" + db.Password + "@tcp(" + db.IP + ":" + db.Port + ")/stock?parseTime=true"
}


type configs struct {
	MysqlDB mysqlDB `yaml:"mysql_db"`
	Env     env     `yaml:"env"`
}

