import React from 'react';
import { Spin } from '@arco-design/web-react';

export default function LoadingSpinner({ tip = '加载中...' }) {
    return (
        <div
            style={{
                minHeight: 300,
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
            }}
        >
            <Spin tip={tip} />
        </div>
    );
}
