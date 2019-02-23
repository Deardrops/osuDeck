# 介绍

怎么说呢，有人需要，我就做了这么个工具。

![Preview](https://ww1.sinaimg.cn/large/007i4MEmgy1g02dgnvnlgj30bi0efjs7.jpg)

读取本地的谱面和收藏夹数据，与别人分享给你的文件（osu!.db 或 collection.db）做比较，批量下载所有缺失的谱面。

有三个谱面下载源可以使用，分别是官方，血猫和 sayobot。

目前正在重构整个项目，去掉 GUI 库，改为 服务端程序+Web前端 的方式进行交互。

具体请参阅 [开发记录](https://github.com/Deardrops/osuDeck/issues/1)。

# 使用指南

第一步，在[腾讯微云](https://share.weiyun.com/5KMHVRY)下载最新的版本。

第二步，下载到本地后，找一个剩余空间比较多的磁盘（比如E盘），新建一个文件夹（可以命名为 osuDeck），将下载好的 exe 文件**剪贴**到这个文件夹中。

第三步，双击运行程序，第一行是以 `Open` 开始三个按钮，选择你的本地 osu! 文件夹，这时你应该可以看到方框中显示了本地的谱面数量等信息。

第四步

A：**osu!.db** 选择别人分享给你的 `osu!db` 文件，然后点击第二行中的 `osu!db` 标签页，你可以看到导入的 db 文件的相关信息，
点击 `download missed beatmaps`，即可开始下载。
下载好的谱面包在程序目录下的 `download` 文件夹中。

B：**collection.db** 选择别人分享给你的 `collection.db` 文件，然后点击第二行中的 `colleciton` 标签页，你可以看到导入的收藏夹的相关信息，
先点击 `load beatmaps in collection`，这一步会向官网API发送请求，解析这个文件中的谱面。（注意，**必须设置自己的 API_KEY**，详见下）。
加载完成后，点击 `download missed beatmaps`，即可开始下载。

> 提示：osu!Deck 默认使用的镜像是 Bloodcat，下载的谱面不带视频和故事版。

## 配置文件

配置文件名称为 **`conf.yaml`** ，请使用“专业”的文本编辑器打开，例如[Sublime Text](https://www.sublimetext.com/)，[Atom](https://atom.io/)，记事本不行。

配置文件会在第一次运行程序后自动创建，**修改配置文件时请关闭程序**，修改完成后记得保存。

这里有一份[配置文件的示例](https://github.com/Deardrops/osuDeck/blob/master/example.conf.yaml)。

## 设置 API_KEY

在[官网上](https://osu.ppy.sh/p/api)申请自己的 API key，然后将 API key 字符串复制到配置文件中 `osu_api_key` 字段中。

更多关于 osu!api 的介绍，请参阅[官方WIKI](https://github.com/ppy/osu-api/wiki)。

## 切换镜像

如果遇到下载不动图的情况，请手动切换镜像后再尝试。

在配置文件中找到 mirror 开头的那一行，将后面的值改成以下三个镜像中的一个。（建议直接复制过去，以防出错）

- `official` - osu! 官方下载，需要在配置文件中输入自己的账号密码
- `bloodcat` - 血猫，国内会有不错的下载速度
- `sayobot` - 小夜的国内谱面镜像网站

> 提示: 从官网下图需要在配置文件中输入自己的账号密码才能下载。

## 代理设置

程序支持代理，不过需要手动设置，程序使用环境变量 `HTTPS_PROXY` 作为代理设定。

# 反馈问题

在使用过程中遇到任何问题，如崩溃啦，下载不动啦，可以通过[开 Issue](https://github.com/Deardrops/osuDeck/issues/new) 的方式反馈给我。

反馈问题时可以附上最新的日志文件（在程序目录下的 logs 文件夹中）。
