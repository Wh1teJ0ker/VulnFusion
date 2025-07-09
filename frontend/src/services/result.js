import request from '../utils/request';

// 当前用户获取某个任务的结果
export function getResultsByTaskId(taskId) {
    return request.get(`/results/task/${taskId}`);
}

// 管理员获取所有结果
export function getAllResults() {
    return request.get('/admin/results');
}

// 获取单个扫描结果详情
export function getResultDetail(id) {
    return request.get(`/results/${id}`);
}

// 导出某任务结果（json）
export function exportResult(taskId) {
    return request.get(`/results/export/${taskId}`);
}
