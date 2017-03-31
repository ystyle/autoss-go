### autoss-go

### 下载autoss
  >下载地址: [https://github.com/ystyle/autoss-go/releases](https://github.com/ystyle/autoss-go/releases)

### 修改配置文件
本软件依赖[shadowsocks](http://www.ishadowsocks.com/)请自行下载,并手动配置shadowsocks的位置.

config.json 文件

配置说明:

1. `cmd`  shadowsocks.exe 的位置

2. `json` shadowsocks的配置文件gui-config.json

  > (如果没有就先打开shadowsocks,然后随便添加一个服务器保存就有了)

3. `timeout` ss的超时时间

4. `local_port` 本地代理端口

5. `args` 程序参数

### 使用
把`autoss`发送到桌面快捷方式

可以右键属性把图标改成 shadowsocks.exe

双击桌面的图标启动shadowsocks , 在启动时会自动获取最新的服务器和密码

### Chrome浏览器设置
下载插件: [SwitchyOmega](https://github.com/FelisCatus/SwitchyOmega/releases)

插件配置备份:
[OmegaOptions.zip](https://github.com/ystyle/autoss-go/files/528625/OmegaOptions.zip)
装好插件直接导入备份
