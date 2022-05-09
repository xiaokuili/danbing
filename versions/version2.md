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
- 合并reader和writer
```
// 这里没有对应规则
// 也就意味着reader和writer只能一一对应关系
for (int i = 0; i < readerTasksConfigs.size(); i++) {
            Configuration taskConfig = Configuration.newDefault();
            taskConfig.set(CoreConstant.JOB_READER_NAME,
                    this.readerPluginName);
            taskConfig.set(CoreConstant.JOB_READER_PARAMETER,
                    readerTasksConfigs.get(i));
            taskConfig.set(CoreConstant.JOB_WRITER_NAME,
                    this.writerPluginName);
            taskConfig.set(CoreConstant.JOB_WRITER_PARAMETER,
                    writerTasksConfigs.get(i));

            if(transformerConfigs!=null && transformerConfigs.size()>0){
                taskConfig.set(CoreConstant.JOB_TRANSFORMER, transformerConfigs);
            }

            taskConfig.set(CoreConstant.TASK_ID, i);
            contentConfigs.add(taskConfig);
        }
```
2. taskgroup分配优化
- 判断分组数量
``` 
// 上面配置了需要x个任务, 也就是x个channel
// 下面配置了一个组有y个channel, 也就是 x/y 并发来执行
int channelsPerTaskGroup = this.configuration.getInt(
                CoreConstant.DATAX_CORE_CONTAINER_TASKGROUP_CHANNEL, 5);

// 要么基于上面的计算， 要么基于配置
int taskNumber = this.configuration.getList(
                CoreConstant.DATAX_JOB_CONTENT).size();

int taskGroupNumber = (int) Math.ceil(1.0 * channelNumber / channelsPerTaskGroup);

// 只配置需要多少并发和拆分多少任务数量	
```
- 划分任务
```
// 实现了这个
    /**
     * /**
     * 需要实现的效果通过例子来说是：
     * <pre>
     * a 库上有表：0, 1, 2
     * b 库上有表：3, 4
     * c 库上有表：5, 6, 7
     *
     * 如果有 4个 taskGroup
     * 则 assign 后的结果为：
     * taskGroup-0: 0,  4,
     * taskGroup-1: 3,  6,
     * taskGroup-2: 5,  2,
     * taskGroup-3: 1,  7
     *
     * </pre>
     */
但是具体使用文档没有找到, 这里直接划分
```


## 其他
1. 资源不同决定了 reader和writer必须一一对应
2. 读写只有两种情况
   1. 读慢写快, 没办法只能等
   2. 读快写慢, 读速度放慢



## Next 
1. 配置反正都改了, 不如不用json
