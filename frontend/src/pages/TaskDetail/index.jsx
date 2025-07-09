import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Descriptions,
    Button,
    Message,
    Tag,
    Space,
    Typography,
    Card,
    Spin,
} from '@arco-design/web-react';
import { getTaskDetail } from '../../services/task';

const { Title } = Typography;

export default function TaskDetail() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [task, setTask] = useState(null);
    const [loading, setLoading] = useState(true);

    const statusColors = {
        pending: 'blue',
        running: 'orange',
        done: 'green',
        failed: 'red',
    };

    const fetchTask = async () => {
        setLoading(true);
        try {
            const res = await getTaskDetail(id);
            if (res && res.ID) {
                setTask(res); // 直接使用原始字段，避免大小写映射出错
            } else {
                Message.error('未获取到任务数据');
            }
        } catch (err) {
            Message.error('加载任务详情失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchTask();
    }, [id]);

    return (
        <div style={{ padding: '24px' }}>
            <Space direction="vertical" size={16} style={{ width: '100%' }}>
                <Button onClick={() => navigate(-1)} type="outline">
                    ← 返回任务列表
                </Button>

                <Card bordered hoverable style={{ maxWidth: 700 }}>
                    <Title heading={5}>任务详情</Title>

                    {loading ? (
                        <Spin style={{ marginTop: 40 }} tip="加载中..." />
                    ) : task ? (
                        <Descriptions
                            column={1}
                            layout="horizontal"
                            style={{ marginTop: 20 }}
                        >
                            <Descriptions.Item label="任务 ID">{task.ID}</Descriptions.Item>
                            <Descriptions.Item label="目标地址">{task.Target}</Descriptions.Item>
                            <Descriptions.Item label="模板名称">{task.Template}</Descriptions.Item>
                            <Descriptions.Item label="状态">
                                <Tag color={statusColors[task.Status] || 'gray'}>
                                    {task.Status}
                                </Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="创建时间">{task.CreatedAt}</Descriptions.Item>
                            <Descriptions.Item label="所属用户 ID">{task.UserID}</Descriptions.Item>
                        </Descriptions>
                    ) : (
                        <div style={{ marginTop: 40, color: '#999' }}>暂无任务信息</div>
                    )}
                </Card>
            </Space>
        </div>
    );
}
