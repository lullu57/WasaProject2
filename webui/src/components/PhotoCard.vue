<template>
  <div class="photo-card">
    <img :src="'data:image/jpeg;base64,' + photo.imageData" alt="Photo" class="photo-image"/>
    <div class="photo-info">
      <h3>{{ photo.username }}</h3>
      <p>{{ formatDate(photo.timestamp) }}</p>
      <div class="photo-actions">
        <button @click="toggleLike">{{ photo.isLiked ? 'Unlike' : 'Like' }} ({{ photo.likesCount }})</button>
        <button @click="toggleComments">Comments ({{ photo.comments.length }})</button>
      </div>
      <div v-if="showComments" class="comments-section">
        <div class="comment-form">
          <input v-model="newComment" placeholder="Write a comment..." class="comment-input"/>
          <button @click="postComment">Post Comment</button>
        </div>
        <div class="comment" v-for="comment in photo.comments" :key="comment.commentId">
          <strong>{{ comment.username }}</strong>: {{ comment.content }}
        </div>
      </div>
    </div>
  </div>
</template>

  
  
  
<script>
import api from '@/services/axios';

export default {
  props: {
    photo: Object
  },
  data() {
    return {
      showComments: true,
      isLiked: false, // Will be updated based on API response
      newComment: ''
    };
  },
  mounted() {
    this.checkIfLiked(); // Check if the photo is liked on component mount
  },
  methods: {
    async checkIfLiked() {
      const config = {
        headers: {
          Authorization: `${localStorage.getItem('userId')}`
        }
      };
      try {
        const response = await api.get(`/likes/${this.photo.photoId}`, config);
        this.isLiked = response.data.liked;
      } catch (error) {
        console.error('Failed to check like status', error);
      }
    },
    async toggleLike() {
      const config = {
        headers: {
          Authorization: `${localStorage.getItem('userId')}`
        }
      };
      try {
        if (!this.isLiked) {
          await api.post(`/photos/${this.photo.photoId}/likes`, {}, config);
          this.photo.likesCount++;
        } else {
          await api.delete(`/photos/${this.photo.photoId}/likes`, config);
          this.photo.likesCount--;
        }
        this.isLiked = !this.isLiked;
      } catch (error) {
        console.error('Failed to toggle like', error);
      }
    },
    toggleComments() {
      this.showComments = !this.showComments;
    },
    async postComment() {
      if (this.newComment.trim() !== '') {
        const config = {
          headers: {
            Authorization: `${localStorage.getItem('userId')}`
          }
        };
        const response = await api.post(`/photos/${this.photo.photoId}/comments`, { content: this.newComment }, config);
        let username = 'You'; // Ideally fetch from server or use global state
        this.photo.comments.push({
          username,
          content: this.newComment,
          commentId: response.data.commentId
        });
        this.newComment = '';
      }
    },
    formatDate(value) {
      return new Date(value).toLocaleString();
    }
  }
}
</script>




<style scoped>
.photo-card {
  border: 1px solid #ccc;
  border-radius: 4px;
  padding: 10px;
  margin-bottom: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  background-color: #fff; /* Adds a clean background color */
}

.photo-image {
  max-width: 90%; /* reduces the maximum width to fit better */
  max-height: 300px; /* reduces the maximum height */
  object-fit: contain; /* ensures the image fits nicely within the constraints */
  border-radius: 4px;
}

.photo-info {
  width: 100%;
  text-align: center;
  padding-top: 10px;
}

.photo-actions {
  display: flex;
  justify-content: space-around; /* changes from space-between to space-around for better distribution */
  padding: 5px 0; /* adds padding around buttons */
}

button {
  background-color: #007BFF; /* bootstrap primary color for consistency */
  color: white;
  border: none;
  border-radius: 4px;
  padding: 5px 10px;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #0056b3; /* darker shade on hover */
}

.comments-section {
  margin-top: 10px;
  width: 100%;
}

.comment-form {
  display: flex;
  margin-top: 5px;
}

.comment-input {
  flex-grow: 1;
  border: 1px solid #ccc;
  border-radius: 4px;
  padding: 5px;
}

.comment {
  background-color: #f0f0f0;
  padding: 5px;
  border-radius: 3px;
  margin-top: 2px;
}
</style>
