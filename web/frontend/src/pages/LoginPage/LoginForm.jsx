// src/pages/LoginPage/LoginForm.jsx
import React from 'react';
import { Form, Input, Button, Typography } from '@arco-design/web-react';
import { IconUser, IconLock } from '@arco-design/web-react/icon';
import './LoginForm.css';

const { Title, Paragraph, Text } = Typography;

export default function LoginForm({ username, password, loading, setUsername, setPassword, handleLogin, handleKeyPress }) {
    return (
        <div className="login-container">
            <div className="login-panel">
                <div className="login-header">
                    <Title heading={4} style={{ marginBottom: 4 }}>🔐 VulnFusion</Title>
                    <Paragraph style={{ fontSize: 14, color: '#888' }}>
                        插件驱动 · 自动化漏洞检测平台
                    </Paragraph>
                </div>

                <Form layout="vertical" onSubmit={handleLogin}>
                    <Form.Item label="用户名" required>
                        <Input
                            prefix={<IconUser />}
                            value={username}
                            onChange={setUsername}
                            onKeyPress={handleKeyPress}
                            placeholder="请输入用户名"
                            allowClear
                        />
                    </Form.Item>

                    <Form.Item label="密码" required>
                        <Input.Password
                            prefix={<IconLock />}
                            value={password}
                            onChange={setPassword}
                            onKeyPress={handleKeyPress}
                            placeholder="请输入密码"
                            allowClear
                        />
                    </Form.Item>

                    <Form.Item>
                        <Button
                            type="primary"
                            htmlType="submit"
                            long
                            loading={loading}
                        >
                            登录
                        </Button>
                    </Form.Item>
                </Form>

                <div className="login-footer">
                    <Text type="secondary" style={{ fontSize: 12 }}>
                        © Wh1teJ0ker
                    </Text>
                </div>
            </div>
        </div>
    );
}
