import requests

# 获取table表对应的schema（-号分割字段）字段记录，需要传入page（当前页数，默认1）、per（每一页大小，默认30）
# GET /api/v1/:table/:schema
print("GET", requests.get(
    "http://127.0.0.1:8080/api/v1/demo_user/demo_id-demo_name?page=1&per=2").text)

# 获取table表对应的schema（-号分割字段）字段记录，约束条件是字段field值为value，需要传入page（当前页数，默认1）、per（每一页大小，默认30）
# GET /api/v1/:table/:field/:value/:schema
print("GET", requests.get(
    "http://127.0.0.1:8080/api/v1/demo_user/demo_id/3/demo_id-demo_name").text)


# 修改table表字段field为value的记录
# PUT /api/v1/:table/:field/:value
print("PUT", requests.put("http://127.0.0.1:8080/api/v1/demo_user/demo_id/3",
                          json={"demo_name": "NewName", "demo_date": "2017-02-01"}).text)

# 新增一条记录到table
# POST /api/v1/:table
print("POST", requests.post("http://127.0.0.1:8080/api/v1/demo_user",
                            json={"demo_name": "New POST", "demo_date": "1998-09-08", "demo_city_name": 1}).text)

# 删除表table中字段field为value的一条记录
# DELETE /api/v1/:table/:field/:value
print("DELETE", requests.delete(
    "http://127.0.0.1:8080/api/v1/demo_user/demo_id/3").text)

input()
