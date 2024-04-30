<template>
  <div class="profile-view">
    <div v-if="userProfile" class="info-container">
      <p>Username: {{ userProfile.username }}</p>
      <p>Followers: {{ userProfile.followers.length }}</p>
      <p>Following: {{ userProfile.following.length }}</p>
      <p>Posts: {{ userProfile.photos.length }}</p>
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
      <!-- Display gallery items if any -->
    </div>
  </div>
</template>


<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router'; // Import useRoute to access route params
import api from "@/services/axios"; 

const route = useRoute(); // Use useRoute to access the route parameters

const userId = route.params.userId; // Access userId from route parameters

const userProfile = ref(null);
const localStorageUserId = localStorage.getItem('userId'); // Access once and use in computed property
const isOwnProfile = computed(() => userId === localStorageUserId);

const fetchUserProfile = async () => {
  try {
    const response = await api.get(`/users/id/${userId}`);
    userProfile.value = response.data;
  } catch (error) {
    console.error("Error fetching user profile:", error);
  }
};

const toggleFollow = async () => {
  const method = userProfile.value.isFollowing ? 'delete' : 'post';
  const endpoint = `/users/follows/${userId}`;
  try {
    await api[method](endpoint);
    userProfile.value.isFollowing = !userProfile.value.isFollowing;
  } catch (error) {
    console.error("Error toggling follow:", error);
  }
};

const toggleBan = async () => {
  const method = userProfile.value.isBanned ? 'delete' : 'post';
  const endpoint = `/users/bans/${userId}`;
  try {
    await api[method](endpoint);
    userProfile.value.isBanned = !userProfile.value.isBanned;
  } catch (error) {
    console.error("Error toggling ban:", error);
  }
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
