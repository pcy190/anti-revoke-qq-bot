QQ防撤回bot
===

- 基于 [IOTQQ](https://github.com/IOTQQ/IOTQQ)

# 功能

~~开源版目前集成了防撤回和权限控制~~。

防撤回支持文本消息、图片、xml等消息格式。

权限控制中,可以添加允许控制撤回的qq号, 并将配置文件conf写入程序运行目录。

# 如何使用

## 环境配置
- 遵照[IOTQQ](https://github.com/IOTQQ/IOTQQ)配置好服务端，并确保安装好mysql
- 修改`model/config.go`文件中的配置
- `go run main.go`

第一次使用时，如果不想手动创建数据库的表,可以使用`util.CreateTable()`,这在main.go中默认启用，后续则可以注释这条。

## 防撤回
进入群聊，发送`!revoke on`即可开启对该群聊的防撤回功能。

如果要关闭，则发送`!revoke off`即可。
管理员可以私聊机器人发送`!status`即可查看当前开启防撤回的群聊与已获取权限的qq号。

# 目录结构
`handler`目录
- `EventMsg.go` : 事件处理回调，包含了撤回事件的回调处理
- `FriendMsg.go` : 私聊消息回调
- `GroupMsg.go` : 群聊消息回调

# 拓展
- go语言对数据库的处理非常灵活，如果要适配其他版本数据库，重写`util/util.go`中的`AddRecord`与`QueryRecord`方法即可。



