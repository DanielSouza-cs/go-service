# Go PDF Report Service

A Go microservice that generates PDF reports for students by consuming a Node.js backend API.

---

## Prerequisites

* **Go 1.23.0** or later.
* A running instance of the required **Node.js backend API**.

---

## Configuration

This project uses a `.env` file for configuration. Create a file named `.env` in the `go-service` root directory and fill it with the correct values for your environment.

**Example `.env` file:**
    ```env
    # Port for this Go service
    PORT=8080

    # URL for the Node.js data API
    NODE_API_URL=http://localhost:5007/api/v1

    # URL for the Node.js authentication endpoint
    NODE_AUTH_URL=http://localhost:5007/api/v1

    # Credentials to log into the Node.js API
    NODE_AUTH_EMAIL=admin@school-admin.com
    NODE_AUTH_PASSWORD=3OU4zn3q6Zh9

    # Application environment
    ENVIRONMENT=development

    # Log level
    LOG_LEVEL=info

## 🚀 Running the Service 

* **make run**:   to build and run the application.

* **make build**: to Compile the application into the /bin directory.

* **make tidy**:  to tidie and formats module dependencies.

* **make clean**: to remove the build directory.