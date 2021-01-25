package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	/*
	 * 起始的时间戳
	 */
	start_stmp int64 = 1480166465631

	/*
	* 每一部分占用的位数
	*/
	sequence_bit   uint64 = 12 //序列号占用的位数
	machine_bit    uint64 = 5  //机器标识占用的位数
	datacenter_bit uint64 = 5  //数据中心占用的位数

	/*
	 * 每一部分的最大值
	 */
	max_datacenter_num = -1 ^ (-1 << datacenter_bit)
	max_machine_num    = -1 ^ (-1 << machine_bit)
	max_sequence       = -1 ^ (-1 << sequence_bit)

	/*
	 * 每一部分向左的位移
	 */
	machine_left    = sequence_bit
	datacenter_left = sequence_bit + machine_bit
	timestmp_left   = datacenter_left + datacenter_bit
)

var (
	datacenterId int64      //数据中心
	machineId    int64      //机器标识
	sequence     int64 = 0  //序列号
	lastStmp     int64 = -1 //上一次时间戳
)

var lock = sync.RWMutex{}

type SnowFlack struct {
	DatacenterId int64 //数据中心
	MachineId    int64 //机器标识
}

/*
 * 初始化函数
 */
func (s *SnowFlack) Init() error {
	if datacenterId > max_datacenter_num || datacenterId < 0 {
		return errors.New("datacenterId can't be greater than MAX_DATACENTER_NUM or less than 0")
	}
	if machineId > max_machine_num || machineId < 0 {
		return errors.New("machineId can't be greater than MAX_MACHINE_NUM or less than 0")
	}
	datacenterId = s.DatacenterId
	machineId = s.DatacenterId
	return nil
}

/*
 * 产生下一个ID
 */
func (s *SnowFlack) NextId() (int64, error) {
	lock.Lock()
	defer lock.Unlock()

	currStmp := s.getNewstmp()
	if currStmp < lastStmp {
		return 0, errors.New("Clock moved backwards.  Refusing to generate id")
	}

	if currStmp == lastStmp {
		//相同毫秒内，序列号自增
		sequence = (sequence + 1) & max_sequence;
		//同一毫秒的序列数已经达到最大
		if sequence == 0 {
			currStmp = s.getNextMill()
		}
	} else {
		//不同毫秒内，序列号置为0
		sequence = 0
	}

	lastStmp = currStmp

	return (currStmp-start_stmp)<<timestmp_left | //时间戳部分
		datacenterId<<datacenter_left | //数据中心部分
		machineId<<machine_left | //机器标识部分
		sequence, //序列号部分
		nil
}

func (s *SnowFlack) getNextMill() int64 {
	var mill = s.getNewstmp()
	for ; mill < lastStmp; {
		mill = s.getNewstmp()
	}
	return mill
}

/*
 * 获取毫秒
 */
func (s *SnowFlack) getNewstmp() int64 {
	return time.Now().UnixNano() / 1e6
}
