# Back-end

## config

存放用来读取配置文件的config.go（用viper）

## api

此文件夹用以实现api接口（也许叫docs会更好一点？）

## controllers

所有在这里的文件，负责选择调用哪个service里的函数，然后返回json数据。

## utils

此文件夹用以存放工具类，logger.go放在里头了

|文件|功能|
|---|---|
|config.go|配置文件|
|jwt.go|jwt相关|
|jsonResponse.go|格式化json输出|
|crypto.go|加密/解密，密钥放在config.yaml里|

- 后端如果有错误日志产生用zerolog搞定吧，初始化那些放在utils/logger.go里了

## middleware

此文件夹用以存放中间件

## models

此文件夹用来存放数据模型

## services

所有需要操作数据库的go文件放在这里

## routes

此文件夹用以存放路径
请放到router.go对应位置
推荐每个种api放一个文件，如user,auth，具体见apifox

## static

静态文件，如icon，css，js等

## upload

上传的文件 暂时不知道是否使用这种形式。
