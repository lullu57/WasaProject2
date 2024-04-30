<template>
    <div class="photo-card">
      <img :src="photo.url" alt="Photo" class="photo-image"/>
      <!-- Add more photo details or actions here -->
    </div>
  </template>
  
  <script>
  export default {
    props: {
      photo: Object
    }
  }
  </script>
  
  <style scoped>
  .photo-card {
    border: 1px solid #ccc;
    border-radius: 4px;
    padding: 10px;
    margin-bottom: 10px;
  }
  
  .photo-image {
    width: 100%;
    height: auto;
    border-radius: 4px;
  }
  </style>
  