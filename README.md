# Data Drive System

A backend service similar to Google Drive, built using **Golang** with the **Gin Web Framework** to manage files and folders with JWT-based user authentication.

## 📌 Features

- User Authentication (Signup, Login) with JWT
- File Upload, Retrieval, and Soft Deletion
- Folder Creation, Retrieval (Nested Structure), and Soft Deletion
- Password Encryption using **bcrypt**
- Hierarchical Folder Structure
- Logging using **Lumberjack**
- Environment Variable Management with **Viper**
- Soft Deletion with `delete_flag` for Data Integrity

## 🔑 Technologies Used

- **Go (Golang)**
- **Gin Framework**
- **GORM** (ORM)
- **MySQL**
- **JWT (github.com/dgrijalva/jwt-go)**
- **bcrypt (golang.org/x/crypto/bcrypt)**
- **Viper** (Configuration Management)
- **Lumberjack** (Logging)

## Folder Structure

```bash
├── Documents/        # Documentation Files
├── config/          # Environment Configurations
├── dao/             # Data Access Layer
├── entities/        # Data Entities/Models
├── handlers/        # API Handlers
├── log/             # Log Files
├── logger/          # Logging Utilities
├── models/          # Database Models
├── router/          # API Routes
├── utility/         # Utility Functions
├── .gitignore       # Git Ignore File
├── README.md        # Project Documentation
├── go.mod           # Go Module File
├── go.sum           # Go Dependencies File
└── main.go          # Application Entry Point
```

## 🔐 User Authentication APIs

| Endpoint  | Method | Description                            |
| --------- | ------ | -------------------------------------- |
| `/signup` | POST   | Create a new user                      |
| `/login`  | POST   | Authenticate user and return JWT token |

## 📄 File Management APIs

| Endpoint              | Method | Description         |
| --------------------- | ------ | ------------------- |
| `/uploadFile`         | POST   | Upload a file       |
| `/getFilesByUserId`   | GET    | Retrieve user files |
| `/deleteFileByUserId` | DELETE | Soft delete file    |

## 📂 Folder Management APIs

| Endpoint                | Method | Description               |
| ----------------------- | ------ | ------------------------- |
| `/createFolderByUserId` | POST   | Create Folder             |
| `/getFolders`           | GET    | Retrieve Folder Hierarchy |
| `/deleteFolderByUserId` | DELETE | Soft Delete Folder        |

## 🔑 JWT Authentication

All API requests (except **/signup** and **/login**) require a JWT token in the request headers.

### Example Header

```bash
Authorization: Bearer <your_token_here>
```

## 🔧 Environment Variables

Create a `.env` file in the root directory with the following keys:

```env
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=datadrive
JWT_SECRET=your_jwt_secret
LOG_FILE=log/logfile.log
```

## 🚀 How to Run the Project

1. Clone the Repository

```bash
git clone <https://github.com/debajyatibanerjee0002/Data-Drive-System_GDrive.git>
cd datadrive
```

2. Install Dependencies

```bash
go mod tidy
```

3. Configure `.env` file

4. Run the Project

```bash
go run main.go
```

### API Testing

You can use **Postman** or **Thunder Client** to test the APIs.

## 📌 Logs

Logs will be stored in `log/logfile.log` using **Lumberjack**.

## 🛠️ Database

- MySQL is used as the primary database.
- Database migrations are handled using GORM.

## 🧠 Improvements

- Add UI (Frontend)
- File Download API
- Search Feature
- Pagination for large datasets

## 📄 Documentation

Refer to the **Documents** folder for detailed API documentation.

## 🤝 Contributing

Feel free to contribute by submitting a pull request or reporting bugs.

## 📌 Author

**Debajyati Banerjee**  
[LinkedIn Profile](https://www.linkedin.com/in/debajyatibanerjee0002/)  
[GitHub Profile](https://github.com/debajyatibanerjee0002/Data-Drive-System_GDrive)
