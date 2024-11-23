
# **CRM System**

A robust and scalable Customer Relationship Management (CRM) system built using **Golang**. This project features modular architecture, microservices, and secure data management, designed to enhance business operations by streamlining customer interactions.

---

## **Features**
- **User Authentication**: OAuth2 and JWT-based secure authentication mechanisms.
- **RESTful APIs and gRPC Services**: For seamless communication between microservices.
- **Scalable Architecture**: Modular design for flexibility and future growth.
- **AWS Integration**: Deployed on AWS using EC2, S3, and RDS with load balancing and auto-scaling.
- **Database Optimization**: Efficient queries and ORM tools for data management.
- **Continuous Integration/Deployment**: Automated pipelines to streamline deployments.

---

## **Tech Stack**
- **Backend**: Golang  
- **Cloud**: AWS (EC2, S3, RDS)  
- **Authentication**: OAuth2, JWT  
- **Database**: PostgreSQL, ORM tools  
- **Microservices**: gRPC, REST APIs  
- **DevOps**: CI/CD pipelines  

---

## **Installation**
1. Clone the repository:
   ```bash
   git clone https://github.com/tamimorif/CRM-Service.git
   cd crm-system
   ```
2. Set up environment variables:
   - Create a `.env` file with required variables for database connection, AWS credentials, etc.
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run the application:
   ```bash
   go run main.go
   ```

---

## **Project Structure**
```plaintext
.
├── authMiddleware.go    # Middleware for authentication
├── database.go          # Database connection and configuration
├── handlers.go          # API request handlers
├── models/              # Data models
├── services/            # Business logic
├── routes/              # API routes
└── main.go              # Entry point of the application
```

---

## **Usage**
1. **Authentication**:
   - Register and log in using secure token-based authentication (JWT).
2. **Manage Users**:
   - Add, edit, and delete customer profiles.
3. **Track Interactions**:
   - Log and manage customer communication data.
4. **Admin Features**:
   - Role-based access control for different user levels.

---

## **Contributing**
Contributions are welcome! Follow these steps:
1. Fork the repository.
2. Create a feature branch:
   ```bash
   git checkout -b feature-name
   ```
3. Commit changes and push:
   ```bash
   git commit -m "Add your message here"
   git push origin feature-name
   ```
4. Create a Pull Request.

---

## **License**
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

---

## **Contact**
**Author**: Tamim Orif  
**GitHub**: [tamimorif](https://github.com/tamimorif)  
Feel free to reach out for questions or collaborations!
