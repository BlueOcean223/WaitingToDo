import api from '@/api'

// 根据id查询好友信息
export function getFiendInfo(id){
    return api.get('/friend/info',{
        params: {
            id: id
        }
    })
}

// 获取好友列表
export function getFriendList(id){
    return api.get('/friend/list',{
        params: {
            id: id
        }
    })
}

// 根据邮箱获取用户信息
export function searchUserInfoByEmail(email){
    return api.get('/friend/search',{
        params: {
            email: email
        }
    })
}

// 发送添加好友请求
export function sendAddFriendRequest(data){
    return api.post('/friend/add', data)
}

// 删除好友
export function deleteFriend(userId,friendId){
    return api.delete('/friend/delete',{
        params:{
            userId: userId,
            friendId: friendId
        }
    })
}