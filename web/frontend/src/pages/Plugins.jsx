import React from 'react';
import BaseLayout from '../components/Layout';
import { Table, Button } from '@arco-design/web-react';

function Plugins() {
    return (
        <BaseLayout>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
                <h2 style={{ fontWeight: 'bold' }}>插件管理</h2>
                <Button type="primary">上传插件</Button>
            </div>
            <Table
                columns={[
                    { title: '插件名', dataIndex: 'name' },
                    { title: '类型', dataIndex: 'type' },
                    { title: '版本', dataIndex: 'version' },
                    { title: '操作', render: () => <a>详情</a> },
                ]}
                data={[]}
                pagination={false}
            />
        </BaseLayout>
    );
}

export default Plugins;
