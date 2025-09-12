# Changelog

All notable changes to the Lunch Delivery System will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2025-09-12

### Added
- **Daily Menu Management Enhancements**
  - Added "Select All" button to select all menu items at once in daily menu form
  - Added "Unselect All" button to deselect all menu items in daily menu form
  - Enhanced user experience for managing large menu lists with bulk selection functionality

- **Order Security Features**
  - Added prevention of editing orders that have already been paid
  - Frontend protection with disabled "Edit Order" button and visual notice for paid orders
  - Backend validation in both orderForm and submitOrder handlers to prevent API bypass

### Fixed
- **Template Error Handling**
  - Fixed nil pointer error in daily menu form when no existing menu is present
  - Added proper nil checking for `.existing` before accessing MenuItemIDs
  - Resolved "Date and menu items required" error during daily menu saving

### Security
- **Order Integrity Protection**
  - Added server-side validation to prevent modification of paid orders
  - Enhanced security against bypassing frontend restrictions via direct API calls

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

[0.2.0]: https://github.com/miftahulmahfuzh/lunch-delivery/releases/tag/v0.2.0
[0.1.0]: https://github.com/miftahulmahfuzh/lunch-delivery/releases/tag/v0.1.0
