# Real-Time Collaboration & Project Management API (Backend - Go)

## Overview
This backend is built using **Go** and **Gin** for a real-time collaboration and project management platform. It includes features for **task management, user authentication, AI-based project idea suggestions, chat, messaging, and real-time collaboration**.

---

## API Routes
### **Authentication Routes** (Public)
Base Path: `/auth`
- `POST /register` - Register a new user
- `POST /login` - User login

### **User Routes**
Base Path: `/user`
#### Public Routes:
- `GET /` - Get user notification preferences
- `PUT /` - Update user notification preferences

#### Protected Routes:
- `PUT /updateMain` - Update main user information
- `GET /secure_information` - Retrieve user security info
- `PUT /changePassword` - Change user password

---

### **Task Management Routes**
Base Path: `/tasks`
- `POST /addTask` - Add a new task
- `GET /` - Retrieve tasks by type
- `PUT /archive/:taskID` - Archive a task
- `PUT /edit/:taskID` - Edit an existing task
- `GET /search` - Search tasks
- `GET /archived` - Get archived tasks
- `PUT /restore/:taskID` - Restore archived tasks
- `DELETE /delete/:taskID` - Delete an archived task

---

### **AI-Based Project Idea Suggestion Routes**
Base Path: `/suggest`
- `GET /suggest` - Get AI-generated project ideas

---

### **Real-Time Collaboration (Publishing & Comments)**
Base Path: `/videos`
- `POST /publish` - Publish a video/project
- `PUT /edit/:publishID` - Edit a published video/project
- `PUT /like/:publishID` - Like a published video/project
- `PUT /dislike/:publishID` - Dislike a published video/project
- `GET /videos` - Get all published videos/projects
- `DELETE /delete/:publishID` - Delete a published video/project

#### Comments:
- `POST /comment` - Add a comment
- `PUT /comment/edit` - Edit a comment
- `DELETE /comment/delete/:publishID` - Delete a comment
- `GET /comment/:publishedID` - Get comments for a specific post
- `PUT /comment/like/:publishID` - Like a comment
- `PUT /comment/dislike/:publishID` - Dislike a comment

---

### **Group Collaboration Routes**
Base Path: `/group`
- `GET /group` - Retrieve group messages

---

### **Chat Routes**
Base Path: `/chats`
- `POST /chat` - Initiate a chat
- `GET /sessions/session` - Fetch specific chat session
- `DELETE /delete` - Delete a chat session
- `GET /sessions` - Get all chat sessions
- `POST /create` - Create a new chat session

---

### **Real-Time Messaging Routes**
Base Path: `/ws`
- `GET /ws` - WebSocket-based real-time messaging

---

## **Middleware**
- **JWT Authentication** - Protected routes require a valid JWT token.

---

## **Setup & Running the Server**
### **Prerequisites**
- Go 1.18+
- MongoDB
- Redis (for caching and real-time updates)

### **Installation**
```sh
# Clone the repository
git clone https://duressa2022/RealColab.git
cd your-repo

# Install dependencies
go mod tidy
```

### **Environment Variables (.env)**
```ini
ACCESS_TOKEN_SECRET=your_secret_key
MONGO_URI=mongodb://localhost:27017
REDIS_URI=redis://localhost:6379
```

### **Run the Server**
```sh
go run main.go
```

---

## **Contributing**
1. Fork the repository
2. Create a feature branch (`git checkout -b feature-name`)
3. Commit your changes (`git commit -m 'Add feature'`)
4. Push to the branch (`git push origin feature-name`)
5. Create a Pull Request

---

## **License**
MIT License

