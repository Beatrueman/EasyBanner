# EasyBanner

## 功能介绍

**通过`@机器人`来触发飞书机器人对特定机器远程封禁恶意IP，依赖于[GitHub - evilsp/xdp_banner: 一个简单的 XDP 小程序，用于 BAN IP](https://github.com/evilsp/xdp_banner)**

1.当机器人未连接服务器时，提示如下消息。

![image-20240827174105360](https://gitee.com/beatrueman/images/raw/master/img/202408271741464.png)

2.当未发现访问次数超过250次的IP，提示如下消息。

![image-20240827174159408](https://gitee.com/beatrueman/images/raw/master/img/202408271741444.png)

3.当发现访问次数超过250次的恶意IP，罗列出恶意IP以及对应的访问次数，并显示封禁BAN按钮。

![image-20240827174315387](https://gitee.com/beatrueman/images/raw/master/img/202408271743440.png)

点击按钮，可以对远程机器进行封禁恶意IP，然后卡片更新，提示封禁完成。

![image-20240827174403213](https://gitee.com/beatrueman/images/raw/master/img/202408271744269.png)

## 逻辑介绍

### 接口

目标机器开放两个接口

- GET /execute：用于查询日志里当前小时内访问量排名前十的IP以及对应的次数。
- POST /ban：用于接收需要封禁的IP，然后执行xdp封禁命令。

机器人设置两个接口

- POST /webhook：用来处理接收消息事件。
- POST /event：用来处理卡片回调

### 运行逻辑

用户给机器人发消息，触发`接收消息v2.0`事件，飞书服务器返回消息体，通过消息体里`messageBody.Event.Message.Mentions`的`mention.ID.UserID`是否为空来判断用户发的消息是否为@机器人，不是则忽略。

是则发送卡片消息。这时调用目标机器`/execute`接口来获得IP以及对应的次数，通过判断 ip 次数，发送不同的模板。

没有返回数据，则发送未连接服务器模板。有返回数据但没有大于250次的IP，返回未检测到恶意IP模板。

有返回数据且有大于250次的IP时，获取所有次数大于250次的IP以及对应的次数动态填充至 JSON 模板，然后发送。

此时用户点击红色按钮BAN，会触发`卡片回传交互`事件，此时会回传数据。接下来对回传消息体进行一些判断：

1. 判断`event_type`是否为`card.action.trigger`
2. 判断`action.Tag`是否为`button`
3. 判断`action.Value`中键`action`对应的值是否为`ban_ip`

如果全部满足，则将需要封禁的IP制作成请求体，对目标机器`POST /ban`进行调用。

API成功调用后，调用飞书`更新卡片`API，对卡片内容进行更新。

## 部署

### 申请机器人

#### *app_id*与*app_secret*获取方法

1.用企业账户，在开发者后台中，**创建企业自建应用**

![image-20230731184331697](https://gitee.com/beatrueman/images/raw/master/img/202307311843815.png)

2.找到app_id与qpp_secret

![image-20230731184507412](https://gitee.com/beatrueman/images/raw/master/img/202408272029722.png)

3.添加应用能力，选择机器人

![image-20230731184549009](https://gitee.com/beatrueman/images/raw/master/img/202408272029399.png)

4.添加以下权限

```
im:message,im:message.group_at_msg,im:message.group_at_msg:readonly,im:message.group_msg,im:message.p2p_msg,im:message.p2p_msg:readonly,im:message:readonly,im:chat:readonly,im:chat,im:message:send_as_bot
```

![image-20230731184637236](https://gitee.com/beatrueman/images/raw/master/img/202408272029978.png)

5.订阅**接收消息**和**卡片回传交互**

若要使机器人有互动对话功能，需要填写请求配置地址，并添加**接收消息v2.0**和**消息已读v2.0**事件

卡片交互需要订阅**卡片回传交互**

![image-20240827203202358](https://gitee.com/beatrueman/images/raw/master/img/202408272032499.png)

![image-20240827203244721](https://gitee.com/beatrueman/images/raw/master/img/202408272032804.png)

### 裸机部署

首先需要在目标机器上执行`EasyBanner/pkgs/data/app_current.py`，保持其稳定运行。

最好将其制作成Service，保证后台持久运行。

这里提供`get_ip.service`文件供参考。

```
[Unit] 
Description=Get IP Service 
After=network.target 
[Service] 
User=root 
WorkingDirectory=/root/yiiong/get_ip  # app_current.py所在目录
Environment="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
ExecStart=/bin/bash -c 'source /root/yiiong/get_ip/venv/bin/activate && exec python3 app_current.py'  # 运行Python程序，注意文件路径
Restart=always 
[Install] 
WantedBy=multi-user.target
```

或者手动运行

```
python3 ../EasyBanner/pkgs/data/app_current.py
```

然后下载依赖后，运行`main.go`

```
go mod tidy
go run main.go
```

### Docker部署

目标机器的API只能裸机部署

容器启动

```
docker run -d -p 8080:8080 \
-e APP_ID=cli_a42011585561100d \ # 填写飞书应用 AppID
-e APP_SECRET=AYRFbDlUH8OKxRweuXM47cLLFwRpO12X \ # 填写飞书应用 AppSecret
-e URL=http://172.20.14.2:9521 \  # 填写检测主机地址
-e GIN_MODE=release beatrueman/easybanner:stable
```

### Kubernetes部署

先创建`Secret`设置环境变量

```
kubectl create secret generic easybanner-secret \
--from-literal=App_ID=your_AddID \
--from-literal=App_Secret=your_AppSecret \
--from-literal=URL=your_url \
--from-literal=GIN_MODE=release \ # 设置gin为生产模式
--namespace=your_namespace
```



然后`apply`../deply目录下的`deployment.yaml`，`service.yaml`。

所有注意修改命名空间。

还有要**注意部署的机器人与目标机器接口的通信问题**

```
kubectl apply deployment.yaml
kuectl apply service.yaml

# 如果有域名需求，还可以添加ingressroute或者ingress，注意修改host
# 这里使用ingressroute
kubectl apply -f ingressroute.yaml
```

## 缺陷

- 机器人可能会重复发送消息。
- 因为会对大日志文件利用bash进行查询，所以执行`GET /execute`速度会比较慢，测试大概会消耗20s左右。
- BAN操作比较慢，需要等待一会才能更新卡片。

## 更新日志

### 2024/9/15

添加了用于打包和推送镜像的**Github Action**

使用时在Settings >> Secrets and varibles >> Actions中添加secrets 

`REGISTRY_USERNAME`和`REGISTRY_PASSWORD`

![image-20240915002948963](https://gitee.com/beatrueman/images/raw/master/img/202409150029068.png)

如果要推送到类似Harbor的自建仓库，请添加varibles

`IMAGE_REGISTRY_SERVICE`：默认为docker.io

`IMAGE_REPOSITORY`：默认为beatrueman/eaesybanner

推送时请指定**tag**，格式为`v1.0.0`，用于指定镜像版本

```
git tag v1.0.0
git push origin v1.0.0
```

或者手动指定**tag**

### 2024/11/23

为了实现机器人在无人值守的情况下也能够自动封禁IP，现在增加了自动封禁IP功能。

结合了`Grafana Alerting`，机器人能够做到Grafana在alert后，自动在飞书群中发送恶意IP列表，并且自动执行封禁IP。

#### 新增了模拟Grafana Alerting发送报警的过程

在`EasyBanner/utils/mock`中新增`mock.go`，用于发送`templates/alert.json`，其中`alert.json`中包含了`Grafana Alerting Webhook`的消息体，详细内容可在仓库或[配置用于警报的 Webhook 通知程序 | Grafana 文档 - Grafana 可观测平台](https://grafana.org.cn/docs/grafana/latest/alerting/configure-notifications/manage-contact-points/integrations/webhook-notifier/)查看。方便了开发过程中的对EasyBanner的测试。

#### 新增了一些获取基本信息的函数

在`EasyBanner/utils/base`下新增`info.go`和`InitConfig.go`

- `info.go`

  ```
  func GetTenantAccessToken() string {}                      // 获取机器人的tenant_access_token
  
  func GetChatID() string {}                                 // 获取机器人所在群聊 chat_id
  
  func GetLatestMessageID(AppID string) (string, error) {}   // 获取群组来自EB的最新一个卡片的message_id
  
  ```

- `InitConfig.go`：将原先的initConfig()封装为包

#### 新增根据Grafana Alerting直接推送恶意IP列表并直接封禁的逻辑

`EasyBanner/pkgs/grafana/receive.go`：用于接收来自Grafana联络点的消息体，提取`Alerts.labels.alertname`。其中alertname有预设的两种：

- aviation telecom带宽持续满载
- 集群内部流量

当`alertname`不为空时，会发送无按键的交互卡片消息，等待一段时间后，机器人会调用远程机器的封禁/ban接口，然后对IP封禁，最后更新卡片，响应封禁是否成功。

新增的函数

```
// 用于不需要按键交互的生成模板
func GenerateNoButtonTemplate(ipDataList []model.IPData, showButton bool, resultText string) (string, error) {}

// 获取所有次数大于250次的IP以及对应的次数填充至 JSON 模板，无按键版
func GetNeedBanIPNoButton() (string, error) {}

// 获取alertname，输出报警类型
func GetAlertName(body *model.Body) (string, error) {}

// 无按键版卡片交互
func HandleAlert(c *gin.Context) {}
```

