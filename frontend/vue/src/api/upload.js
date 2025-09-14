import api from '@/api'

export function uploadImg(data){
    return api.post('/upload/img', data)
}

export function uploadFile(id,data){
    return api.post(`/upload/${id}/file`,data,{
        headers: {
            'Content-Type':'multipart/form-data'
        }
    })
}

export function deleteFile(data){
    return api.delete('/upload/deletefile',{
        data: {
            delete_ids: data
        }
    })
}