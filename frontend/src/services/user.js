import request from '../utils/request';

/**
 * 获取当前用户信息
 * @param {string} token - 访问令牌
 */
export function getUserInfo(token) {
    return request.get('/user/info', {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });
}

/**
 * 获取所有用户（管理员权限）
 */
export function getAllUsers() {
    return request.get('/admin/users');

}

/**
 * 删除指定用户（管理员权限）
 * @param {number} id - 用户ID
 */
export function deleteUserById(id) {
    return request.delete(`/admin/users/${id}`);
}

/**
 * 更新用户信息（管理员权限）
 * @param {number} id - 用户ID
 * @param {Object} data - 更新字段（支持 password, role）
 */
export function updateUserById(id, data) {
    return request.put(`/admin/users/${id}`, data);
}

// 管理员重置用户密码
export function resetUserPassword(id, newPassword) {
    return request.put(`/admin/users/${id}/password`, { password: newPassword });
}
