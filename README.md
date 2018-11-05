# log-alert
日志文件监控警报

### 用途：
  用于简单的监控服务器上日志文件，根据配置的规则（正则），当日志内容符合该规则，则触发警报
 
### 使用方法
1. 复制config_example.yaml 到config.yaml
2. 修改config.yaml
  - receivers: 通知接送者的邮箱
  - notify: 发送者邮件配置
  - rules: 监控规则列表
    - file: 要监控的文件
    - rule: 正则表达式
    - desc: 描述
    - duration: 警报次数累计周期时长
    - times: 一个duration期间内，达到多少次后触发警报。当duration为0时，则和周期无关，累计达到该值就触发
    - interval: 2次警报最低间隔时间
