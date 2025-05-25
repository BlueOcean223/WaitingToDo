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