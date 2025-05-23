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