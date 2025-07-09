import React, { useEffect, useState } from 'react';
import { Table, Button, Message, Tag } from '@arco-design/web-react';
import { useUserStore } from '../../store/user';
import { getTasks, getAllTasks } from '../../services/task';
import { useNavigate } from 'react-router-dom';

export default function TaskList() {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(false);
    const { role } = useUserStore();
    const navigate = useNavigate();

    const fetchTasks = async () => {
        setLoading(true);
        try {
            const res = role === 'admin' ? await getAllTasks() : await getTasks();

            // 字段规范化，转换为小写 key
            const normalized = (res || []).map(item => ({
                id: item.ID,
                target: item.Target,
                template: item.Template,
                status: item.Status,
                userId: item.UserID,
                createdAt: item.CreatedAt,
            }));

            setData(normalized);
        } catch (err) {
            Message.error('加载任务失败');
        } finally {
            setLoading(false);
        }
    };


    useEffect(() => {
        fetchTasks();
    }, [role]);

    const columns = [
        { title: '任务 ID', dataIndex: 'id', width: 80 },
        { title: '目标地址', dataIndex: 'target' },
        { title: '模板路径', dataIndex: 'template' },
        {
            title: '状态',
            dataIndex: 'status',
            render: (status) => {
                const colorMap = {
                    pending: 'blue',
                    running: 'orange',
                    done: 'green',
                    failed: 'red',
                };
                return <Tag color={colorMap[status] || 'gray'}>{status}</Tag>;
            },
        },
        {
            title: '操作',
            render: (_, record) => (
                <Button
                    size="mini"
                    onClick={() => navigate(`/task/${record.id}`)}
                >
                    查看详情
                </Button>
            ),
        },
    ];

    return (
        <div>
            <h2 style={{ marginBottom: 16 }}>任务列表</h2>
            <Table
                rowKey="id"
                columns={columns}
                data={data}
                loading={loading}
                pagination={{ pageSize: 10 }}
            />
        </div>
    );
}
