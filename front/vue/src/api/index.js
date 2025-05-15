// 封装aixos

import axios from 'axios';
import { ElMessage } from 'element-plus';

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
    response => response,
    error => {
        // token过期
        if (error.response && error.response.status === 401){
            ElMessage.error('登录过期，请重新登录');
        }
        return Promise.reject(error);
    }
)

export default service;