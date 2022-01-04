## 1、中间件部署约定
中间件部署统一采用docker方式，IP为固定的IP，例如：

| 中间件 | ip地址 |
| --- | --- |
| MySQL | 172.100.0.100 |
| zookeeper | 172.100.0.101 |
| kafka | 172.100.0.102 |
| canal | 172.100.0.103 |



## 
## 2、相关Docker 命令
```
//创建自定义网络，设定固定IP,便于部署
sudo docker network create --subnet=172.100.0.0/24 bigdata-network

//查看各容器的ip地址
sudo docker inspect -f '{{.Name}} - {{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(sudo docker ps -aq)
```

<br />
<br />

## 2、MySQL安装部署
```
// pull docker image
docker pull mysql:5.7


//映射mysql配置文件到宿主机
docker run -d --net bigdata-network --ip 172.100.0.100 --name mysql5.7 -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7
mkdir -p ~/docker/mysql/5.7/data
mkdir -p ~/docker/mysql/5.7/conf
docker cp mysql5.7:/var/lib/mysql  ~/docker/mysql/5.7/data
docker cp mysql5.7:/etc/mysql  ~/docker/mysql/5.7/conf

//rm mysql containor
docker stop mysql5.7
docker rm mysql5.7

//create mysql containor with host dir
sudo docker run  -d --net bigdata-network --ip 172.100.0.100 --name mysql5.7 -p 3306:3306  -v ~/docker/mysql/5.7/data/mysql:/var/lib/mysql  -v ~/docker/mysql/5.7/conf/mysql:/etc/mysql  -e MYSQL_ROOT_PASSWORD=123456  mysql:5.7

//启动binlog，添加以下内容到： ~/docker/mysql/5.7/conf/mysql/mysql.conf.d/mysqld.cnf
log-bin=/var/lib/mysql/mysql-bin
server-id=123454
binlog-format=ROW

//restart mysql containor
docker restart mysql5.7

//check binlog config
mysql -h 127.0.0.1 -uroot -p123456
mysql> SHOW  GLOBAL VARIABLES LIKE '%log_bin%';
+---------------------------------+--------------------------------+
| Variable_name                   | Value                          |
+---------------------------------+--------------------------------+
| log_bin                         | ON                             |
| log_bin_basename                | /var/lib/mysql/mysql-bin       |
| log_bin_index                   | /var/lib/mysql/mysql-bin.index |
| log_bin_trust_function_creators | OFF                            |
| log_bin_use_v1_row_events       | OFF                            |
+---------------------------------+--------------------------------+

//create canal account and auth
CREATE USER 'canal'@'%' IDENTIFIED BY 'canal';
ALTER USER 'canal'@'%' IDENTIFIED WITH mysql_native_password BY 'canal';
grant all privileges on *.* to 'canal'@'%' identified by 'canal' with grant option;
FLUSH PRIVILEGES;

```


## 3、Kafka安装部署
```
//安装zookeeper
sudo docker run -d --net bigdata-network --ip 172.100.0.101 --name zookeeper -p 2181:2181  wurstmeister/zookeeper:latest

//安装配置kafka
docker run  -d --net bigdata-network --ip 172.100.0.102 --name kafka -p 9092:9092 -e KAFKA_BROKER_ID=0 -e KAFKA_ZOOKEEPER_CONNECT=172.100.0.101:2181 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://172.100.0.102:9092 -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 -t wurstmeister/kafka:latest

//清空单个topic消息
kafka-topics.sh --zookeeper 172.18.100.101:2181 --delete --topic stock_stock_info
kafka-topics.sh --zookeeper 172.18.100.101:2181 --delete --topic stock_stock_quotation
```
## 4、canal安装部署
```
sudo docker pull canal/canal-server:v1.1.5

mkdir -p  ~/docker/canal-server/v115/instance
mkdir -p  ~/docker/canal-server/v115/conf
sudo docker run -p 11111:11111 --name canal -d canal/canal-server:v1.1.5
sudo docker cp canal:/home/admin/canal-server/conf/example/instance.properties ~/docker/canal-server/v115/instance
sudo docker cp canal:/home/admin/canal-server/conf/canal.properties ~/docker/canal-server/v115/conf

sudo docker stop canal
sudo docker rm canal


sudo docker run -p 11111:11111 --name canal -v ~/docker/canal-server/v115/instance/instance.properties:/home/admin/canal-server/conf/example/instance.properties -v ~/docker/canal-server/v115/conf/canal.properties:/home/admin/canal-server/conf/canal.properties -d --net bigdata-network --ip 172.100.0.103 canal/canal-server:v1.1.5

```


## 5、canal配置
```
vim  ~/docker/canal-server/v115/instance/instance.properties 修改为以下配置：

canal.instance.master.address=172.100.0.100:3306
canal.mq.dynamicTopic=stock\\.stock_info,stock\\.stock_quotation


vim ~/docker/canal-server/v115/conf/canal.properties 修改为以下配置：
canal.serverMode = kafka
kafka.bootstrap.servers = 172.100.0.102:9092

重启生效
docker restart canal

```
