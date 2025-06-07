<template>
  <div class="not-found">
    <h1>这里没有东西噢</h1>
    <p>您访问的页面不存在</p>
    <router-link to="/" class="home-link">返回首页</router-link>
    
    <div class="image-container">
      <img 
        src="@/assets/404syz.jpg" 
        class="not-found-image"
      >
    </div>

  </div>
</template>

<script setup name="NotFound">
</script>

<style scoped>
.not-found {
  text-align: center;
  padding: 2rem;
  margin-top: 0.5rem;
}

.image-container {
  margin: 2rem 0;
}
.not-found-image {
  width: 30%;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.home-link {
  display: inline-block;
  padding: 0.5rem 1rem;
  background-color: #42b983;
  color: white;
  text-decoration: none;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.home-link:hover {
  background-color: #3aa876;
}
</style>