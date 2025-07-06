import React from 'react';
import { Card, Typography, Grid, Statistic } from '@arco-design/web-react';
import BaseLayout from '../components/Layout';

const { Title, Text } = Typography;
const Row = Grid.Row;
const Col = Grid.Col;

function Dashboard() {
    return (
        <BaseLayout>
            <Title heading={3}>仪表盘</Title>
            <Text style={{ display: 'block', marginBottom: 16 }}>
                欢迎使用 VulnFusion 漏洞扫描平台 👋
            </Text>

            <Row gutter={16}>
                <Col span={8}>
                    <Card bordered>
                        <Statistic title="今日任务数" value={8} />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card bordered>
                        <Statistic title="扫描中的任务" value={3} />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card bordered>
                        <Statistic title="发现漏洞数" value={17} />
                    </Card>
                </Col>
            </Row>

            <Card style={{ marginTop: 24 }} title="系统说明">
                <Text>
                    后续将在此展示统计图、最新任务、漏洞概览等内容，敬请期待。
                </Text>
            </Card>
        </BaseLayout>
    );
}

export default Dashboard;
