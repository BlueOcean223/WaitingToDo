import { defineStore } from 'pinia'

// 全局管理用户信息
export const useUserStore = defineStore('user',{
    state: ()  => ({
        userInfo: JSON.parse(localStorage.getItem('user')) || {},
    }),
    actions: {
        updateUserInfo(userInfo){
            this.userInfo = userInfo
            localStorage.setItem('user', JSON.stringify(userInfo))
        }
    }
})