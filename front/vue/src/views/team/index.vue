<template>
    <Navbar :activeIndex="$route.path"/>
    欢迎来到小组页
</template>

<script>
    import Navbar from '@/components/Navbar.vue'
    export default {
        name: 'TeamPage',
        components: {
            Navbar
        },
    }
</script>