import React from 'react';
import { Routes, Route, Navigate, useLocation } from 'react-router-dom';
import Cookies from 'js-cookie';

import LoginPage from '../pages/LoginPage';
import Dashboard from '../pages/Dashboard';
import Tasks from '../pages/Tasks';
import Plugins from '../pages/Plugins';

function AppRouter() {
    const location = useLocation();
    const token = Cookies.get('auth_token');

    return (
        <Routes>
            {/* 根据是否有 token 判断重定向目标 */}
            <Route path="/" element={<Navigate to={token ? "/dashboard" : "/login"} replace />} />

            <Route path="/login" element={<LoginPage />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/tasks" element={<Tasks />} />
            <Route path="/plugins" element={<Plugins />} />
        </Routes>
    );
}

export default AppRouter;
