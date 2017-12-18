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




参考：https://blog.jcoglan.com/2017/02/12/the-myers-diff-algorithm-part-1/