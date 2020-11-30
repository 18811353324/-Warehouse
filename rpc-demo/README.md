> ./build.sh

> docker build --no-cache --disable-content-trust=true -t docker.wanpinghui.com/rpc_demo .

> docker push docker.wanpinghui.com/rpc_demo

> docker run -d --name=test -p 8080:8080 -p 10000:10000 docker.wanpinghui.com/rpc_demo

> curl http://localhost:8080/rpc/test/rpc_call