import React from 'react';
import { Layout } from '@arco-design/web-react';
import AppHeader from './Header';
import Sidebar from './Sidebar';

const { Content } = Layout;

function BaseLayout({ children }) {
    return (
        <Layout style={{ minHeight: '100vh' }}>
            {/* 固定 Header */}
            <div style={{ position: 'fixed', top: 0, left: 0, right: 0, zIndex: 1000 }}>
                <AppHeader />
            </div>

            <Layout hasSider style={{ marginTop: 64 }}>
                {/* 固定 menuConfig.jsx */}
                <div style={{ position: 'fixed', top: 64, bottom: 0, left: 0, width: 220, zIndex: 999 }}>
                    <Sidebar />
                </div>

                {/* 主内容区域自适应宽高，左边空出 menuConfig.jsx 宽度 */}
                <Layout
                    style={{
                        marginLeft: 220,
                        padding: 24,
                        backgroundColor: '#f7f8fa',
                        height: 'calc(100vh - 64px)',
                        overflow: 'auto',
                    }}
                >
                    <Content
                        style={{
                            backgroundColor: '#ffffff',
                            padding: 24,
                            borderRadius: 12,
                            boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
                            minHeight: '100%',
                            boxSizing: 'border-box',
                        }}
                    >
                        {children}
                    </Content>
                </Layout>
            </Layout>
        </Layout>
    );
}

export default BaseLayout;
