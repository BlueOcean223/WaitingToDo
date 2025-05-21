import api from '@/api'

// 分页查询
export function getList(page,pageSize){
    return api.get('/task/list',{
        params: {
            page: page,
            pageSize: pageSize
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