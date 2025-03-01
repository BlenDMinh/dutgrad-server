# Da Nang University Graduation Project - GO Server

This is the **backend server** for the graduation project at Da Nang University. The server acts as the **main hub**, connecting multiple microservices (e.g., AI services, databases, etc.) and handling communication between them, including requests from the frontend.

## Table of Contents
- [Project Overview](#project-overview)
- [Tech Stack](#tech-stack)
- [Features](#features)
- [Setup & Installation](#setup--installation)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Contributing](#contributing)
- [License](#license)

## Project Overview

The server provides an interface for various services used in the project, such as:
- Microservices (e.g., AI service)
- Databases (e.g., PostgreSQL, MongoDB)
- Frontend connections (React, Angular, etc.)

This project demonstrates the **microservices architecture** and **RESTful API communication** between services. The server is designed to be **highly scalable** and **performant**, using **Go (Golang)** as the main programming language.

## Tech Stack

- **Language**: Go (Golang)
- **Web Framework**: Gin/Echo (Choose your preferred framework)
- **API**: RESTful APIs (with possible gRPC integration)
- **Database**: PostgreSQL / MongoDB (depending on your data needs)
- **Containerization**: Docker
- **Orchestration**: Kubernetes (optional, for large-scale deployments)
- **Version Control**: Git (GitHub/GitLab)
- **CI/CD**: GitHub Actions / GitLab CI / Jenkins

## Features

- **Microservices Integration**: Connects to different backend microservices (e.g., AI server).
- **REST API**: Provides a simple and robust API for communication with the frontend.
- **Database Management**: Integrates with relational and/or NoSQL databases.
- **Scalability**: Supports horizontal scaling with minimal overhead.
- **Error Handling & Logging**: Built-in error handling and logging for smooth operations.
- **Authentication**: Supports token-based authentication (JWT or OAuth).
- **Containerization**: Easily deployable via Docker for consistency across environments.

## Setup & Installation

### Prerequisites

Before setting up the project, make sure you have the following installed on your system:
- [Go](https://golang.org/doc/install)
- [Docker](https://www.docker.com/get-started)
- [PostgreSQL / MongoDB](https://www.postgresql.org/download/ or https://www.mongodb.com/try/download/community)
- [Git](https://git-scm.com/)

### Clone the Repository

Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/graduation-project-server.git
cd graduation-project-server
