import React, { useEffect, useState } from 'react';
import {
    Card,
    Grid,
    Statistic,
    Message,
    Typography,
    Divider,
} from '@arco-design/web-react';
import {
    IconBug,
    IconList,
    IconUser,
} from '@arco-design/web-react/icon';
import { useUserStore } from '../../store/user';
import { getTasks } from '../../services/task';
import { getResultsByTaskId, getAllResults } from '../../services/result';
import ReactECharts from 'echarts-for-react';
import dayjs from 'dayjs';

const Row = Grid.Row;
const Col = Grid.Col;
const normalizeSeverity = (s) => (s || '').trim().toLowerCase();



export default function Dashboard() {
    const { username, role } = useUserStore();
    const [taskCount, setTaskCount] = useState(0);
    const [resultCount, setResultCount] = useState(0);
    const [severityStats, setSeverityStats] = useState({
        low: 0,
        medium: 0,
        high: 0,
        critical: 0,
    });
    const [taskTrend, setTaskTrend] = useState({}); // key: date, value: count

    const loadData = async () => {
        try {
            const tasks = await getTasks();
            setTaskCount(tasks.length);

            const severityCounter = {
                low: 0,
                medium: 0,
                high: 0,
                critical: 0,
            };

            const trendCounter = {};
            const today = dayjs();

            // 初始化近 7 天日期
            for (let i = 6; i >= 0; i--) {
                const date = today.subtract(i, 'day').format('YYYY-MM-DD');
                trendCounter[date] = 0;
            }

            let totalResults = 0;

            if (role === 'admin') {
                const allResults = await getAllResults();
                totalResults = allResults.length;

                for (const r of allResults) {
                    const level = normalizeSeverity(r.severity);
                    const date = dayjs(r.created_at).format('YYYY-MM-DD');

                    if (severityCounter[level] !== undefined) {
                        severityCounter[level]++;
                    }
                    if (date in trendCounter) {
                        trendCounter[date]++;
                    }
                }
            } else {
                for (const task of tasks) {
                    const date = dayjs(task.createdAt).format('YYYY-MM-DD');
                    if (date in trendCounter) {
                        trendCounter[date]++;
                    }

                    const results = await getResultsByTaskId(task.ID);
                    totalResults += results.length;

                    for (const r of results) {
                        const level = normalizeSeverity(r.severity);
                        if (severityCounter[level] !== undefined) {
                            severityCounter[level]++;
                        }
                    }
                }
            }

            setResultCount(totalResults);
            setSeverityStats(severityCounter);
            setTaskTrend(trendCounter);
        } catch (err) {
            console.error('加载数据失败:', err);
            Message.error('加载图表数据失败');
        }
    };

    useEffect(() => {
        loadData();
    }, []);

    const pieOption = {
        title: { text: '风险等级分布', left: 'center' },
        tooltip: { trigger: 'item' },
        legend: { bottom: '0%', left: 'center' },
        series: [
            {
                type: 'pie',
                radius: '50%',
                data: [
                    { value: severityStats.low, name: '低风险' },
                    { value: severityStats.medium, name: '中风险' },
                    { value: severityStats.high, name: '高风险' },
                    { value: severityStats.critical, name: '严重风险' },
                ],
                emphasis: { itemStyle: { shadowBlur: 10, shadowOffsetX: 0 } },
            },
        ],
    };

    const lineOption = {
        title: { text: '近7天任务数量趋势', left: 'center' },
        xAxis: {
            type: 'category',
            data: Object.keys(taskTrend),
        },
        yAxis: { type: 'value' },
        series: [
            {
                data: Object.values(taskTrend),
                type: 'line',
                smooth: true,
            },
        ],
    };

    return (
        <div>
            <Typography.Title heading={4} style={{ marginBottom: 24 }}>
                欢迎回来，{username}（{role}）
            </Typography.Title>

            <Row gutter={16}>
                <Col span={8}>
                    <Card bordered>
                        <Statistic title="任务总数" value={taskCount} icon={<IconList />} />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card bordered>
                        <Statistic title="漏洞结果总数" value={resultCount} icon={<IconBug />} />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card bordered>
                        <Statistic title="当前用户" value={username} icon={<IconUser />} />
                    </Card>
                </Col>
            </Row>

            <Divider style={{ margin: '30px 0' }} />

            <Row gutter={16}>
                <Col span={12}>
                    <Card>
                        <ReactECharts option={pieOption} style={{ height: 400 }} />
                    </Card>
                </Col>
                <Col span={12}>
                    <Card>
                        <ReactECharts option={lineOption} style={{ height: 400 }} />
                    </Card>
                </Col>
            </Row>
        </div>
    );
}
