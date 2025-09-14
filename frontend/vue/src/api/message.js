import api from '@/api'

// 获取未读消息数量
export function getUnreadMessageCount(id) {
    return api.get('/message/unreadCount',{
        params: {
            id: id
        }
    })
}

// 获取消息列表
export function getMessageList(page,pageSize,id) {
    return api.get('/message/list',{
        params: {
            page: page,
            pageSize: pageSize,
            id: id
        }
    })
}

// 更新消息
export function updateMessage(data){
    return api.put('/message/update',data)
}

// 删除消息
export function deleteMessage(id){
    return api.delete('/message/delete',{
        params: {
            id: id
        }
    })

}

// 一键已读所有消息
export function readAllMessage(id){
    return api.put('/message/readAll',id)
}

// 处理请求
export function handleRequest(data){
    return api.post('/message/handle',data)
}

// 添加消息
export function addMessage(data){
    return api.post('/message/add',data)
}