## 功能实现
1. split函数切分
2. taskgroup分配优化


## datax分析
1. split函数
- 定义任务数量
```
// 1.基于数据大小计算
// byte流控制
needChannelNumberByByte = (int) (globalLimitedByteSpeed / channelLimitedByteSpeed);
// record流控制
needChannelNumberByRecord = (int) (globalLimitedRecordSpeed / channelLimitedRecordSpeed);
// needChannelNumber 取这两个的最小值
this.needChannelNumber = needChannelNumberByByte < needChannelNumberByRecord ? needChannelNumberByByte : needChannelNumberByRecord;
// 2. 直接从配置文件定义
this.needChannelNumber = this.configuration.getInt(CoreConstant.DATAX_JOB_SETTING_SPEED_CHANNEL);
```
- 基于插件执行split
```
this.jobReader.split(adviceNumber)
// 拿一个stream看看
// 可以看到直接基于adviceNumber复制任务
public List<Configuration> split(int adviceNumber) {
	List<Configuration> configurations = new ArrayList<Configuration>();

	for (int i = 0; i < adviceNumber; i++) {
        configurations.add(this.originalConfig.clone());
	}
	return configurations;
}
```


## 其他
1. 资源不同决定了 reader和writer必须一一对应
2. 读写只有两种情况
   1. 读慢写快, 没办法只能等
   2. 读快写慢, 读速度放慢



## Next 
1. 配置反正都改了, 不如不用json
