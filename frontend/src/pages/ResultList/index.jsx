import React, { useEffect, useState } from 'react';
import { Table, Tag, Message, Button } from '@arco-design/web-react';
import { getResultsByTaskId, getAllResults } from '../../services/result';
import { useUserStore } from '../../store/user';
import { useNavigate } from 'react-router-dom';

export default function ResultList() {
    const { role } = useUserStore();
    const [results, setResults] = useState([]);
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    // 统一格式化字段名为小写驼峰
    const formatResults = (raw) =>
        raw.map((item) => ({
            id: item.ID,
            taskId: item.TaskID,
            vulnerability: item.Vulnerability,
            severity: item.Severity,
            target: item.Target,
            timestamp: item.Timestamp,
        }));

    const fetchData = async () => {
        setLoading(true);
        try {
            const raw =
                role === 'admin'
                    ? await getAllResults()
                    : await getResultsByTaskId(localStorage.getItem('currentTaskId'));

            const formatted = formatResults(raw);
            setResults(formatted);
        } catch (err) {
            console.error('获取结果失败:', err);
            Message.error('加载漏洞结果失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, [role]);

    const columns = [
        { title: '结果 ID', dataIndex: 'id', width: 80 },
        { title: '漏洞名称', dataIndex: 'vulnerability' },
        { title: '目标地址', dataIndex: 'target' },
        {
            title: '风险等级',
            dataIndex: 'severity',
            render: (level) => {
                const colorMap = {
                    low: 'green',
                    medium: 'orange',
                    high: 'red',
                    critical: 'purple',
                };
                return <Tag color={colorMap[level] || 'gray'}>{level}</Tag>;
            },
        },
        { title: '时间戳', dataIndex: 'timestamp', width: 200 },
        {
            title: '操作',
            width: 100,
            render: (_, record) => (
                <Button size="mini" onClick={() => navigate(`/result/${record.id}`)}>
                    查看详情
                </Button>
            ),
        },
    ];

    return (
        <div style={{ padding: 24 }}>
            <h2 style={{ marginBottom: 16 }}>漏洞扫描结果</h2>
            <Table
                rowKey="id"
                columns={columns}
                data={results}
                loading={loading}
                pagination={{ pageSize: 10 }}
                border
            />
        </div>
    );
}
