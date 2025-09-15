# Changelog

All notable changes to the Lunch Delivery System will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.0] - 2025-09-15

### Added
- **🔄 Comprehensive CI/CD Pipeline Implementation**
  - Complete GitHub Actions workflow for automated testing and CI/CD
  - Automated unit test execution with coverage reporting
  - Go formatting checks with gofmt validation
  - golangci-lint integration for code quality enforcement
  - Multi-step build verification process

- **🧪 Comprehensive Unit Testing Suite**
  - Complete test coverage for entire internal codebase
  - Unit tests for all handlers (auth, admin, orders)
  - Comprehensive middleware testing with authentication scenarios
  - Database layer testing with mock implementations
  - LLM client testing with mock interfaces
  - Services layer testing for nutritionist functionality
  - Models and repository testing with fixtures
  - Utility functions testing (email, token, validation)

- **📚 Extensive Testing Documentation**
  - Master unittest guide (UNITTEST_MASTER_GUIDE.md) with comprehensive testing strategies
  - Individual unittest guides for each internal package
  - Detailed testing conventions and best practices
  - Mock implementation patterns and usage guidelines
  - Test fixture management and common utilities

- **🏗️ Testing Infrastructure & Utilities**
  - Centralized test utilities and helpers
  - Mock implementations for LLM and repository interfaces
  - Test fixtures for consistent data setup
  - Common testing patterns and shared utilities
  - Enhanced test database management

### Fixed
- **🔒 Security Vulnerabilities Resolution**
  - Fixed all gosec security vulnerabilities
  - Improved security posture across the codebase
  - Enhanced input validation and sanitization
  - Secure error handling improvements

- **🔧 Code Quality & Compilation Issues**
  - Resolved all golangci-lint issues and warnings
  - Fixed circular import problems in test packages
  - Resolved all unit test failures
  - Fixed CI build compilation errors
  - Enhanced error handling patterns

- **🤖 AI Nutritionist Improvements**
  - Fixed AI nutritionist caching null constraint violations
  - Resolved repository method implementation gaps
  - Enhanced nutritionist service reliability

### Changed
- **🎨 Template System Refactoring**
  - Migrated all admin templates to shared styles system
  - Consolidated template styling with reusable components
  - Extracted shared styles into dedicated template
  - Eliminated styling redundancy across templates

- **📁 Project Structure Improvements**
  - Removed duplicate favicon files
  - Enhanced project organization and file structure
  - Improved import paths and dependency management

- **⚙️ CI/CD Infrastructure Updates**
  - Updated deprecated GitHub Actions (upload-artifact v3 → v4)
  - Fixed gosec action configuration with specific version v2.22.8
  - Replaced invalid security scanning actions with proper alternatives
  - Enhanced CI pipeline reliability and performance

### Security
- **🛡️ Enhanced Security Scanning**
  - Integration of gosec static security analyzer
  - Automated security vulnerability detection in CI
  - Enhanced security best practices enforcement
  - Improved secret and credential handling

### Technical Infrastructure
- **🧪 Testing Framework Enhancement**
  - Comprehensive mock system for external dependencies
  - Advanced testing utilities for database operations
  - Enhanced test isolation and cleanup mechanisms
  - Improved test performance and reliability

- **📈 Code Quality Improvements**
  - Standardized coding conventions across entire codebase
  - Enhanced error handling and logging practices
  - Improved code documentation and comments
  - Better separation of concerns in testing

## [0.5.0] - 2025-01-14

### Added
- **🔐 Complete Password Reset System**
  - Forgot password functionality with email-based token system
  - Secure password reset tokens with expiration (1 hour)
  - SMTP email integration with configurable settings
  - Dedicated forgot password and reset password templates
  - Comprehensive email utilities and token management

- **📊 Individual Order Tracking Enhancement**
  - Individual order preparation status tracking system
  - Real-time order status updates with visual indicators
  - Enhanced order management with preparation workflow
  - Status-based UI updates and notifications

- **🎯 Advanced Order Session Management**
  - Duplicate order session detection and prevention
  - Smart redirect functionality for existing sessions
  - Enhanced session validation with user-friendly popups
  - Improved session creation workflow with conflict handling

