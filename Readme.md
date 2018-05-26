![build status](https://travis-ci.org/sysu-SAAD-project/service-end.svg?branch=master)
# 服务后端说明

## 文件结构

- main.go

    代码入口文件，应只支持启动服务，而不应包含过多实现细节

- static/

    存放前端代码以及静态文件

- router/

    Store router information

- models/

    后台数据处理代码

- service/

    存放各种服务，通过调用以上两个包中的代码完成服务

- init/

    存放各种静态变量和上下文环境，用于初始化后台服务

- controller/

    Store logic codes

## 测试说明
