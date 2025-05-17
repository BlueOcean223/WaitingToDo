<template>
    <NavBar :activeIndex="$route.path" />
    欢迎来到好友页
</template>

<script>
import NavBar from '@/components/NavBar.vue'
export default {
  name: 'FriendsPage',
  components: {
    NavBar
  },
  data(){
    return {
        
    }
  },
}
</script>