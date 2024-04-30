<template>
    <div class="login-container">
      <h1>Login</h1>
      <input v-model="username" placeholder="Enter username" />
      <button @click="login">Login</button>
      <p v-if="error">{{ error }}</p>
    </div>
  </template>
  
  <script>
  import axios from './services/axios'
  
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
          const response = await axios.post('/session', { name: this.username });
          localStorage.setItem('userId', response.data.identifier);
          this.$router.push('/'); // Redirect to home or dashboard
          console.log('Logged in successfully');
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
  