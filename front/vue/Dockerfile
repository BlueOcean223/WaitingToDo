# 使用官方 Nginx 镜像作为基础镜像
FROM nginx:alpine

# 删除默认的 nginx 配置
RUN rm /etc/nginx/conf.d/default.conf

# 复制自己的 nginx 配置文件到容器中
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 复制打包后的 dist 文件夹到 nginx 的 html 目录
COPY dist /usr/share/nginx/html

# 暴露端口
EXPOSE 80

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]
