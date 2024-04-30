<template>
    <div class="photo-card">
      <img :src="'data:image/jpeg;base64,' + photo.imageData" alt="Photo" class="photo-image"/>
      <div class="photo-info">
        <h3>{{ photo.username }}</h3>
        <p>{{ photo.timestamp | formatDate }}</p>
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
  import api from '@/services/axios'; // Assuming axios is set up to handle API calls
  
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
        if (this.isLiked) {
          this.photo.likesCount++;
          await api.post(`/photos/likes/${this.photo.photoId}`);
        } else {
          this.photo.likesCount--;
          await api.delete(`/photos/likes/${this.photo.photoId}`);
        }
        this.isLiked = !this.isLiked;
      },
      toggleComments() {
        this.showComments = !this.showComments;
      },
      async postComment() {
        if (this.newComment.trim() !== '') {
          const response = await api.post(`/photos/comments/${this.photo.photoId}`, { content: this.newComment });
          this.photo.comments.push({ // Simulating adding the comment to the list
            username: "YourUsername", // This should ideally come from the server response or a global state
            content: this.newComment,
            commentId: response.data.commentId // Assuming your server returns the new comment ID
          });
          this.newComment = ''; // Reset the comment input field
        }
      }
    },
    filters: {
      formatDate(value) {
        return new Date(value).toLocaleString();
      }
    }
  }
  </script>
  
  <style scoped>
  .comment-form {
    display: flex;
    justify-content: space-between;
    margin-top: 10px;
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
  