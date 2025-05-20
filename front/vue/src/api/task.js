import api from '@/api'

export function getList(page,pageSize){
    return api.get('/task/list',{
        params: {
            page: page,
            pageSize: pageSize
        }
    })
}


export function add(data){
    return api.post('/task/add', data)
}