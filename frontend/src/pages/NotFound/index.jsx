import React from 'react';
import { Result, Button } from '@arco-design/web-react';
import { IconExclamationCircle } from '@arco-design/web-react/icon';
import { useNavigate } from 'react-router-dom';

export default function NotFound() {
    const navigate = useNavigate();

    return (
        <div
            style={{
                height: '100vh',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                background: '#f8f9fa',
            }}
        >
            <Result
                status="404"
                icon={<IconExclamationCircle />}
                title="页面未找到"
                subTitle="你访问的页面不存在，或链接已失效。"
                extra={
                    <Button type="primary" onClick={() => navigate('/dashboard')}>
                        返回首页
                    </Button>
                }
            />
        </div>
    );
}
