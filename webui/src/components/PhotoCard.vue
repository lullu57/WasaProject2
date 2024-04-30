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
      showComments: false,
      isLiked: false,
      newComment: ''
    };
  },
  methods: {
    async toggleLike() {
      
      const config = {
        headers: {
          Authorization: `${localStorage.getItem('userId')}`  // Adjusted to use a proper token header if applicable
        }
      };
      console.log(this.isLiked)
      console.log(this)
      if (!this.isLiked) {
        this.photo.likesCount++;
        await api.post(`/photos/${this.photo.photoId}/likes`, {}, config);
      } else {
        this.photo.likesCount--;
        await api.delete(`/photos/${this.photo.photoId}/likes`, config);
      }
      this.isLiked = !this.isLiked;
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
        let username = 'You';  // Default username
        try {
          const userRes = await api.get(`/username/${localStorage.getItem('userId')}`, config);
          username = userRes.data.username;  // Fetch the username dynamically
        } catch (error) {
          console.error('Failed to fetch username', error);
        }
        this.photo.comments.push({
          username: username,
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
  }
  
  .photo-image {
    max-width: 100%; /* ensures the image is not wider than the card */
    max-height: 400px; /* sets a maximum height for the image */
    width: auto; /* maintains the aspect ratio */
    height: auto; /* maintains the aspect ratio */
    object-fit: contain; /* ensures the image fits nicely within the constraints, without stretching */
    border-radius: 4px;
  }
  
  .photo-info {
    width: 100%;
    text-align: center;
    padding-top: 10px;
  }
  
  .photo-actions {
    display: flex;
    justify-content: space-between;
    margin-top: 10px;
  }
  
  .comments-section {
    margin-top: 10px;
    width: 100%; /* ensures the comments section uses the full width of the card */
  }
  
  .comment-form {
    display: flex;
    justify-content: space-between;
    width: 100%; /* ensures the form uses the full width of the card */
  }
  
  .comment-input {
    flex-grow: 1;
    margin-right: 10px;
  }
  
  .comment {
    background-color: #f0f0f0;
    padding: 5px;
    border-radius: 3px;
    margin-top: 2px;
  }
  </style>
  
  