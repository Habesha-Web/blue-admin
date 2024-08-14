# Blue API Role Management System

Welcome to the **Blue API Role Management System** documentation. This document provides an overview of the features and endpoints available for managing users, roles, features, endpoints, pages, and applications within the system.

## Features Overview

### Authentication
- **login**: Authenticate a user and start a session.
- **check_login**: Verify if the user is currently logged in.

### User Management
- **Create User**: Add a new user to the system.
- **Get Users List**: Retrieve a list of all users.
- **Update User Detail**: Modify details of an existing user.
- **Delete User**: Remove a user from the system.
- **Activate/Deactivate User**: Enable or disable a user's account.
- **Get User Roles**: List all roles assigned to a user.
- **Update User Role**: Change the role assigned to a user.
- **Delete User Role**: Remove a role from a user.

### Role Management
- **Create Role**: Define a new role within the system.
- **Update Role**: Modify details of an existing role.
- **Delete Role**: Remove a role from the system.
- **Activate/Deactivate Role**: Enable or disable a role.

### Features Management
- **Create Feature**: Define a new feature.
- **Update Feature Details**: Modify details of an existing feature.
- **Activate/Deactivate Feature**: Enable or disable a feature.
- **Delete Feature**: Remove a feature from the system.
- **Map Feature with Role**: Associate a feature with a specific role.
- **Map Feature with Endpoints**: Associate a feature with specific endpoints.

### Endpoints
- **Endpoint (Auto Populate for Self)**: Automatically populate endpoints for internal use.
- **Auto Populate (gRPC Endpoint for External Apps)**: Automatically populate gRPC endpoints for external applications.
- **Get List of Endpoints**: Retrieve a list of all endpoints.

### Page Management
- **Create Page**: Define a new page.
- **Activate/Deactivate Page**: Enable or disable a page.
- **Map Features with Page**: Associate features with a specific page.

### Application Management
- **Create App**: Define a new application.
- **Activate/Deactivate App**: Enable or disable an application.
- **Map Features with App**: Associate features with a specific application.

## API Reference

### Authentication
- **POST /login**: Logs in a user.
- **GET /check_login**: Checks the login status of the user.

### Users
- **POST /users**: Create a new user.
- **GET /users**: Get a list of all users.
- **PUT /users/{id}**: Update a user's details.
- **DELETE /users/{id}**: Delete a user.
- **PATCH /users/{id}/activate**: Activate a user.
- **PATCH /users/{id}/deactivate**: Deactivate a user.
- **GET /users/{id}/roles**: Get roles assigned to a user.
- **PUT /users/{id}/roles**: Update a user's roles.
- **DELETE /users/{id}/roles/{role_id}**: Delete a user's role.

### Roles
- **POST /roles**: Create a new role.
- **PUT /roles/{id}**: Update a role.
- **DELETE /roles/{id}**: Delete a role.
- **PATCH /roles/{id}/activate**: Activate a role.
- **PATCH /roles/{id}/deactivate**: Deactivate a role.

### Features
- **POST /features**: Create a new feature.
- **PUT /features/{id}**: Update a feature's details.
- **DELETE /features/{id}**: Delete a feature.
- **PATCH /features/{id}/activate**: Activate a feature.
- **PATCH /features/{id}/deactivate**: Deactivate a feature.
- **POST /features/{id}/roles**: Map a feature with a role.
- **POST /features/{id}/endpoints**: Map a feature with endpoints.

### Endpoints
- **GET /endpoints**: Get a list of endpoints.
- **POST /endpoints/self**: Auto populate endpoints for internal use.
- **POST /endpoints/grpc**: Auto populate gRPC endpoints for external applications.

### Pages
- **POST /pages**: Create a new page.
- **PATCH /pages/{id}/activate**: Activate a page.
- **PATCH /pages/{id}/deactivate**: Deactivate a page.
- **POST /pages/{id}/features**: Map features with a page.

### Applications
- **POST /apps**: Create a new application.
- **PATCH /apps/{id}/activate**: Activate an application.
- **PATCH /apps/{id}/deactivate**: Deactivate an application.
- **POST /apps/{id}/features**: Map features with an application.

## Getting Started

1. **Authentication**: Start by logging in to the system to obtain the necessary session token.
2. **User Management**: Create and manage users according to your needs.
3. **Role Management**: Define and adjust roles for users.
4. **Features Management**: Add and configure features, and map them to roles and endpoints.
5. **Endpoints**: Manage internal and external endpoints as required.
6. **Page Management**: Create and configure pages, and associate features.
7. **Application Management**: Define and manage applications, and map features.

For further details or assistance, please refer to the API documentation or contact support.
