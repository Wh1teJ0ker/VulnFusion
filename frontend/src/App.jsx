// src/App.jsx
import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import LoginPage from './pages/LoginPage';          // 登录页
import Dashboard from './pages/Dashboard';          // 仪表盘（待开发）
import NotFound from './pages/NotFound';            // 404 页面
import Layout from './components/Layout/Layout';    // 通用布局（侧边栏 + Header）

import TaskList from './pages/TaskList';
import TaskDetail from './pages/TaskDetail';
import ResultList from './pages/ResultList';
import ResultDetail from './pages/ResultDetail';
import UserManagement from './pages/UserManagement';


import { useUserStore } from './store/user';         // 用户状态管理

export default function App() {
    const { token } = useUserStore();

    return (
        <Router>
            <Routes>
                <Route path="/login" element={<LoginPage />} />
                <Route
                    path="/"
                    element={
                        token ? (
                            <Layout />
                        ) : (
                            <Navigate to="/login" replace />
                        )
                    }
                >
                    {/* ✅ 添加默认重定向 */}
                    <Route index element={<Navigate to="/dashboard" replace />} />

                    <Route path="dashboard" element={<Dashboard />} />
                    <Route path="tasks" element={<TaskList />} />
                    <Route path="task/:id" element={<TaskDetail />} />
                    <Route path="results" element={<ResultList />} />
                    <Route path="result/:id" element={<ResultDetail />} />
                    <Route path="admin/users" element={<UserManagement />} />
                </Route>


                <Route path="*" element={<NotFound />} />
            </Routes>
        </Router>
    );
}
