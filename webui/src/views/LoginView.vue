<template>
    <div class="login-container">
      <h1>Login</h1>
      <input v-model="username" placeholder="Enter username" />
      <button @click="login">Login</button>
      <p v-if="error">{{ error }}</p>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  import api from "@/services/axios"; 
  
  export default {
    data() {
      return {
        username: '',
        error: ''
      };
    },
    methods: {
      async login() {
        try {
          const response = await api.post('/session', { name: this.username });
          console.log(response.data)
          console.log(response.data.token)
          localStorage.setItem("userId", response.data.token);
          axios.defaults.headers.common['Authorization'] = response.data;
          this.$router.push('/'); // Redirect to home after successful login
          location.reload()
        } catch (err) {
          console.log(err);
          this.error = 'Failed to login. Please try again.';
          console.error(err);
        }
      }
    }
  }
  </script>
  
  <style>
  /* Add styles here */
  </style>
  