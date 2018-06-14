## Raft 

分布式一致性算法的一种，相比较Paxos，Raft理解起来容易多了。

Raft三种角色：  
+ follow：（群众）
+ candidate：（候选人）  
+ leader：（领袖）  


#### 参考 
1. raft实现源码：https://github.com/goraft/raft,  `go get github.com/goraft/raft`
2. 容易理解的演示动画： http://thesecretlivesofdata.com/raft/