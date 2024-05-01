<script setup>
import { ref, computed } from 'vue';
import { RouterLink, RouterView, useRoute } from 'vue-router';
import UploadImage from './components/UploadImage.vue';

const userId = ref(localStorage.getItem('userId'));
const isAuthenticated = computed(() => !!userId.value);

function logout() {
  localStorage.removeItem('userId');
  userId.value = null;
  window.location.href = '/';
}
</script>

<template>
  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASA Photo</a>
    <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
  </header>

  <div class="container-fluid">
    <div class="row">
      <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
        <div class="position-sticky pt-3 sidebar-sticky">
          <ul class="nav flex-column">
            <li class="nav-item">
              <RouterLink to="/stream" class="nav-link">
                <svg class="feather"><use xlink:href="/feather-sprite.svg#home"/></svg>
                Home
              </RouterLink>
            </li>
            <li class="nav-item" v-if="isAuthenticated">
				<RouterLink :to="'/profile/' + userId" class="nav-link">
				<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
				My Profile
				</RouterLink>
			</li>
            <li class="nav-item" v-if="!isAuthenticated">
              <RouterLink to="/" class="nav-link">
                <svg class="feather"><use xlink:href="/feather-sprite.svg#log-in"/></svg>
                Login
              </RouterLink>
            </li>
            <li class="nav-item" v-if="isAuthenticated">
              <button @click="logout" class="btn nav-link">
                <svg class="feather"><use xlink:href="/feather-sprite.svg#log-out"/></svg>
                Logout
              </button>
            </li>
            <li class="nav-item">
              <RouterLink to="/discover" class="nav-link" v-if="isAuthenticated">
                <svg class="feather"><use xlink:href="/feather-sprite.svg#users"/></svg>
                Discover Users
              </RouterLink>
            </li>
            <li class="nav-item">
              <UploadImage v-if="isAuthenticated"/>
            </li>
          </ul>
        </div>
      </nav>

      <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
        <RouterView />
        
      </main>
    </div>
  </div>
</template>

<style>
/* Main CSS adjustments */
.navbar {
  background-color: #007bff; /* Example of a brighter color */
  padding: 10px 15px; /* Increased padding */
}

.navbar-brand {
  font-size: 1.25rem; /* Larger font size */
}

.sidebar {
  width: 250px; /* Increased width */
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 10px; /* Adds space between icon and text */
  transition: background-color 0.3s ease-in-out; /* Smooth transition for hover */
}

.nav-link:hover {
  background-color: #f8f9fa; /* Light background on hover */
}

/* Use Bootstrap's card class for photo presentations */
.photo-card {
  display: block;
  border: none;
  margin-bottom: 20px;
  box-shadow: 0 0 10px rgba(0,0,0,0.1); /* Soft shadow for depth */
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .sidebar {
    width: 100%;
    position: absolute;
    z-index: 1000;
  }
}

/* Adding footer */
footer {
  padding: 20px;
  text-align: center;
  background-color: #f5f5f5;
}

</style>
