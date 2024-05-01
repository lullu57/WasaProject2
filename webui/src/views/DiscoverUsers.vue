<template>
    <div class="discover-users">
      <h2>Discover Users</h2>
      <ul class="user-list">
        <li v-for="user in users" :key="user.userId">
          {{ user.username }}
          <button @click="toggleFollow(user)" :disabled="user.processing">
            {{ user.isFollowing ? 'Unfollow' : 'Follow' }}
          </button>
        </li>
      </ul>
    </div>
  </template>
  
  <script>
  import api from '@/services/axios';
  
  export default {
    data() {
      return {
        users: []
      };
    },
    async mounted() {
      await this.fetchUsers();
    },
    methods: {
      async fetchUsers() {
        try {
          const response = await api.get('/users', {
            headers: { Authorization: localStorage.getItem('userId') }
          });
          this.users = response.data.map(user => ({
            ...user,
            isFollowing: false, // default to false, will be updated
            processing: false
          }));
          await this.checkFollowStatus();
        } catch (error) {
          console.error('Failed to fetch users:', error);
        }
      },
      async checkFollowStatus() {
        try {
          // Get following status for each user
          for (const user of this.users) {
            const res = await api.get(`/follows/${user.userId}`, {
              headers: { Authorization: localStorage.getItem('userId') }
            });
            user.isFollowing = res.data.isFollowed;
          }
        } catch (error) {
          console.error('Failed to fetch follow status:', error);
        }
      },
      async toggleFollow(user) {
        user.processing = true; // Indicate processing
        try {
          if (user.isFollowing) {
            await api.delete(`/users/follows/${user.userId}`, {
              headers: { Authorization: localStorage.getItem('userId') }
            });
            user.isFollowing = false;
          } else {
            await api.post(`/users/follows/${user.userId}`, {}, {
              headers: { Authorization: localStorage.getItem('userId') }
            });
            user.isFollowing = true;
          }
        } catch (error) {
          console.error('Failed to toggle follow:', error);
        } finally {
          user.processing = false; // Processing done
        }
      }
    }
  }
  </script>
  
  <style scoped>
  .discover-users {
    padding: 20px;
  }
  
  .user-list {
    list-style-type: none;
    padding: 0;
  }
  
  .user-list li {
    margin-bottom: 10px;
    display: flex;
    align-items: center;
  }
  
  button {
    margin-left: 10px;
  }
  </style>
  