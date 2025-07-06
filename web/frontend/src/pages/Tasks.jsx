import React from 'react';
import BaseLayout from '../components/Layout';
import { Table, Button } from '@arco-design/web-react';

function Tasks() {
    return (
        <BaseLayout>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
                <h2 style={{ fontWeight: 'bold' }}>扫描任务</h2>
                <Button type="primary">新建任务</Button>
            </div>
            <Table
                columns={[
                    { title: '任务名称', dataIndex: 'name' },
                    { title: '目标地址', dataIndex: 'target' },
                    { title: '状态', dataIndex: 'status' },
                    { title: '操作', render: () => <a>查看</a> },
                ]}
                data={[]}
                pagination={false}
            />
        </BaseLayout>
    );
}

export default Tasks;
