import request from '../utils/request';

// 当前用户的任务列表
export function getTasks() {
    return request.get('/tasks');
}

// 管理员：获取所有任务
export function getAllTasks() {
    return request.get('/admin/tasks');
}

// 创建任务
export function createTask(data) {
    return request.post('/tasks', data);
}

// 获取任务详情
export function getTaskDetail(id) {
    return request.get(`/tasks/${id}`);
}

// 删除任务
export function deleteTask(id) {
    return request.delete(`/tasks/${id}`);
}

// 批量删除
export function batchDeleteTasks(ids) {
    return request.post('/tasks/batch_delete', { ids });
}

// 更新任务状态
export function updateTaskStatus(id, status) {
    return request.post('/tasks/status', { id, status });
}
