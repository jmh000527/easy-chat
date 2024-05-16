#!/bin/bash
reso_addr='registry.cn-hangzhou.aliyuncs.com/easy-chat-jmh/social-rpc-dev'
tag='latest'

pod_ip="192.168.199.138"

container_name="easy-chat-social-rpc-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_easy-im -v /easy-im/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -p 10001:10001 -e POD_IP=${pod_ip}  --name=${container_name} -d ${reso_addr}:${tag}