package db

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/guanaitong/crab/gconf"
	"github.com/guanaitong/crab/util"
	"github.com/guanaitong/crab/util/format"
	"strings"
)

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
	DbName            string         `json:"dbName"`
	Username          string         `json:"username"`
	EncryptedPassword string         `json:"encryptedPassword"`
	Password          string         `json:"password"`
	GroupName         string         `json:"groupName"`
	MysqlServers      []*MysqlServer `json:"mysqlServers"`
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
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		dataSourceConfig.Username,
		pwd,
		host,
		mysqlServer.Port,
		dataSourceConfig.DbName,
	)
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
	publicKey := gconf.GetConfigCollection("golang").GetConfig("publicKey")
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

func GetDefaultDataSourceConfig() *DataSourceConfig {
	dataSourceConfig := new(DataSourceConfig)
	configValue := gconf.GetCurrentConfigCollection().GetConfig("datasource.json")
	err := format.AsJson(configValue, dataSourceConfig)
	if err != nil {
		panic(err.Error())
	}
	return dataSourceConfig
}
