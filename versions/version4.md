## 监控
1. datax实现
2. golang类库选型


## datax 实现
### 模块设计
##### Communication: record static  
    - singleton: Map<Integer, Communication> taskGroupCommunicationMap = new ConcurrentHashMap<Integer, Communication>();   
    - updateTaskGroupCommunication(final int taskGroupId, final Communication communication)  
##### AbstractCollector: manage **map[taskgroup]communication** and **map[job]communication**  
    -  private Map<Integer, Communication> taskCommunicationMap = new ConcurrentHashMap<Integer, Communication>();  
    -  getTGCommunication(Integer taskGroupId) -> taskGroupCommunicationMap  

##### AbstractContainerCommunicator  
    - attr:collector -> AbstractCollector  
    - collect()  
    - report()  
    - registerCommunication()  


### 流程
1. 生成流程
```
//  1. init scheduler -> AloneJobContainerCommunicator(configuration)
scheduler = initStandaloneScheduler(this.configuration);
AbstractContainerCommunicator containerCommunicator = new StandAloneJobContainerCommunicator(configuration);

// 2. register taskGroupContainer -> Communication.singleton
this.containerCommunicator.registerCommunication(configurations);

// 3. init taskgroup communication 
initCommunicator(configuration);
super.setContainerCommunicator(new StandaloneTGContainerCommunicator(configuration));

// 4. register communication -> map[task]communication 
this.containerCommunicator.registerCommunication(taskConfigs);


```
2. 调用流程
```
// schedule, 剩下部分不展示
// 1. collect job stat
// 收集taskGroupCommunicationMap数据
Communication nowJobContainerCommunication = this.containerCommunicator.collect();
// 2. get report then report 
Communication reportCommunication = CommunicationTool.getReportCommunication(nowJobContainerCommunication, lastJobContainerCommunication, totalTasks);

this.containerCommunicator.report(reportCommunication);
```

## TODO
1. package 
   1. https://github.com/rcrowley/go-metrics
2. 功能
   1. 页面, Prometheus 
3. 其他资料
   1. https://talks.golang.org/2012/10things.slide#3 -> 比较简短的小的项目