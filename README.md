# go-workerpool

📝 go 实现的轻量级线程池（最小可行实现）

该 demo 主要用于熟悉 channel 的使用（语法 & 程序设计）😊

纯 channel + select 的实现方案 ✅
没有使用 sync 包的同步结构（Mutex、RWMutex、Cond..） ❌

workerpool 的 3 个主要部分：
· pool 的创建 New 与销毁 Free
· pool 中的 worker / goroutine 的管理
· task 的提交与调度 Schedule

❗️Special：
· 使用了“功能选项 functional option”方案，让 workerpool 支持行为定制机制
· Option 实质是一个接受 \*Pool 类型参数的函数类型
