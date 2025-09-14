import api from '@/api'

// 分页查询
export function getList(page,pageSize,status){
    return api.get('/task/list',{
        params: {
            page: page,
            pageSize: pageSize,
            status: status
        }
    })
}

// 添加任务
export function add(data){
    return api.post('/task/add', data)
}

// 修改任务
export function update(data) {
    return api.put('/task/update',data)
}

// 删除任务
export function remove(id) {
    return api.delete('/task/delete',{
        params: {
            id: id
        }
    })
}

// 获取紧急任务
export function getUrgent() {
    return api.get('/task/urgent')
}

// 获取小组任务列表
export function getTeamTaskList(page,pageSize,id){
    return api.get('/task/teamList',{
        params: {
            page: page,
            pageSize: pageSize,
            id: id
        }
    })
}

// 删除小组任务
export function removeTeamTask(taskId,userId){
    return api.delete('/task/team/delete',{
        params: {
            taskId: taskId,
            userId: userId
        }
    })
}

// 添加小组任务
export function addTeamTask(data){
    return api.post( '/task/team/add', data)
}

// 小组成员完成任务
export function completeTeamTask(data){
    return api.put('/task/team/complete', data)
}

// 邀请成员
export function inviteMember(data){
    return api.post('/task/team/invite', data)
}

// 获取小组任务邀请码
export function getTeamTaskInviteCode(taskId) {
    return api.get('/task/team/inviteCode', {
        params: {
            taskId: taskId
        }
    })
}

// 根据邀请码加入小组
export function joinTeamByCode(inviteCode) {
    return api.post('/task/team/codejoin', {
        inviteCode: inviteCode
    })
}