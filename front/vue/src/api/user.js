import api from '@/api'

// 校验验证码
export function checkCaptcha(data){
    return api.post('/user/checkCaptcha', data)
}

// 重置密码
export function reset(data){
    return api.post('/user/reset', data)
}

// 更新个人信息
export function update(data){
    return api.post('/user/update', data)
}