## 监控
1. datax实现
2. golang类库选型


## datax 实现
- 有点像插件的实现
```
statistic
    -- communication: 收集器数据及结构定义, counter, state, msg, throwable, timestamp,（单例？）
    -- container: 收集器行为定义, 注册以及收集行为的实现
        -- collect: 维护这样数据结构Map<Integer, Communication> 
        -- communicator: 初始化collect, 然后将communicator实例扔进去
            -- taskgroup: registerCommunication()
            -- job
    -- report: report函数实现

collect 

container
    class: AbstractContainerCommunicator 

```
- 各个组件注册
```
// job
// 初始化一个job的communicator
tempContainerCollector = new StandAloneJobContainerCommunicator(configuration);
super.setContainerCommunicator(tempContainerCollector);

// scheduler
// 处理一些
AbstractContainerCommunicator containerCommunicator = new StandAloneJobContainerCommunicator(configuration);
super.setContainerCommunicator(containerCommunicator)

// tg 
this.containerCommunicator.registerCommunication(configurations);

// task
private TaskMonitor taskMonitor = TaskMonitor.getInstance();

```