import React from 'react';
import {
    IconDashboard,
    IconBug,
    IconSettings,
    IconApps,
} from '@arco-design/web-react/icon';

const menuItems = [
    {
        key: '/dashboard',
        label: '仪表盘',
        icon: <IconDashboard />,
    },
    {
        key: '/tasks',
        label: '扫描任务',
        icon: <IconBug />,
        children: [
            {
                key: '/tasks/new',
                label: '新建任务',
                icon: <IconApps />,
            },
            {
                key: '/tasks/history',
                label: '历史记录',
                icon: <IconApps />,
            },
        ],
    },
    {
        key: '/plugins',
        label: '插件管理',
        icon: <IconSettings />,
    },
];

export default menuItems;
