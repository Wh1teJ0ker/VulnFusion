// src/pages/LoginPage/index.jsx
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Card,
    Form,
    Input,
    Button,
    Message,
    Typography,
} from '@arco-design/web-react';
import { IconUser, IconLock } from '@arco-design/web-react/icon';
import ReactCanvasNest from 'react-canvas-nest';

import { login } from '../../services/auth';
import { getUserInfo } from '../../services/user';
import { useUserStore } from '../../store/user';

export default function LoginPage() {
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();
    const setUser = useUserStore((state) => state.setUser);

    const handleSubmit = async (values) => {
        setLoading(true);
        try {
            const { token } = await login(values);
            if (!token) {
                Message.error('登录失败：未返回 token');
                return;
            }

            localStorage.setItem('token', token);
            const userInfo = await getUserInfo(token);
            setUser({ token, ...userInfo });

            Message.success('登录成功');
            navigate('/dashboard');
        } catch (err) {
            Message.error(err?.message || '登录失败');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div
            style={{
                height: '100vh',
                width: '100%',
                position: 'relative',
                overflow: 'hidden',
                background: '#e6f0ff',
            }}
        >
            {/* 背景粒子特效 */}
            <ReactCanvasNest
                className="canvasNest"
                config={{
                    pointColor: '66, 90, 128',   // 柔和蓝色线条
                    lineColor: '66, 90, 128',
                    pointOpacity: 0.5,
                    pointR: 2,
                    count: 120,
                }}
                style={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    zIndex: 1,
                    width: '100%',
                    height: '100%',
                }}
            />

            {/* 登录卡片 */}
            <div
                style={{
                    position: 'relative',
                    zIndex: 10,
                    display: 'flex',
                    height: '100%',
                    alignItems: 'center',
                    justifyContent: 'center',
                }}
            >
                <Card style={{ width: 380 }}>
                    <Typography.Title heading={5}>VulnFusion 登录</Typography.Title>

                    <Form layout="vertical" onSubmit={handleSubmit}>
                        <Form.Item label="用户名" field="username" rules={[{ required: true }]}>
                            <Input prefix={<IconUser />} placeholder="请输入用户名" />
                        </Form.Item>

                        <Form.Item label="密码" field="password" rules={[{ required: true }]}>
                            <Input.Password prefix={<IconLock />} placeholder="请输入密码" />
                        </Form.Item>

                        <Form.Item>
                            <Button type="primary" htmlType="submit" loading={loading} long>
                                登录
                            </Button>
                        </Form.Item>
                    </Form>
                </Card>
            </div>
        </div>
    );
}
