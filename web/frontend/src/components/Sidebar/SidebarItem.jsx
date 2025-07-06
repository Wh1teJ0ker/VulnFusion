// components/Sidebar/SidebarItem.jsx
import React from 'react';
import { Menu } from '@arco-design/web-react';

export default function SidebarItem({ item, activePath, onClick }) {
    if (item.children && item.children.length > 0) {
        return (
            <Menu.SubMenu key={item.key} title={item.label} icon={item.icon}>
                {item.children.map((sub) => (
                    <Menu.Item
                        key={sub.key}
                        icon={sub.icon}
                        onClick={() => onClick(sub.key)}
                        style={{
                            paddingLeft: 24,
                            fontWeight: activePath === sub.key ? 'bold' : 'normal',
                        }}
                    >
                        {sub.label}
                    </Menu.Item>
                ))}
            </Menu.SubMenu>
        );
    }

    return (
        <Menu.Item
            key={item.key}
            icon={item.icon}
            onClick={() => onClick(item.key)}
            style={{
                fontWeight: activePath === item.key ? 'bold' : 'normal',
            }}
        >
            {item.label}
        </Menu.Item>
    );
}