- **🎨 Enhanced User Interface & Navigation**
  - Clickable order summary with scroll-to-item functionality
  - Improved order navigation and user experience
  - Toggle switches replacing traditional action buttons
  - Refresh button integration for real-time updates
  - Priority-based Edit Orders link with session status awareness

- **🗂️ Complete SQL Scripts Reorganization**
  - Organized SQL scripts into logical directory structure
  - Separated schema, seeds, updates, and deletion scripts
  - Enhanced documentation for all script categories
  - Migration from legacy migrations folder to structured scripts/sql/

- **📧 SMTP Testing and Documentation**
  - Comprehensive SMTP testing utilities
  - Gmail integration setup documentation
  - Email testing scripts for password reset functionality
  - Complete SMTP configuration examples

- **🛡️ Enhanced Admin Features**
  - Menu validation for order session creation
  - Improved admin page footer consistency
  - Standardized admin interface across all pages
  - Enhanced menu management with validation checks

- **📖 Documentation Improvements**
  - Enhanced README with forgot password feature documentation
  - Comprehensive scripts directory documentation
  - Detailed setup instructions for email functionality
  - Updated seed data documentation

### Changed
- **🔄 File Structure Modernization**
  - Migrated from migrations/ to scripts/sql/ structure
  - Reorganized SQL files by purpose (schema/seeds/updates/deletions)
  - Updated .gitignore for improved project structure
  - Enhanced .env.example with email configuration

- **🎯 UI/UX Enhancements**
  - Replaced traditional buttons with modern toggle switches
  - Improved order status layout with better visual hierarchy
  - Enhanced session orders interface with status prioritization
  - Streamlined order management interface

- **📊 Order Management Improvements**
  - Enhanced order preparation tracking workflow
  - Improved session status handling and validation
  - Better integration between order creation and session management
  - Enhanced order summary display with navigation features

### Fixed
- **🔧 Session Management Issues**
  - Fixed duplicate order session creation problems
  - Resolved session validation edge cases
  - Improved error handling for session conflicts
  - Enhanced session status update reliability

- **🎨 UI Consistency & Layout**
  - Standardized footer layouts across admin pages
  - Improved toggle switch functionality and appearance
  - Fixed order status display inconsistencies
  - Enhanced responsive design elements

### Security
- **🔐 Authentication System Enhancements**
  - Secure password reset token implementation
  - Time-limited token validation (1-hour expiration)
  - Email-based verification for password resets
  - Enhanced input validation for authentication flows

### Technical Infrastructure
- **📧 Email System Integration**
  - SMTP client implementation with configurable settings
  - Token generation and validation utilities
  - Email template system for password reset notifications
  - Comprehensive error handling for email operations

- **🗄️ Database Schema Updates**
  - Password reset tokens table implementation
  - Individual order status tracking enhancement
  - Updated seed data with realistic menu items
  - Enhanced data validation and constraints

## [0.4.0] - 2025-09-13

### Added
- **🎨 Modern UI/UX Complete Redesign**
  - Comprehensive visual overhaul with Tailwind CSS design system
  - Chrome-style tab navigation across all templates
  - Glassmorphism effects and streamlined UX patterns
  - Smart floating submit button with improved positioning
  - Modern admin panel with enhanced visual hierarchy

- **🧠 Enhanced AI Nutritionist Features**
  - Improved user experience with better positioning and close functionality
  - User-specific stock filtering integration
  - Enhanced floating button UX with AI nutritionist integration
  - Real-time UI rendering improvements

- **📋 Advanced Order Management**
  - Enhanced order form with search functionality
  - Select All and Unselect All buttons for bulk operations
  - Date range filtering for order history
  - Enhanced order session management with status controls
  - Smart out-of-stock handling improvements

- **📊 Comprehensive Stock Management System**
  - User-specific stock tracking and filtering
  - Admin controls for stock management
  - Real-time stock synchronization
  - Notification system for stock alerts
  - Global stock management removal in favor of user-specific tracking

- **🔔 Advanced Notification System**
  - Modal interface for notification management
  - Bulk notification operations
  - Enhanced notification management workflow
  - Real-time notification updates

- **💰 Currency System Updates**
  - Complete conversion from US cents to Indonesian Rupiah
  - Updated price representation across all components
  - Migration scripts for currency conversion
  - Menu price statistics and reporting

