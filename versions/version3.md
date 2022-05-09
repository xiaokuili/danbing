## 功能
1. 实现reader和writer接口 taskgroup中
2. 注意goroutine在并发过程中对指针的影响

## datax实现分析
- 开线程池
```
for (Configuration taskGroupConfiguration : configurations) {
    TaskGroupContainerRunner taskGroupContainerRunner = newTaskGroupContainerRunner(taskGroupConfiguration);
    this.taskGroupContainerExecutorService.execute(taskGroupContainerRunner);
}
```
- 执行reader和writer
```
// 其实这里是通过channel控制流量, 下一版本考虑
public void doStart() {
    this.writerThread.start();

  
    this.readerThread.start();
}
 
```



## Next 
1. 监控
2. 日志
3. 这两部分如果需要单独设计接口 
4. 在这部分考虑taskgroup是否分割
5. 流量控制