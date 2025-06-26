// 封装aixos

import axios from 'axios';
import { ElMessage } from 'element-plus';
import router from '@/router'

const baseURL = import.meta.env.VITE_APP_API_BASE_URL;


// 创建axios实例
const service = axios.create({
  baseURL,
});

// 请求拦截器
service.interceptors.request.use(
    config => {
        // 携带令牌
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    error => {
        return Promise.reject(error);
    }
)

// 响应拦截器
service.interceptors.response.use(
    response => {
        // 检查响应头是否有新令牌
        const newToken = response.headers['new-access-token'];
        if(newToken) {
            // 存储新令牌
            localStorage.setItem('token', newToken);
            // 更新请求头
            service.defaults.headers.common['Authorization'] = `Bearer ${newToken}`;
        }
        return response
    },
    error => {
        // token过期
        if (error.response && error.response.status === 401){
            // 清除本地信息
            localStorage.removeItem('user');
            localStorage.removeItem('token');
            // 重定向到登录页面
            router.push('/login');
            ElMessage.error('登录过期，请重新登录');
        }
        return Promise.reject(error);
    }
)

export default service;