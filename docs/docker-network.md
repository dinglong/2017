# 0. 启动两个容器，并设置网络（即只使用本地回环）
```
docker run -d --net=none --name=c1 nginx:1.13.3-alpine
docker run -d --net=none --name=c2 nginx:1.13.3-alpine
```

# 1. 构建容器网络
> 创建启动网桥、给网桥加入地址
```
ip link add br0 type bridge
ip link set br0 up
ip address add 172.18.0.1/16 dev br0
```

> 查看容器的网络状态
```
docker exec -it c1 ip link list
docker exec -it c2 ip link list
```

> 创建veth pair并启动
```
ip link add veth1 type veth peer name eth0
ip link set veth1 up
```

> 将veth1加入到网桥
```
ip link set veth1 master br0
```

> 设置eth0的网络名字空间
```
ip link set eth0 netns $(docker inspect -f {{.State.Pid}} c1)
```

> 启动容器内eth0网卡，并设置ip和默认网关
```
nsenter -n -t $(docker inspect -f {{.State.Pid}} c1) ip link set eth0 up
nsenter -n -t $(docker inspect -f {{.State.Pid}} c1) ip address add 172.18.0.2/16 dev eth0
nsenter -n -t $(docker inspect -f {{.State.Pid}} c1) ip route add default via 172.18.0.1 dev eth0
```

> 查看网桥状态
```
brctl show br0
```

> 创建veth pair并启动（为第二个容器准备）
```
ip link add veth2 type veth peer name eth0
ip link set veth2 up
ip link set veth2 master br0
ip link set eth0 netns $(docker inspect -f {{.State.Pid}} c2)
nsenter -n -t $(docker inspect -f {{.State.Pid}} c2) ip link set eth0 up
nsenter -n -t $(docker inspect -f {{.State.Pid}} c2) ip address add 172.18.0.3/16 dev eth0
nsenter -n -t $(docker inspect -f {{.State.Pid}} c2) ip route add default via 172.18.0.1 dev eth0
```

> 设置iptable的forward
```
iptables -I FORWARD -i br0 -o br0 -j ACCEPT
```

> 测试容器间通讯
```
docker exec -it c1 wget -O - 172.18.0.3
docker exec -it c2 wget -O - 172.18.0.2
```

# 2. 配置容器访问外部网络
> 允许br0和ens33（宿主机与外界通讯网卡）数据相互转发
```
iptables -I FORWARD -i br0 -o ens33 -j ACCEPT
iptables -I FORWARD -i ens33 -o br0 -j ACCEPT
```

> 由于对外的源ip是172.18.0.0/16网段地址，而外部对该地址不返回，故进行ip伪装
```
iptables -t nat -I POSTROUTING -s 172.18.0.0/16 -o ens33 -j MASQUERADE
```

> 测试外部访问
```
docker exec -it c1 wget -O - www.baidu.com
docker exec -it c2 wget -O - www.baidu.com
```

# 3. 通过hostport访问容器
> 外部访问宿主机5678端口，则转到容器地址上（执行以下命令则外部机器可以访问nginx）
```
iptables -t nat -I PREROUTING ! -i br0 -p tcp -m tcp --dport 5678 -j DNAT --to-destination 172.18.0.2:80
```

> 但是在宿主机上还不能访问，需要在OUTPUT链上增加规则
```
iptables -t nat -I OUTPUT ! -o br0 -p tcp -m tcp --dport 5678 -j DNAT --to-destination 172.18.0.2:80
```

# 4. 访问宿主机服务
> 设置veth1的hairpin模式，hairpin使veth1发送到br0的数据还可以返回给veth1。如果不设置网桥默认不会把数据发送到接收端口。同时增加了--to-ports以改变源端口
```
bridge link set dev veth1 hairpin on
```

> 到目前，容器还不能访问宿主机提供的服务（docker在宿主机启动proxy完成功能），可以使用PREROUTING和POSTROUTING替换源和目的地址实现。
```
iptables -t nat -I PREROUTING -i br0 -p tcp -m tcp --dport 5678 -m addrtype --dst-type LOCAL -j DNAT --to-destination 172.18.0.2:80
iptables -t nat -I POSTROUTING -o br0 -p tcp -m tcp --dport 80 -d 172.18.0.2 -j MASQUERADE --to-ports 1024-65535
```

> 测试（容器访问宿主机服务）
```
docker exec -it c1 wget -O - http://192.168.1.61:5678
```
