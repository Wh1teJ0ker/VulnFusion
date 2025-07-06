import React from 'react';
import { Menu } from '@arco-design/web-react';
import {
    IconDashboard,
    IconBug,
    IconSettings,
} from '@arco-design/web-react/icon';
import { useNavigate, useLocation } from 'react-router-dom';

const MenuItems = [
    { key: '/dashboard', label: '仪表盘', icon: <IconDashboard /> },
    { key: '/tasks', label: '扫描任务', icon: <IconBug /> },
    { key: '/plugins', label: '插件管理', icon: <IconSettings /> },
];

function Sidebar() {
    const navigate = useNavigate();
    const location = useLocation();

    return (
        <Menu
            style={{
                width: 220,
                height: '100vh',
                borderRight: '1px solid #E5E6EB',     // gray-3
                background: '#F7F8FA',                // gray-1
                color: '#1D2129',                     // gray-10
                paddingTop: 64,                       // 留出 Header 空间
            }}
            selectedKeys={[location.pathname]}
            onClickMenuItem={key => navigate(key)}
        >
            {MenuItems.map(item => (
                <Menu.Item
                    key={item.key}
                    icon={item.icon}
                    style={{
                        paddingLeft: 20,
                        height: 50,
                        lineHeight: '50px',
                        fontSize: 15,
                        background: '#ffffff',
                        color: '#1D2129',             // 主文字色
                    }}
                >
                    {item.label}
                </Menu.Item>
            ))}
        </Menu>

    );
}

export default Sidebar;
