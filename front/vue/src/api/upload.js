import api from '@/api'

export function uploadImg(data){
    return api.post('/upload/img', data)
}