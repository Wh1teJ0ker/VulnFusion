import React, { useState } from 'react';
import { Button, Drawer, Form, Input, Message } from '@arco-design/web-react';
import { createTask } from '../../services/task';

export default function CreateTaskForm({ onSuccess }) {
    const [visible, setVisible] = useState(false);
    const [loading, setLoading] = useState(false);

    const [form] = Form.useForm();

    const handleSubmit = async () => {
        const values = form.getFieldsValue();
        setLoading(true);
        try {
            await createTask(values);
            Message.success('任务创建成功');
            setVisible(false);
            form.resetFields();
            onSuccess(); // 通知父组件刷新任务列表
        } catch (err) {
            Message.error(err?.message || '任务创建失败');
        } finally {
            setLoading(false);
        }
    };

    return (
        <>
            <Button type="primary" onClick={() => setVisible(true)} style={{ marginBottom: 20 }}>
                创建任务
            </Button>

            <Drawer
                title="新建扫描任务"
                width={400}
                visible={visible}
                onCancel={() => setVisible(false)}
                footer={
                    <Button type="primary" onClick={handleSubmit} loading={loading}>
                        提交
                    </Button>
                }
            >
                <Form form={form} layout="vertical" autoComplete="off">
                    <Form.Item
                        label="目标地址"
                        field="target"
                        rules={[{ required: true, message: '请输入目标地址' }]}
                    >
                        <Input placeholder="https://example.com" />
                    </Form.Item>

                    <Form.Item
                        label="Nuclei 模板路径"
                        field="template"
                        rules={[{ required: true, message: '请输入模板路径' }]}
                    >
                        <Input placeholder="cves/2021/*.yaml" />
                    </Form.Item>
                </Form>
            </Drawer>
        </>
    );
}
