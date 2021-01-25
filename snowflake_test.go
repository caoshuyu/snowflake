package snowflake

import (
	"fmt"
	"testing"
)

var snowFlack SnowFlack

func init() {
	snowFlack = SnowFlack{
		DatacenterId: 1,
		MachineId:    1,
	}
	err := snowFlack.Init()
	if nil != err {
		panic(err)
		return
	}
}

func TestSnowFlack_Init(t *testing.T) {
	id, err := snowFlack.NextId()
	if nil != err {
		fmt.Println(err)
		return
	}
	fmt.Println(id)
}
