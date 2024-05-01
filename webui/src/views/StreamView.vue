<template>
    <div class="stream-view">
      <div v-if="photos.length > 0" class="photo-grid">
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
        photos: [],
        error: ''
      };
    },
    async mounted() {
      await this.fetchStreamPhotos();
    },
    methods: {
      async fetchStreamPhotos() {
        try {
          const response = await api.get('/stream', {
            headers: { Authorization: localStorage.getItem('userId') }
          });
          const photoIds = response.data; // Assuming this is an array of photo IDs
          await this.fetchPhotoDetails(photoIds);
        } catch (error) {
          console.error('Failed to fetch stream photos:', error);
          this.error = "Failed to load photos. Please try again later.";
        }
      },
      async fetchPhotoDetails(photoIds) {
        this.photos = await Promise.all(photoIds.map(async (photoId) => {
          try {
            const res = await api.get(`/photos/${photoId}`, {
              headers: { Authorization: localStorage.getItem('userId') }
            });
            const photo = res.data;
            // Process comments
            photo.comments = await Promise.all(photo.comments.map(async (comment) => {
              const userResponse = await api.get(`/username/${comment.userId}`);
              comment.username = userResponse.data.username;
              comment.isOwner = comment.userId === localStorage.getItem('userId');
              return comment;
            }));
            return photo;
          } catch (error) {
            console.error("Error fetching photo details:", error);
            return {}; // Return empty object or handle errors appropriately
          }
        }));
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
    padding: 20px;
  }
  
  .photo-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 20px;
  }
  
  p {
    color: #666;
    text-align: center;
  }
  
  @media (max-width: 768px) {
    .photo-grid {
      grid-template-columns: 1fr; /* Adjust grid for smaller screens */
    }
  }
  </style>
  