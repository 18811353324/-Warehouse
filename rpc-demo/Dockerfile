FROM ccr.ccs.tencentyun.com/wanpinghui/alpine

ENV     GIN_MODE release
RUN     mkdir -p /opt/rpc_demo /log/rpc_demo
WORKDIR /opt/rpc_demo
COPY    ./bin/ ./
RUN     chmod +x ./rpc_demo
EXPOSE 8080
EXPOSE 10000
ENV PATH $PATH:/sbin:/opt/rpc_demo

ENTRYPOINT ["tini", "--"]
CMD ["rpc_demo", "-c", "config", ">>", "/log/rpc_demo/log.log", "2>&1", "&"]