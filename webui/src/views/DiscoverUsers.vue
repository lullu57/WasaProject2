<template>
    <div class="discover-users">
      <h2>Discover Users</h2>
      <ul class="user-list">
        <li v-for="user in users" :key="user.userId">
          {{ user.username }}
          <button @click="followUser(user.userId)">Follow</button>
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
          this.users = response.data;
        } catch (error) {
          console.error('Failed to fetch users:', error);
        }
      },
      async followUser(userId) {
        try {
          await api.post(`/users/follows/${userId}`, {}, {
            headers: { Authorization: localStorage.getItem('userId') }
          });
          alert('Followed successfully!');
        } catch (error) {
          console.error('Failed to follow user:', error);
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
  }
  </style>
  