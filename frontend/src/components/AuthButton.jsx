import React from 'react';
import { Button, Dropdown, Menu } from '@arco-design/web-react';
import { IconDown } from '@arco-design/web-react/icon';
import { useUserStore } from '../store/user';
import { useNavigate } from 'react-router-dom';

export default function AuthButton() {
    const { username, clearUser } = useUserStore();
    const navigate = useNavigate();

    const handleLogout = () => {
        clearUser();  // 清空 token
        navigate('/login');
    };

    const menu = (
        <Menu>
            <Menu.Item onClick={handleLogout}>退出登录</Menu.Item>
        </Menu>
    );

    return (
        <Dropdown droplist={menu} trigger="click" position="bl">
            <Button type="primary">
                {username || '未登录'}
                <IconDown />
            </Button>
        </Dropdown>
    );
}
