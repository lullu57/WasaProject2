<template>
    <div class="stream-view">
      <div v-if="photos.length > 0">
        <PhotoCard 
          v-for="photo in photos" 
          :key="photo.photoId"
          :photo="photo"
        />
      </div>
      <div v-else>
        <p>No photos to display. Start following people to see their photos here.</p>
      </div>
    </div>
  </template>
  
  <script>
  import PhotoCard from '@/components/PhotoCard.vue';
  import api from '@/services/axios';
  
  export default {
    components: {
      PhotoCard
    },
    data() {
      return {
        photos: []
      };
    },
    async mounted() {
      await this.fetchStream();
    },
    methods: {
      async fetchStream() {
        try {
          const response = await api.get('/stream', {
            headers: { Authorization: localStorage.getItem('userId') }
          });
          this.photos = response.data;
        } catch (error) {
          console.error('Failed to fetch stream:', error);
        }
      }
    }
  }
  </script>
  
  <style scoped>
  .stream-view {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
  }
  </style>
  