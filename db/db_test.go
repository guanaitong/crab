package db

import (
	"fmt"
	"github.com/guanaitong/crab/system"
	"testing"
)

func TestDecrypt(t *testing.T) {
	epwd := "E+lDullrAU/qV1MVqR7L0GrbBkFHWaftsKTVni3ooL90/PZyH/VpcKF/HqJqzAyzoHI8vR+tawW/kE5sgRcpVkYivugNhWhEtnQpbRNjvnkCd8OcyuhjEVnrzDg4iNtJ4+RWKq37vb4aXU1/skmXDLd1Jf2ZNYndzTgHM1EbP6Ac0KqWzpeS4o2QxtX4E1nzdrxCOtEYtTewtXiaxA4kHdVb6fIkLa/OvY2xDNOQZKhlw9IU6LC3Ypq8qqQPq1dCW+Y/TzktZcbKVmQ0aHchPLuWpiO2VNwojHu7hiD7ZiNsELiDvose8iNNSwwpfTKbIODqjtgBrRWD/VLjCbMcxg=="
	pwd := decrypt(epwd)
	if pwd == "" {
		t.Error("error")
	}
	fmt.Println(pwd)
}

func TestGetDefaultDataSourceConfig(t *testing.T) {
	system.SetupAppName("userdoor")
	d := GetDefaultDataSourceConfig()
	m := d.MasterDataSourceName()
	if m == "" {
		t.Error("error")
	}
	s := d.SlaveDataSourceName()
	if s == "" {
		t.Error("error")
	}
}
