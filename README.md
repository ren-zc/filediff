# filediff

一个用于比对文本文件差异的程序，使用golang实现。  
输出形式如下：
```    
    (   1,   1)   Header
    (   2,   2)   
    (   3,   3)   paragraph 1
    (   4,   4)     p1 line1
    (   5,   5)     p1 line2
    (   6,   6)     p1 line3
    (    ,   7) +   p1 line4
    (   7,   8)   p1 end
    (   8,   9)   
    (   9,    ) - paragraph 2
    (  10,    ) -   {
    (  11,    ) -       p2 line1
    (  12,    ) -       p2 line2
    (  13,    ) -       p2 line3
    (  14,    ) -   }
    (    ,  10) + //paragraph 2
    (    ,  11) + //    {
    (    ,  12) + //        p2 line1
    (    ,  13) + //        p2 line2
```

##filediff的实现过程：  
假设有两个文本文件：  
src: ABCABBA  
dst: CBABAC

如何实现从src到dst的修改过程？  
有很多实现过程，我们的目标是找出src与dst之间的最大相同行的集合。  

为了方便src与dst之间的比较，我们将其表示成如下形式的图表：  
![chart1](https://github.com/jacenr/filediff/blob/master/Screenshots/0.png)  
图中向右移动一格表示从src中"-"去一行，向下移动一格表示在dst中"+"上一行，  
斜线表示src和dst中的内容相同，不做改变。  
通过观察图表我们发现，这就是一个每条路径的权重均为1的有向图，我们的目标就是找出  
从(0,0)到(7,6)最短路径，最短路径也就是经过斜线最多的路径，或者找出斜边最多的集合，  
集合中的每个斜边组合在一起的路径方向必须是向右或者向下。

实现算法一: 每个父节点(最末的顶点除外)必有1-3个子节点，这1-3个子节点中至少一个是  
最短路径，而父节点的每个子节点也同样具有1-3个子节点，子节点的子节点中同样至少一  
个是最短路径，...就这样层层递归...，直到到达最末一个没有子节点的节点返回。
diffV1便是此方法。  
此方法最简单也最直观，当文件行数不超过1000行时也勉强可以输出结果，但是其cpu及  
内存占用过高，因为每个节点都要经过至少一次计算，文件的行越多，程序的执行时间越长  
cpu及内存占用率就越高，当文件行超过1000行的时候，cup及内存基本已经不堪重负。 

实现算法二：  
既然是找出斜边最多的有序集合，那么对于一条已知的斜边，当查找下一条斜边时，对于此斜  
边对应一个搜索范围,比如[(2,0),(3,1)],[(1,3),(2,4)]都有各自的最短路径查找范围，  
为了形式上的统一，我增加了一个节点(-1,-1)，斜边[(-1,-1),(0,0)]的查找范围是整个图：  
![chart2](https://github.com/jacenr/filediff/blob/master/Screenshots/1.png)  
而且因为每个节点的路径方向均为向右、向下、斜向下，我们可以按层查找：  
![chart2](https://github.com/jacenr/filediff/blob/master/Screenshots/2.png)  
图示为斜边[(-1,-1),(0,0)]的两个查找层次。  
对于每一个层，我们分3个方向查找，斜下，右，下，在每个方向上我们只需要保留遇到的  
第一条斜边，因为同一个方向上其他斜边能够延伸到的节点，第一条斜边也可以延伸到。  
当斜下找到一条斜边时即终止此层的查找，当斜下没有斜边时，  
再**分别**向右，向下查找，直到找到一条斜边或者到达边界  
（但不包含边界），终止此方向的查找，当向下和向右都查找完成时，  
进行下一层的查找；对于每个范围的查找，都有一个边界，对于每层每个方向的查找，  
每找到一条斜边后更新此边界，比如对于斜边[(-1,-1),(0,0)]的查找范围，  
其第一层的边界为xlimit=7，ylimit=6，  
第二层的边界为xlimit=2，ylimit=6：  
![chart2](https://github.com/jacenr/filediff/blob/master/Screenshots/3.png)  
为什么第二层的边界会被修改？  
因为对于斜边[(-1,-1),(0,0)]每查找到一条斜边，其查找范围便发生了变化：
![chart2](https://github.com/jacenr/filediff/blob/master/Screenshots/4.png)  
根据这样一个查找过程，依次递归，直到最末节点。  
为了方便实现，我们使用斜边左上的顶点代表斜边。  
此实现方法基本解决了查找最短路径的问题。实现代码是diffV2。  

diffV2补充：  
diffV2有一个缺点：  
为了节省内存，point定义的过于简单，我们需要使用递归的方法才能把完整路径描绘出来，  
这会导致在某些情况下cpu运算时间过长。  
最终的diff版本对其进行了改进，牺牲少量的内存，换取cpu时间，在point中增加了  
parent和depth字段，在输出最后的完整路径时无须递归，通过节点的parent字段便可以  
输出完整的路径。  
depth字段是为了防止重复创建point而设置的：  
```go
func checkNew(x, y int, p *point) *point {
    xyStr := strconv.Itoa(x) + strconv.Itoa(y)
    v, ok := newed[xyStr]
    if !ok {
        pNew := newpoint(x, y)
        pNew.parent = p
        pNew.depth = p.depth + 1
        newed[xyStr] = pNew
        return pNew
    }
    // 若一个节点被多条斜边查找到，仅把depth最大的一条斜边(节点)，作为其parent节点。
    if v.depth < p.depth+1 { 
        v.parent = p 
        v.depth = p.depth + 1
        return v
    }
    return nil
}
```


其他参考：https://blog.jcoglan.com/2017/02/12/the-myers-diff-algorithm-part-1/
