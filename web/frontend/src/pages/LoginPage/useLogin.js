// src/pages/LoginPage/useLogin.js
import { useState } from 'react';
import axios from 'axios';
import { Message } from '@arco-design/web-react';
import { useNavigate } from 'react-router-dom';

export default function useLogin() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    const handleLogin = async () => {
        if (!username || !password) {
            Message.warning('请输入用户名和密码');
            return;
        }

        setLoading(true);
        const hideLoading = Message.loading({ content: '正在登录中...', duration: 0 });

        try {
            await axios.post('/api/v1/auth/login', { username, password }, { withCredentials: true });
            hideLoading();
            Message.success('登录成功，正在跳转...');
            setTimeout(() => navigate('/dashboard'), 500);
        } catch (err) {
            hideLoading();
            const msg =
                (err.response?.data?.error) ||
                err.message ||
                '登录失败，请重试';
            Message.error(`登录失败：${msg}`);
        } finally {
            setLoading(false);
        }
    };

    return {
        username,
        setUsername,
        password,
        setPassword,
        loading,
        handleLogin,
    };
}
