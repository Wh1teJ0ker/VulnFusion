import React from 'react';
import { Typography, Space } from '@arco-design/web-react';
import { IconUser } from '@arco-design/web-react/icon';
import AuthButton from '../AuthButton.jsx';

export default function Header() {
    return (
        <div
            style={{
                height: 60,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                padding: '0 20px',
            }}
        >
            <Space>
                <IconUser />
                <Typography.Text>VulnFusion 漏洞平台</Typography.Text>
            </Space>

            <AuthButton />
        </div>
    );
}
