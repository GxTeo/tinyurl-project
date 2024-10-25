# TinyURL Project
A URL shortening service built with Golang.

## Tech Stack
- **Golang**: The backend programming language used for building the application.
- **React**: The frontend framework used for building the application
- **PostgreSQL**: A powerful, open-source relational database used to store the URLs and their mappings.
- **Redis**: An in-memory data structure store used for caching and improving the performance of URL lookups.
- **Docker**: Used for containerization, allowing the application and its dependencies to run consistently across different environments.
- **Docker Compose**: A tool for defining and running multi-container Docker applications, making it easier to manage the services.

## Setup Instructions with Docker Compose

To set up the TinyURL project using Docker Compose, follow these steps:

1. **Clone the Repository**
   ```bash
   git clone https://github.com/yourusername/tinyurl-project.git
   cd tinyurl-project
   ```

2. **Install Docker and Docker Compose**
   Ensure you have Docker and Docker Compose installed on your machine. You can download them from the [Docker website](https://www.docker.com/get-started).

3. **Start the Application**
   Use Docker Compose to build and start the application along with the PostgreSQL database:
   ```bash
   docker-compose up --build
   ```

## Contributing
Feel free to submit issues for any bugs you see. Thank you!
