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
        <button @click="postComment" class="post-comment">Post</button> 
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
  border: 2px solid #ccc;
  border-radius: 4px;
  padding: 10px;
  margin-bottom: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  width: 100%;
  max-width: 1200px;
  
}

.photo-image {
  max-width: 100%; /* ensures the image is not wider than the card */
  max-height: 900px; /* sets a maximum height for the image */
  object-fit: contain; /* ensures the image fits nicely within the constraints */
  border-radius: 4px;
}

.photo-info {
  width: 100%;
  text-align: center;
  padding-top: 10px;
}

.photo-actions button {
  background-color: #ff5e8a; /* a light reddish-pink color for like button */
  color: white;
  border: none;
  border-radius: 4px;
  padding: 5px 10px;
  cursor: pointer;
  transition: background-color 0.3s;
  object-fit: contain;
}

.photo-actions button:hover {
  background-color: #d11a3e; /* a deeper shade on hover */
}

.comments-section {
  margin-top: 10px;
  width: 100%; /* ensures the comments section uses the full width of the card */
}

.comment-form {
  display: flex;
  justify-content: space-between;
  width: 100%; /* ensures the form uses the full width of the card */
  margin-top: 5px;
}

.comment-input {
  flex-grow: 1;
  margin-right: 10px;
  padding: 5px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

button.post-comment {
  background-color: #1980ed; /* primary button color */
  padding: 5px 8px; /* smaller padding */
  font-size: 0.8rem; /* smaller font size */
}

.comment {
  background-color: #d7d6d6;
  padding: 5px;
  border-radius: 3px;
  margin-top: 2px;
}

.username {
  font-size: 0.7rem; /* Smaller font size for usernames */
  color: #555; /* Dark gray for better readability */
  margin-bottom: 2px;
}
</style>

