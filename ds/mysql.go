package ds

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/guanaitong/crab/gconf"
	"github.com/guanaitong/crab/util"
	"github.com/guanaitong/crab/util/format"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultDataSourceKey = "datasource.json"

type MysqlServer struct {
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	Ip        string `json:"ip"`
	Port      string `json:"port"`
	Version   string `json:"version"`
	GroupName string `json:"groupName"`
	Role      string `json:"role"`
}

type DataSourceConfig struct {
	DbName            string            `json:"dbName"`
	Username          string            `json:"username"`
	EncryptedPassword string            `json:"encryptedPassword"`
	Password          string            `json:"password"`
	GroupName         string            `json:"groupName"`
	MysqlServers      []*MysqlServer    `json:"mysqlServers"`
	Params            map[string]string `json:"params"`
}

func (dataSourceConfig *DataSourceConfig) OpenMaster() (db *sql.DB, err error) {
	return dataSourceConfig.open(false)
}

func (dataSourceConfig *DataSourceConfig) OpenSlave() (db *sql.DB, err error) {
	return dataSourceConfig.open(true)
}

func (dataSourceConfig *DataSourceConfig) open(preferSlave bool) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", dataSourceConfig.dataSourceName(preferSlave))
	if err == nil {
		db.SetMaxOpenConns(dataSourceConfig.getParamValue("maxOpenConns", 20))
		db.SetMaxIdleConns(dataSourceConfig.getParamValue("maxIdleConns", 3))
		db.SetConnMaxLifetime(time.Second * time.Duration(dataSourceConfig.getParamValue("maxIdleConns", 1200)))
	}
	return
}

func (dataSourceConfig *DataSourceConfig) MasterDataSourceName() string {
	return dataSourceConfig.dataSourceName(false)
}

func (dataSourceConfig *DataSourceConfig) SlaveDataSourceName() string {
	return dataSourceConfig.dataSourceName(true)
}

func (dataSourceConfig *DataSourceConfig) dataSourceName(preferSlave bool) string {
	var pwd = decrypt(dataSourceConfig.EncryptedPassword)
	if pwd == "" {
		pwd = dataSourceConfig.Password
	}
	mysqlServer := dataSourceConfig.getMysqlServer(preferSlave)
	if mysqlServer == nil {
		panic("there is no mysql server")
	}
	var host = mysqlServer.Domain
	if host == "" {
		host = mysqlServer.Ip
	}
	timezone := "'Asia/Shanghai'"
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dataSourceConfig.Username,
		pwd,
		host,
		mysqlServer.Port,
		dataSourceConfig.DbName,
	) + "?charset=utf8mb4&parseTime=true&loc=Local&time_zone=" + url.QueryEscape(timezone)
}

func (dataSourceConfig *DataSourceConfig) getParamValue(key string, defaultValue int) int {
	v, ok := dataSourceConfig.Params[key]
	if ok {
		i, err := strconv.Atoi(v)
		if err == nil {
			return i
		}
	}
	return defaultValue
}

func (dataSourceConfig *DataSourceConfig) getMysqlServer(preferSlave bool) *MysqlServer {
	size := len(dataSourceConfig.MysqlServers)
	if size == 1 {
		return dataSourceConfig.MysqlServers[0]
	} else if size > 1 {
		var master *MysqlServer
		var slave *MysqlServer
		for _, ms := range dataSourceConfig.MysqlServers {
			if "master" == strings.ToLower(ms.Role) {
				master = ms
			} else if "slave" == strings.ToLower(ms.Role) {
				slave = ms
			}
		}

		if preferSlave && slave != nil {
			return slave
		}
		return master
	}
	return nil
}

func decrypt(encryptedPassword string) string {
	if encryptedPassword == "" {
		return ""
	}
	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return ""
	}
	publicKey := gconf.GetGlobalConfigCollection().GetConfig("publicKey")
	key, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return ""
	}
	pubKey, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return ""
	}
	pub := pubKey.(*rsa.PublicKey)
	return string(util.RsaPublicDecrypt(pub, encryptedDecodeBytes))
}

func GetDataSourceConfig(key string) *DataSourceConfig {
	if key == "" {
		panic(errors.New("data source is null"))
	}

	dataSourceConfig := new(DataSourceConfig)
	configValue := gconf.GetCurrentConfigCollection().GetConfig(key)
	err := format.AsJson(configValue, dataSourceConfig)
	if err != nil {
		panic(err.Error())
	}
	return dataSourceConfig
}

func GetDefaultDataSourceConfig() *DataSourceConfig {
	return GetDataSourceConfig(defaultDataSourceKey)
}
