# Web Management Dockerfile
FROM node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制package文件
COPY web/package*.json ./

# 安装依赖
RUN npm ci --only=production

# 复制源代码
COPY web/ .

# 构建应用
RUN npm run build

# 最终镜像
FROM nginx:alpine

# 复制构建的文件
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制nginx配置
COPY web/nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 3000

# 运行nginx
CMD ["nginx", "-g", "daemon off;"]