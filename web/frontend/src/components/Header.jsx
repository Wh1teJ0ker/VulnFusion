import React from 'react';
import { Layout, Typography, Button, Avatar, Space } from '@arco-design/web-react';
import { IconExport, IconUser } from '@arco-design/web-react/icon';
import Cookies from 'js-cookie';
import { useNavigate } from 'react-router-dom';

const { Header } = Layout;
const { Title } = Typography;

function AppHeader() {
    const navigate = useNavigate();

    const handleLogout = () => {
        Cookies.remove('auth_token');
        navigate('/login');
    };

    return (
        <Header
            style={{
                height: 64,
                background: '#ffffff',
                borderBottom: '1px solid #e5e6eb',
                display: 'flex',
                alignItems: 'center',
                padding: '0 24px',
                boxShadow: '0 2px 8px rgba(0,0,0,0.04)',
            }}
        >
            {/* 左侧平台 logo 与标题 */}
            <Space align="center">
                <img
                    src="/vite.svg"
                    alt="logo"
                    style={{ height: 36, marginRight: 12 }}
                />
                <Title heading={5} style={{ color: '#00c8c8', margin: 0 }}>
                    VulnFusion 安全平台
                </Title>
            </Space>

            {/* 右侧用户操作 */}
            <Space size={16}>
                <Avatar style={{ backgroundColor: '#00c8c8' }} icon={<IconUser />} />
                <Button
                    type="text"
                    icon={<IconExport />}
                    onClick={handleLogout}
                    style={{ color: '#ffffff' }}
                >
                    退出登录
                </Button>
            </Space>
        </Header>
    );
}

export default AppHeader;
