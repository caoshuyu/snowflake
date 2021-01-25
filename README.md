# snowflake

### 概述
分布式系统中，有一些需要使用全局唯一ID的场景，这种时候为了防止ID冲突可以使用36位的UUID，但是UUID有一些缺点，首先他相对比较长，另外UUID一般是无序的。

有些时候我们希望能使用一种简单一些的ID，并且希望ID能够按照时间有序生成。

而twitter的snowflake解决了这种需求，最初Twitter把存储系统从MySQL迁移到Cassandra，因为Cassandra没有顺序ID生成机制，所以开发了这样一套全局唯一ID生成服务。

### 结构
snowflake的结构如下(每部分用-分开):

```text
0 - 0000000000 0000000000 0000000000 0000000000 0 - 00000 - 00000 - 000000000000
```
第一位为未使用，接下来的41位为毫秒级时间(41位的长度可以使用69年)，然后是5位datacenterId和5位workerId(10位的长度最多支持部署1024个节点） ，最后12位是毫秒内的计数（12位的计数顺序号支持每个节点每毫秒产生4096个ID序号）

一共加起来刚好64位，为一个Long型。(转换成字符串长度为18)

snowflake生成的ID整体上按照时间自增排序，并且整个分布式系统内不会产生ID碰撞（由datacenter和workerId作区分），并且效率较高。据说：snowflake每秒能够产生26万个ID。

### 调用方法
```go
package main

import (
	"fmt"
	"github.com/caoshuyu/snowflake"
)

var snowFlack snowflake.SnowFlack

func init() {
	snowFlack = snowflake.SnowFlack{
		DatacenterId: 1,
		MachineId:    1,
	}
	err := snowFlack.Init()
	if nil != err {
		panic(err)
		return
	}
}

func main() {
	id, err := snowFlack.NextId()
	if nil != err {
		fmt.Println(err)
		return
	}
	fmt.Println(id)
}
```
