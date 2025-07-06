// src/pages/LoginPage/index.jsx
import React from 'react';
import useLogin from './useLogin';
import LoginForm from './LoginForm';

export default function LoginPage() {
    const login = useLogin();

    return (
        <div
            style={{
                position: 'relative',
                width: '100vw',
                height: '100vh',
                overflow: 'hidden',
            }}
        >

            {/* 登录表单（居中） */}
            <div
                style={{
                    position: 'absolute',
                    top: '50%',
                    left: '50%',
                    transform: 'translate(-50%, -50%)',
                    zIndex: 1,
                }}
            >
                <LoginForm {...login} />
            </div>
        </div>
    );
}
