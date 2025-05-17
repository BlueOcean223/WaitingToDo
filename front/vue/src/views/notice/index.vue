<template>
    <Navbar :activeIndex="$route.path"/>
    欢迎来到通知页
</template>

<script>
    import Navbar from '@/components/Navbar.vue'
    export default {
        name: 'NoticePage',
        components: {
            Navbar
        },
    }
</script>