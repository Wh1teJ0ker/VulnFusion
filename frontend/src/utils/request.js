// src/utils/request.js
import axios from 'axios';

const instance = axios.create({
    baseURL: '/api/v1',
    timeout: 8000,
});

// 请求拦截器
instance.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token && !config.headers.Authorization) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// 响应拦截器
instance.interceptors.response.use(
    (res) => res.data,
    (err) => {
        throw err.response?.data || { message: '网络错误' };
    }
);

export default instance;
