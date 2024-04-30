<script setup>
import { ref, onMounted, computed } from 'vue';
import { RouterLink, RouterView, useRoute } from 'vue-router';

const userId = ref(localStorage.getItem('userId'));
console.log(userId)
const isAuthenticated = computed(() => !!userId.value);

function logout() {
  localStorage.removeItem('userId');
  userId.value = null;
  window.location.href = '/';
}
</script>

<template>
  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">Example App</a>
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
				<RouterLink :to="{ name: 'Profile', params: { profileId: userId } }" class="nav-link">
					<svg class="feather"><use xlink:href="/feather-sprite.svg#user"/></svg>
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
</style>
