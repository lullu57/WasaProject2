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
      photos: [],
      error: '' // To handle errors and display messages
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
        if (response && response.data) {
          this.photos = await Promise.all(response.data.map(async photo => {
            // Fetch usernames for each comment on the photo
            photo.comments = await Promise.all(photo.comments.map(async (comment) => {
              const userResponse = await api.get(`/username/${comment.userId}`);
              comment.username = userResponse.data.username;
              comment.isOwner = comment.userId === localStorage.getItem('userId');
              return comment;
            }));
            return photo;
          }));
        }
      } catch (error) {
        console.error('Failed to fetch stream:', error);
        this.error = "Failed to load photos. Please try again later.";
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
  padding: 20px;
}
p {
  color: #666;
}
</style>
