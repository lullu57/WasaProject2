<template>
  <div class="profile-view">
    <div v-if="userProfile" class="info-container">
      <p>Username: {{ userProfile.username }}</p>
      <p>Followers: {{ userProfile.followers?.length || '0' }}</p>
      <p>Following: {{ userProfile.following?.length || '0' }}</p>
      <p>Posts: {{ detailedPhotos.length || '0' }}</p>
      <button v-if="!isOwnProfile" @click="toggleFollow">
        {{ userProfile.isFollowing ? 'Unfollow' : 'Follow' }}
      </button>
      <button v-if="!isOwnProfile" @click="toggleBan">
        {{ userProfile.isBanned ? 'Unban' : 'Ban' }}
      </button>
    </div>
    <div v-else>
      <p>No profile data available.</p>
    </div>
    <div class="gallery">
      <PhotoCard 
        v-for="photo in detailedPhotos" 
        :key="photo.photoId"
        :photo="photo"
      />
    </div>
  </div>
</template>


<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import api from "@/services/axios";
import PhotoCard from '@/components/PhotoCard.vue';

const route = useRoute();
const userId = route.params.profileId;
const userProfile = ref(null);
const detailedPhotos = ref([]);
const localStorageUserId = localStorage.getItem('userId');
const isOwnProfile = computed(() => userId === localStorageUserId);

const fetchUserProfile = async () => {
  try {
    const response = await api.get(`/users/id/${userId}`);
    userProfile.value = response.data;
    if (userProfile.value && userProfile.value.photos) {
      fetchPhotoDetails(userProfile.value.photos);
    }
  } catch (error) {
    console.error("Error fetching user profile:", error);
  }
};

const fetchPhotoDetails = async (photoIds) => {
  detailedPhotos.value = await Promise.all(photoIds.map(async (id) => {
    try {
      const res = await api.get(`/photos/${id}`);
      return res.data;
    } catch (error) {
      console.error("Error fetching photo details:", error);
      return null; // Handle errors or missing data gracefully
    }
  }));
};

const toggleFollow = async () => {
  const method = userProfile.value.isFollowing ? 'delete' : 'post';
  const endpoint = `/users/follows/${userId}`;
  await api[method](endpoint);
  userProfile.value.isFollowing = !userProfile.value.isFollowing;
};

const toggleBan = async () => {
  const method = userProfile.value.isBanned ? 'delete' : 'post';
  const endpoint = `/users/bans/${userId}`;
  await api[method](endpoint);
  userProfile.value.isBanned = !userProfile.value.isBanned;
};

onMounted(fetchUserProfile);
</script>




<style scoped>
.profile-view {
  padding: 20px;
}

.info-container {
  background-color: #f4f4f4;
  padding: 20px;
  margin-bottom: 20px;
}

.gallery {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 10px;
}
</style>