- **🛠️ Development Infrastructure**
  - Comprehensive .gitignore configuration
  - Environment configuration examples
  - Enhanced migration system with test data
  - Favicon and branding assets

### Fixed
- **🎯 UI Positioning and Layout**
  - Multiple floating button positioning fixes
  - Chrome tab navigation design improvements
  - Header blur effects optimization
  - Sticky positioning and layout adjustments
  - Template routing bug resolution (admin/daily-menu)

- **🔧 Technical Improvements**
  - SQL LIMIT formatting corrections in notifications
  - Enhanced error handling and validation
  - Improved modal dialog functionality
  - Better responsive design implementations

### Changed
- **🎨 Design System Standardization**
  - Standardized chrome tab navigation with green color scheme
  - Simplified navigation design patterns
  - Streamlined admin operations (removed confirmation dialogs)
  - Enhanced edit functionality with modal dialogs

- **📈 Performance Optimizations**
  - Dynamic preview functionality in daily menu forms
  - Improved search functionality across forms
  - Better data loading and caching strategies
  - Enhanced template rendering performance

### Security
- **🔒 Enhanced Data Validation**
  - Improved input validation across all forms
  - Better error handling for edge cases
  - Enhanced SQL injection prevention
  - Secure file handling improvements

## [0.3.0] - 2025-09-12

### Added
- **🧠 AI Nutritionist Selection Feature**
  - Smart auto-selection button powered by DeepSeek AI
  - Structured JSON responses for reliable menu item selection
  - Daily caching system for optimal performance and cost efficiency
  - Real-time nutritional analysis and reasoning display
  - Auto-selection of menu items based on AI recommendations

- **📊 Daily Cache Optimization System**
  - Single LLM API call per day across all users for same menu
  - PostgreSQL-based caching with nutritionist_selections table
  - Menu change invalidation with automatic cache clearing
  - Cost-effective scaling for high user volumes

- **🔄 Menu Update Reset Mechanism**
  - Admin menu updates trigger nutritionist_reset flag
  - User notification system for menu changes affecting AI selections
  - Tracking of users who utilized nutritionist recommendations
  - Support for multiple daily menu updates with proper cache invalidation

- **🎨 Enhanced User Experience**
  - Beautiful teal "Nutritionist Selection" button with loading states
  - Detailed AI reasoning display with nutritional breakdown
  - Visual feedback for protein, vegetables, and carbohydrates assessment
  - Warning notifications for users when menus are updated after AI selection

- **🏗️ New Technical Infrastructure**
  - LLM integration layer with DeepSeek API support
  - Configuration management system for AI credentials
  - Nutritionist service with structured prompt engineering
  - User selection tracking for notification system
  - Database migrations for new schema additions

### Technical Implementation
- **New API Endpoints**: `POST /order/:company/:date/nutritionist-select`
- **Database Schema**: 2 new tables (nutritionist_selections, nutritionist_user_selections) + 1 new column
- **New Dependencies**: langchaingo, tiktoken-go, zerolog for LLM integration
- **JavaScript Frontend**: Auto-selection logic and API integration
- **Error Handling**: JSON parsing with regex fallback for LLM responses
- **Validation**: Index bounds checking and menu item validation

### Performance & Architecture
- **Daily Cache Strategy**: Dramatically reduces LLM API costs by sharing results
- **Structured Responses**: Eliminates text parsing errors with JSON schema enforcement
- **Reset Flag System**: Handles dynamic menu updates throughout the day
- **User Tracking**: Enables targeted notifications for affected users
- **Graceful Degradation**: Fallback mechanisms for LLM service interruptions

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

[0.6.0]: https://github.com/miftahulmahfuzh/lunch-delivery/releases/tag/v0.6.0
[0.5.0]: https://github.com/miftahulmahfuzh/lunch-delivery/releases/tag/v0.5.0
[0.4.0]: https://github.com/miftahulmahfuzh/lunch-delivery/releases/tag/v0.4.0
[0.3.0]: https://github.com/miftahulmahfuzh/lunch-delivery/releases/tag/v0.3.0
[0.2.0]: https://github.com/miftahulmahfuzh/lunch-delivery/releases/tag/v0.2.0
[0.1.0]: https://github.com/miftahulmahfuzh/lunch-delivery/releases/tag/v0.1.0
