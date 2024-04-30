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
      </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import api from "@/services/axios"; 

const props = defineProps({
  userId: {
    type: String,
    required: true
  }
});

const userProfile = ref(null);
const isOwnProfile = computed(() => props.userId === localStorage.getItem('userId'));

const fetchUserProfile = async () => {
  console.log(props.userId)
  console.log(localStorage.getItem('userId'))  
  try {
    const response = await api.get(`/users/id/${props.userId}`);
    console.log(response.data)
    console.log(response)
    userProfile.value = response.data;
  } catch (error) {
    console.error("Error fetching user profile:", error);
  }
};

const toggleFollow = async () => {
  const endpoint = `/users/follows/${props.userId}`;
  try {
    await api({
      method: userProfile.value.isFollowing ? 'delete' : 'post',
      url: endpoint
    });
    userProfile.value.isFollowing = !userProfile.value.isFollowing;
  } catch (error) {
    console.error("Error toggling follow:", error);
  }
};

const toggleBan = async () => {
  const endpoint = `/users/bans/${props.userId}`;
  try {
    await api({
      method: userProfile.value.isBanned ? 'delete' : 'post',
      url: endpoint
    });
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
