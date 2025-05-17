<template>
    <Navbar :activeIndex="$route.path"/>
    欢迎来到个人详情页
</template>

<script>
    import Navbar from '@/components/Navbar.vue'

    export default {
        name: 'PersonalPage',
        components: {
            Navbar
        },
    }
</script>