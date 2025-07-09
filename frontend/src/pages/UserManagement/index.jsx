import React, { useEffect, useState } from 'react';
import {
    Table,
    Button,
    Message,
    Typography,
    Space,
    Modal,
    Select,
    Input,
} from '@arco-design/web-react';
import {
    getAllUsers,
    deleteUserById,
    updateUserById,
    resetUserPassword,
} from '../../services/user';
import { useUserStore } from '../../store/user';
import { register } from '../../services/auth';

const Option = Select.Option;

export default function UserManagement() {
    const { role: currentUserRole } = useUserStore(); // 当前登录用户的角色
    const [loading, setLoading] = useState(false);
    const [users, setUsers] = useState([]);
    const [editUser, setEditUser] = useState(null);
    const [resetTarget, setResetTarget] = useState(null);
    const [newPassword, setNewPassword] = useState('');
    const [addModalVisible, setAddModalVisible] = useState(false);
    const [newUser, setNewUser] = useState({ username: '', password: '', role: 'user' });

    const fetchUsers = async () => {
        setLoading(true);
        try {
            const res = await getAllUsers();
            const normalized = (res || []).map(u => ({
                id: u.ID,
                username: u.Username,
                role: u.Role,
            }));
            setUsers(normalized);
        } catch (err) {
            console.error('❌ 获取用户失败:', err);
            Message.error(err?.message || '获取用户列表失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchUsers();
    }, []);

    const handleDelete = (id) => {
        Modal.confirm({
            title: '确认删除',
            content: '删除后该用户将无法登录，是否继续？',
            onOk: async () => {
                try {
                    await deleteUserById(id);
                    Message.success('删除成功');
                    await fetchUsers();
                } catch (err) {
                    console.error('❌ 删除失败:', err);
                    Message.error(err?.message || '删除失败');
                }
            },
        });
    };

    const handleEdit = (record) => {
        setEditUser(record);
    };

    const handleUpdate = async () => {
        try {
            await updateUserById(editUser.id, { role: editUser.role });
            Message.success('更新成功');
            setEditUser(null);
            await fetchUsers();
        } catch (err) {
            console.error('❌ 更新失败:', err);
            Message.error(err?.message || '更新失败');
        }
    };

    const handleResetPassword = (record) => {
        setResetTarget(record);
        setNewPassword('');
    };

    const submitResetPassword = async () => {
        if (!newPassword) {
            Message.warning('请输入新密码');
            return;
        }

        try {
            await resetUserPassword(resetTarget.id, newPassword);
            Message.success('密码重置成功');
            setResetTarget(null);
            setNewPassword('');
        } catch (err) {
            console.error('❌ 密码重置失败:', err);
            Message.error(err?.message || '密码重置失败');
        }
    };

    const handleAddUser = async () => {
        const { username, password, role } = newUser;

        if (!username || !password) {
            Message.warning('用户名和密码不能为空');
            return;
        }

        try {
            await register({ username, password, role });
            Message.success('用户添加成功');
            setAddModalVisible(false);
            setNewUser({ username: '', password: '', role: 'user' });
            await fetchUsers();
        } catch (err) {
            console.error('❌ 添加用户失败:', err);
            Message.error(err?.message || '添加用户失败');
        }
    };

    const columns = [
        {
            title: '用户ID',
            dataIndex: 'id',
            width: 80,
        },
        {
            title: '用户名',
            dataIndex: 'username',
        },
        {
            title: '角色',
            dataIndex: 'role',
            render: (val) => (val === 'admin' ? '管理员' : '普通用户'),
        },
        {
            title: '操作',
            render: (_, record) => (
                <Space>
                    <Button size="mini" onClick={() => handleEdit(record)}>编辑</Button>
                    <Button
                        size="mini"
                        status="warning"
                        onClick={() => handleResetPassword(record)}
                    >
                        重置密码
                    </Button>
                    <Button
                        size="mini"
                        status="danger"
                        onClick={() => handleDelete(record.id)}
                    >
                        删除
                    </Button>
                </Space>
            ),
        },
    ];

    if (currentUserRole !== 'admin') {
        return <Typography.Text>无权访问该页面</Typography.Text>;
    }

    return (
        <div>
            <Typography.Title heading={5}>用户管理</Typography.Title>
            <Typography.Text type="secondary">
                当前用户数量：{users.length}
            </Typography.Text>

            <div style={{ margin: '12px 0' }}>
                <Button type="primary" onClick={() => setAddModalVisible(true)}>
                    添加新用户
                </Button>
            </div>

            <Table
                rowKey="id"
                columns={columns}
                data={users}
                loading={loading}
                pagination={false}
            />

            {/* 添加新用户 */}
            <Modal
                title="添加新用户"
                visible={addModalVisible}
                onOk={handleAddUser}
                onCancel={() => setAddModalVisible(false)}
            >
                <Input
                    placeholder="用户名"
                    value={newUser.username}
                    onChange={(val) => setNewUser({ ...newUser, username: val })}
                    style={{ marginBottom: 10 }}
                />
                <Input.Password
                    placeholder="密码"
                    value={newUser.password}
                    onChange={(val) => setNewUser({ ...newUser, password: val })}
                    style={{ marginBottom: 10 }}
                />
                <Select
                    value={newUser.role}
                    onChange={(val) => setNewUser({ ...newUser, role: val })}
                    style={{ width: '100%' }}
                >
                    <Option value="admin">管理员</Option>
                    <Option value="user">普通用户</Option>
                </Select>
            </Modal>

            {/* 编辑角色 */}
            <Modal
                title="编辑用户角色"
                visible={!!editUser}
                onOk={handleUpdate}
                onCancel={() => setEditUser(null)}
            >
                <Select
                    value={editUser?.role}
                    onChange={(val) => setEditUser({ ...editUser, role: val })}
                    style={{ width: '100%' }}
                >
                    <Option value="admin">管理员</Option>
                    <Option value="user">普通用户</Option>
                </Select>
            </Modal>

            {/* 重置密码 */}
            <Modal
                title={`重置用户密码（${resetTarget?.username}）`}
                visible={!!resetTarget}
                onOk={submitResetPassword}
                onCancel={() => setResetTarget(null)}
            >
                <Input.Password
                    placeholder="请输入新密码"
                    value={newPassword}
                    onChange={setNewPassword}
                />
            </Modal>
        </div>
    );
}
