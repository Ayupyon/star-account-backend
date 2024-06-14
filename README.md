# STAR-ACCOUNT后端

## 前置条件
1. docker
2. go v1.21及以上
3. golang-sqlc v1.26.0及以上
4. golang-migrate v4.17.1及以上
5. python3

## 使用方法

```shell
python script.py init # 拉取postgresql docker，并建立15432端口到5432端口的映射
python script.py createdb # 创建数据库
python script.py dropdb # 删除数据库
python script.py migrateup # 初始化数据库表
python script.py migratedown # 删除数据库表
python script.py sqlc # 生成数据库操作模板（一般不需要执行）
python scirpt.py server # 启动服务器
```