import api from '@/api'

//  登录
export function login(data){
    return api.post('/auth/login',data);
}

//  注册
export function register(data){
    return api.post('/auth/register', data)
}

//  忘记密码
export function forget(data){
    return api.post('/auth/forget', data)
}

//  获取验证码
export function captcha(email){
    return api.get('/auth/captcha',{
        params: {email}
    });
}