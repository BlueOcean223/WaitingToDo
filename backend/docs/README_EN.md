# WaitingToDo

English | [中文](../../README.md)

A feature-rich modern todo management platform that supports personal task management, team collaboration, friend system, and many other functions.

## 📖 Project Introduction

WaitingToDo is a full-stack todo management application developed with Vue 3 + Go. It not only provides traditional personal task management features but also integrates modern functionalities such as team collaboration, friend system, real-time message notifications, and file management, aiming to provide users with an efficient and convenient task management solution.

## ✨ Key Features

### 🔐 User System
- User registration, login, password reset
- JWT authentication
- User profile management
- Avatar upload

### 📝 Task Management
- Create, edit, delete tasks
- Task priority settings (High/Medium/Low)
- Task status management (Pending/In Progress/Completed)
- Task deadline reminders
- Task tag categorization
- Task search and filtering

### 👥 Team Collaboration
- Create and manage teams
- Team task assignment
- Join teams with invite codes
- Team member management
- Collaborative task progress tracking

### 🤝 Friend System
- Add friends
- Friend request management
- Friend task sharing
- Online status display

### 🔔 Message Notifications
- Real-time message push
- Task reminder notifications
- Team collaboration notifications
- Friend request notifications
- System messages

### 📁 File Management
- Task attachment upload
- File preview
- File download
- Storage space management

### 📊 Data Statistics
- Task completion rate statistics
- Personal efficiency analysis
- Team collaboration data
- Visual chart display

### 🎨 User Experience
- Responsive design, mobile support
- Dark/Light theme switching
- Internationalization support (Chinese/English)
- Intuitive user interface

## 🚀 Quick Start Guide

### Environment Requirements

#### Backend Environment
- Go 1.19+
- MySQL 8.0+
- Redis 6.0+
- MinIO (Object Storage)
- RabbitMQ (Message Queue)

#### Frontend Environment
- Node.js 16+
- npm 8+ or yarn 1.22+

### Installation Steps

#### 1. Clone the Project
```bash
git clone https://github.com/yourusername/WaitingToDo.git
cd WaitingToDo
```

#### 2. Backend Setup

```bash
# Enter backend directory
cd backend

# Install dependencies
go mod download

# Copy configuration file
cp config/config.example.yaml config/config.yaml

# Edit configuration file, set database connection and other information
vim config/config.yaml

# Run database migration
go run main.go migrate

# Start backend service
go run main.go
```

#### 3. Frontend Setup

```bash
# Enter frontend directory
cd frontend

# Install dependencies
npm install
# or use yarn
yarn install

# Start development server
npm run dev
# or use yarn
yarn dev
```

#### 4. Database Setup

```sql
-- Create database
CREATE DATABASE waitingtodo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create user (optional)
CREATE USER 'waitingtodo'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON waitingtodo.* TO 'waitingtodo'@'localhost';
FLUSH PRIVILEGES;
```

#### 5. Using Docker (Recommended)

```bash
# Start with Docker Compose one-click
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f
```

### Access the Application

- Frontend Application: http://localhost:3000
- Backend API: http://localhost:8080
- API Documentation: http://localhost:8080/swagger/index.html

## 📚 Usage Instructions

### Basic Usage Flow

1. **Register Account**: Visit the application homepage, click the register button to create a new account
2. **Login System**: Use the registered email and password to login
3. **Create Tasks**: Click the "Add Task" button, fill in task information
4. **Manage Tasks**: View, edit, complete, or delete tasks in the task list
5. **Team Collaboration**: Create teams, invite members, assign team tasks
6. **Add Friends**: Search users, send friend requests, share tasks with friends

### Advanced Features

- **Task Filtering**: Filter tasks using status, priority, tags, and other conditions
- **Batch Operations**: Select multiple tasks for batch deletion or status updates
- **Data Export**: Export task data in CSV or PDF format
- **API Integration**: Integrate with third-party applications using RESTful API

### Mobile Usage

The application uses responsive design to provide a good user experience on mobile devices:

- Touch-friendly interface design
- Gesture operation support
- Offline data caching
- Push notification support

