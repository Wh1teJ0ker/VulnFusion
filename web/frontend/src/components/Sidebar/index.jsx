import React from 'react';
import { Menu } from '@arco-design/web-react';
import { useNavigate, useLocation } from 'react-router-dom';
import SidebarItem from './SidebarItem';
import menuItems  from './menuConfig.jsx';

export default function Sidebar() {
    const navigate = useNavigate();
    const location = useLocation();
    const currentPath = location.pathname;

    return (
        <Menu
            style={{
                width: 240,
                height: '100vh',
                borderRight: '1px solid #e5e6eb',
                background: '#ffffff',
                paddingTop: 64,
            }}
            selectedKeys={[currentPath]}
            autoOpen
        >
            {menuItems.map((item) => (
                <SidebarItem
                    key={item.key}
                    item={item}
                    activePath={currentPath}
                    onClick={navigate}
                />
            ))}
        </Menu>
    );
}
