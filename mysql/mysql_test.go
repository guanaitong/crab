package mysql_test

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/guanaitong/crab/mysql"
	"github.com/guanaitong/crab/system"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
func TestDecrypt(t *testing.T) {
	epwd := "E+lDullrAU/qV1MVqR7L0GrbBkFHWaftsKTVni3ooL90/PZyH/VpcKF/HqJqzAyzoHI8vR+tawW/kE5sgRcpVkYivugNhWhEtnQpbRNjvnkCd8OcyuhjEVnrzDg4iNtJ4+RWKq37vb4aXU1/skmXDLd1Jf2ZNYndzTgHM1EbP6Ac0KqWzpeS4o2QxtX4E1nzdrxCOtEYtTewtXiaxA4kHdVb6fIkLa/OvY2xDNOQZKhlw9IU6LC3Ypq8qqQPq1dCW+Y/TzktZcbKVmQ0aHchPLuWpiO2VNwojHu7hiD7ZiNsELiDvose8iNNSwwpfTKbIODqjtgBrRWD/VLjCbMcxg=="
	pwd := decrypt(epwd)
	if pwd == "" {
		t.Error("error")
	}
	fmt.Println(pwd)
}*/

func TestGetDataSourceConfig(t *testing.T) {
	system.SetupAppName("approval")
	d := mysql.GetDefaultDataSourceConfig()
	//d := mysql.GetDataSourceConfig("datasource.json")
	assert.NotNil(t, d)

	m := d.MasterDataSourceName()
	assert.NotEmpty(t, m)

	s := d.SlaveDataSourceName()
	assert.NotEmpty(t, s)

	db, err := d.OpenMaster()
	if !assert.NoError(t, err) {
		t.Error(err.Error())
		return
	}
	t.Log(db)

	err = db.Ping()
	if !assert.NoError(t, err) {
		t.Error(err.Error())
		return
	}

	//db.Exec()
}
