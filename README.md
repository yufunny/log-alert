# log-alert
日志文件监控警报

### 用途：
  用于简单的监控服务器上日志文件，根据配置的规则（正则），当日志内容符合该规则，则触发警报
 
### 使用方法
####1.部署
  - 从release中下载[二级制文件](https://github.com/yufunny/log-alert/releases/download/v0.1.1/log-alert-linux-amd64)
  - 上传log-alert-linux-amd64到服务器上
  - 执行 chmod 755 log-alert-linux-amd64 修改文件权限
####2.配置
- 复制config_example.yaml 到config.yaml
- 修改config.yaml
  - receivers: 通知接送者的邮箱
  - notify: 发送者邮件配置
    - driver: 发送通知驱动，mail为邮件，其他通知方式需要自行开发
    - url: 通知配置，邮件通知的格式为 发送邮箱|密码|smtp服务器|smtp端口号
  - rules: 监控规则列表
    - file: 要监控的文件
    - rule: 正则表达式
    - desc: 描述
    - duration: 警报次数累计周期时长
    - times: 一个duration期间内，达到多少次后触发警报。当duration为0时，则和周期无关，累计达到该值就触发
    - interval: 2次警报最低间隔时间
    
####3.执行
   
   可以通过nohup直接执行， nohup ./log-alert-linux-amd64 &
   
   也可以搭配pm2或其他进程管理工具使用。
    