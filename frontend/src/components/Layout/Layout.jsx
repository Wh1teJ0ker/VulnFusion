import React from 'react';
import { Outlet } from 'react-router-dom';
import { Layout as ArcoLayout } from '@arco-design/web-react';
import Header from './Header.jsx';
import Sidebar from './Sidebar.jsx';

const { Sider, Header: ArcoHeader, Content } = ArcoLayout;

export default function Layout() {
    return (
        <ArcoLayout style={{ height: '100vh' }}>
            <Sider collapsible breakpoint="xl" width={200}>
                <Sidebar />
            </Sider>

            <ArcoLayout>
                <ArcoHeader>
                    <Header />
                </ArcoHeader>

                <Content style={{ padding: 20, background: '#f6f8fa' }}>
                    <Outlet />
                </Content>
            </ArcoLayout>
        </ArcoLayout>
    );
}
