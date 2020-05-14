## 依赖库：

```sh
go get -u gopkg.in/ini.v1

go get -u github.com/kataras/iris

go get -u github.com/jinzhu/gorm 

go get -u github.com/unknwon/com
```



## 服务端优化思路

![image-20200512092901398](img/README/image-20200512092901398.png)

#### 分布式验证

- 一致性 Hash 算法

  > 用途：快速定位资源，均匀分布；
  >
  > 场景：分布式存储，分布式缓存，负载均衡；

  

  **[一致性hash算法原理及golang实现](https://segmentfault.com/a/1190000013533592)**

  

  结构示意图

  ![image-20200514143301254](img/README/image-20200514143301254.png)

