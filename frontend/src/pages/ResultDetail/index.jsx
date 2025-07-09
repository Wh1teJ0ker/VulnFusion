import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Descriptions,
    Typography,
    Button,
    Tag,
    Message,
} from '@arco-design/web-react';
import LoadingSpinner from '../../components/LoadingSpinner';
import { getResultDetail } from '../../services/result';

export default function ResultDetail() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [result, setResult] = useState(null);
    const [loading, setLoading] = useState(false);

    const fetchData = async () => {
        setLoading(true);
        try {
            const res = await getResultDetail(id);
            setResult(res);
        } catch (err) {
            Message.error('加载结果详情失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, [id]);

    const severityColor = {
        low: 'green',
        medium: 'orange',
        high: 'red',
        critical: 'purple',
    };

    if (loading || !result) return <LoadingSpinner tip="加载结果详情中..." />;

    return (
        <div>
            <Button onClick={() => navigate(-1)} style={{ marginBottom: 16 }}>
                返回
            </Button>

            <Descriptions
                title={`扫描结果 ID：${result.id}`}
                column={1}
                layout="horizontal"
                style={{ maxWidth: 800 }}
            >
                <Descriptions.Item label="漏洞名称">{result.vulnerability}</Descriptions.Item>
                <Descriptions.Item label="目标地址">{result.target}</Descriptions.Item>
                <Descriptions.Item label="风险等级">
                    <Tag color={severityColor[result.severity] || 'gray'}>
                        {result.severity}
                    </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="所属任务 ID">{result.taskID}</Descriptions.Item>
                <Descriptions.Item label="时间戳">{result.timestamp}</Descriptions.Item>
                <Descriptions.Item label="详细信息">
                    <Typography.Paragraph copyable style={{ whiteSpace: 'pre-wrap' }}>
                        {result.detail}
                    </Typography.Paragraph>
                </Descriptions.Item>
            </Descriptions>
        </div>
    );
}
