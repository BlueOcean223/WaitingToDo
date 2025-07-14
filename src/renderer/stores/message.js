import { defineStore } from 'pinia'

export const useMessageStore = defineStore('message', {
  state: () => ({
    unreadCount: 0
  }),
  actions: {
    // 获取未读消息数量
    setUnreadCount(count){
        this.unreadCount = count
    },
    // 用户读取一条消息
    readMessage(){
        this.unreadCount--
    },
    // 用户读取全部消息
    readAllMessage(){
        this.unreadCount = 0
    }
  }
})