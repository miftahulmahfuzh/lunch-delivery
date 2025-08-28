# Changelog

All notable changes to the Lunch Delivery System will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-08-28

### Added
- **Core System Architecture**
  - Go backend with Gin framework
  - PostgreSQL database with sqlx
  - Server-side rendered HTML templates

- **Admin Management Panel**
  - Menu items CRUD (create, read, update, delete)
  - Daily menu selection from master menu
  - Company management with full CRUD operations
  - Employee management per company
  - Order sessions management with date filtering
  - Real-time order tracking and payment status

- **Employee Self-Service Portal**
  - User registration with company selection
  - Login/logout with cookie-based sessions
  - Personal dashboard showing order history
  - Interactive order placement with real-time price calculation
  - Order modification capability until session closes

- **Order Management System**
  - Date-based order sessions per company
  - Session status lifecycle (OPEN → CLOSED_FOR_ORDERS → DELIVERED → PAYMENT_PENDING → COMPLETED)
  - Session close/reopen functionality for admin flexibility
  - Individual order tracking with payment status
  - Order summary with revenue calculation

- **Database Schema**
  - 6 core tables with proper relationships
  - PostgreSQL arrays for menu item selections
  - Soft delete functionality
  - Unique constraints preventing duplicate sessions

- **Security Features**
  - bcrypt password hashing
  - HTTP-only cookie sessions
  - Input validation and SQL injection prevention
  - Authentication middleware

- **User Experience**
  - Responsive design for desktop and mobile
  - Real-time order total calculation
  - Visual feedback for selected menu items
  - Comprehensive error handling
  - Intuitive admin navigation

### Technical Implementation
- **25+ API endpoints** covering all functionality
- **11 HTML templates** with standalone designs
- **Repository pattern** for clean data access
- **Middleware-based authentication**
- **Template functions** for data formatting
- **Sample data scripts** for testing

### Business Logic
- **B2B lunch delivery model** - companies order for employees
- **Daily menu subset** selection from master catalog
- **Time-bounded ordering** with admin control
- **Individual payment tracking** per employee
- **WhatsApp contact integration** for communication
- **Flexible session management** with reopen capability

[0.1.0]: https://github.com/miftahulmahfuzh/lunch-delivery/releases/tag/v0.1.0
