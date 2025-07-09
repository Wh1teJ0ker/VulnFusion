import React from 'react';
import { Menu } from '@arco-design/web-react';
import {
    IconHome,
    IconList,
    IconFile,
    IconUser,
} from '@arco-design/web-react/icon';
import { useNavigate, useLocation } from 'react-router-dom';
import { useUserStore } from "../../store/user";


export default function Sidebar() {
    const navigate = useNavigate();
    const { pathname } = useLocation();
    const { role } = useUserStore(); // 获取当前角色

    const menus = [
        { key: '/dashboard', icon: <IconHome />, label: '仪表盘' },
        { key: '/tasks', icon: <IconList />, label: '扫描任务' },
        { key: '/results', icon: <IconFile />, label: '漏洞结果' },
    ];

    // 只有 admin 显示用户管理
    if (role === 'admin') {
        menus.push({ key: '/admin/users', icon: <IconUser />, label: '用户管理' });
    }

    return (
        <Menu
            selectedKeys={[pathname]}
            onClickMenuItem={(key) => navigate(key)}
            style={{ height: '100%' }}
        >
            {menus.map((item) => (
                <Menu.Item key={item.key} icon={item.icon}>
                    {item.label}
                </Menu.Item>
            ))}
        </Menu>
    );
}