## 🛠️ Technology Stack

### Frontend Technologies
- **Framework**: Vue 3 + TypeScript
- **Build Tool**: Vite
- **State Management**: Pinia
- **Routing**: Vue Router 4
- **UI Components**: Element Plus
- **Styling**: SCSS
- **HTTP Client**: Axios
- **Charts**: ECharts

### Backend Technologies
- **Language**: Go 1.19+
- **Framework**: Gin
- **Database**: MySQL 8.0
- **Cache**: Redis
- **Object Storage**: MinIO
- **Message Queue**: RabbitMQ
- **Authentication**: JWT
- **API Documentation**: Swagger

### Development Tools
- **Version Control**: Git
- **Containerization**: Docker + Docker Compose
- **Code Quality**: ESLint + Prettier (Frontend), golangci-lint (Backend)
- **Testing**: Jest (Frontend), Go testing (Backend)

## 📁 Project Structure

```
WaitingToDo/
├── frontend/                 # Frontend project
│   ├── src/
│   │   ├── components/       # Components
│   │   ├── views/           # Pages
│   │   ├── store/           # State management
│   │   ├── router/          # Router configuration
│   │   ├── api/             # API interfaces
│   │   └── utils/           # Utility functions
│   ├── public/              # Static resources
│   └── package.json         # Dependency configuration
├── backend/                 # Backend project
│   ├── api/                 # API routes
│   ├── config/              # Configuration files
│   ├── internal/            # Internal modules
│   │   ├── handler/         # Handlers
│   │   ├── service/         # Business logic
│   │   ├── repository/      # Data access
│   │   └── model/           # Data models
│   ├── pkg/                 # Common packages
│   └── main.go              # Entry file
├── docs/                    # Project documentation
├── docker-compose.yml       # Docker orchestration
└── README.md               # Project description
```

## 🤝 Contributing Guide

We welcome all forms of contributions! Whether it's reporting bugs, suggesting new features, or submitting code improvements.

### How to Contribute

1. **Fork the Project**
   ```bash
   # Click the Fork button in the upper right corner of the GitHub page
   ```

2. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

4. **Push to Branch**
   ```bash
   git push origin feature/your-feature-name
   ```

5. **Create Pull Request**
   - Create a Pull Request on GitHub
   - Describe your changes in detail
   - Wait for code review

### Code Standards

#### Frontend Code Standards
- Use ESLint + Prettier for code formatting
- Follow Vue 3 Composition API best practices
- Component naming uses PascalCase
- File naming uses kebab-case

#### Backend Code Standards
- Follow Go official code standards
- Use golangci-lint for code checking
- Function and variable naming uses camelCase
- Package names use lowercase letters

#### Commit Message Standards

Use [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
type(scope): description

[optional body]

[optional footer]
```

Type descriptions:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation update
- `style`: Code format adjustment
- `refactor`: Code refactoring
- `test`: Test related
- `chore`: Build process or auxiliary tool changes

### Reporting Issues

If you find bugs or have feature suggestions, please:

1. Check if there are related issues in [Issues](https://github.com/yourusername/WaitingToDo/issues)
2. If not, create a new Issue
3. Describe the problem or suggestion in detail
4. If it's a bug, please provide reproduction steps

### Development Environment Setup

1. **Install Development Tools**
   ```bash
   # Install frontend development tools
   npm install -g @vue/cli
   npm install -g eslint prettier
   
   # Install backend development tools
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. **Configure Git Hooks**
   ```bash
   # Install husky
   npm install -g husky
   
   # Set pre-commit hook
   husky add .husky/pre-commit "npm run lint"
   ```

3. **Run Tests**
   ```bash
   # Frontend tests
   cd frontend
   npm run test
   
   # Backend tests
   cd backend
   go test ./...
   ```

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Thanks to all developers and users who have contributed to this project!

## 📞 Contact Us

- Project Homepage: https://github.com/yourusername/WaitingToDo
- Issue Feedback: https://github.com/yourusername/WaitingToDo/issues
- Email: your-email@example.com

---

⭐ If this project helps you, please give us a Star!