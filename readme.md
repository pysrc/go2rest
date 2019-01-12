## Mysql数据库到Restful api一键生成

简单配置一下generator.go的数据库连接信息，即可生成restful api



## 接口说明



### GET

#### 方式一



**路由：** /api/v1/:table/:schema

**说明：**获取table表对应的schema（-号分割字段）字段记录，需要传入page（当前页数，默认1）、per（每一页大小，默认30）

**举例：** /api/v1/demo_user/demo_id-demo_name?page=1&per=2



#### 方式二

**路由：** /api/v1/:table/:field/:value/:schema

**说明：** 获取table表对应的schema（-号分割字段）字段记录，约束条件是字段field值为value，需要传入page（当前页数，默认1）、per（每一页大小，默认30）

**举例：** /api/v1/demo_user/demo_id/3/demo_id-demo_name?page=1&per=2



### PUT

**路由：** /api/v1/:table/:field/:value

**说明：** 修改table表字段field为value的记录

**举例：** /api/v1/demo_user

**数据：** {"demo_name": "New POST", "demo_date": "1998-09-08", "demo_city_name": 1}



### POST

**路由：** /api/v1/:table

**说明：** 新增一条记录到table

**举例：** /api/v1/demo_user

**数据：** {"demo_name": "New POST", "demo_date": "1998-09-08", "demo_city_name": 1}



### DELETE



**路由：** /api/v1/:table/:field/:value

**说明：** 删除表table中字段field为value的一条记录

**举例：** /api/v1/demo_user/demo_id/3



