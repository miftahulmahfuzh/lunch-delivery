# Lunch Delivery System

A comprehensive web application built in Go for managing lunch deliveries to corporate clients. The system handles menu management, employee registration, order sessions, and payment tracking.

## System Overview

This application manages a restaurant's lunch delivery service that operates on a B2B model, delivering packaged orders to company offices. Each company order contains individual meals for registered employees.

### Key Features
- **Company Management**: Register and manage corporate clients
- **Employee Self-Registration**: Company employees can create their own accounts
- **Daily Menu Management**: Admin can set available menu items for each day
- **Order Sessions**: Time-bound ordering windows for each company
- **Individual Order Tracking**: Track each employee's order and payment status
- **Real-time Order Management**: Close/reopen order sessions as needed
- **🤖 AI Nutritionist**: AI-powered meal recommendations with nutritional analysis and intelligent caching

## Architecture

### Technology Stack
- **Backend**: Go 1.21+ with Gin framework
- **Database**: PostgreSQL with sqlx for query handling
- **Frontend**: Server-side rendered HTML templates with vanilla CSS/JavaScript
- **Authentication**: Cookie-based sessions with bcrypt password hashing
- **AI Integration**: LLM-powered nutritionist service with smart caching

### Project Structure
```
lunch-delivery/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── models/
│   │   ├── models.go              # Data structures
│   │   └── repository.go          # Database operations
│   ├── handlers/
│   │   ├── handlers.go            # Route setup
│   │   ├── admin.go               # Admin functionality
│   │   ├── auth.go                # Authentication
│   │   ├── orders.go              # Customer orders
│   │   └── employees.go           # Employee management
│   ├── services/
│   │   └── nutritionist.go        # AI nutritionist service
│   ├── middleware/
│   │   └── auth.go                # Authentication middleware
│   └── database/
│       └── db.go                  # Database connection
├── migrations/
│   └── 001_initial.sql            # Database schema
├── templates/                     # HTML templates
└── static/                        # Static assets (CSS, JS, images)
```

## Database Schema

### Core Tables

**menu_items**
- Stores all available menu items with fixed prices
- Fields: id, name, price (in cents), active, created_at

**companies**
- Corporate clients who order lunch for their employees
- Fields: id, name, address, contact, active, created_at

**employees**
- Individual users who can place orders
- Fields: id, company_id, name, email, wa_contact, password_hash, active, created_at

**daily_menus**
- Subset of menu items available on specific dates
- Fields: id, date, menu_item_ids (array), created_at

**order_sessions**
- Time-bound ordering windows for companies on specific dates
- Fields: id, company_id, date, status, created_at, closed_at
- Status values: OPEN, CLOSED_FOR_ORDERS, DELIVERED, PAYMENT_PENDING, COMPLETED

**individual_orders**
- Individual employee orders within sessions
- Fields: id, session_id, employee_id, menu_item_ids (array), total_price, paid, created_at

**nutritionist_selections**
- AI-generated nutritional recommendations cached by date
- Fields: id, date, menu_item_ids (array), selected_indices (array), reasoning, nutritional_summary (JSONB), created_at

**nutritionist_user_selections**
- Tracks users who have used AI nutritionist recommendations
- Fields: id, employee_id, date, order_id, created_at

## Installation & Setup

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 12+
- Git

### Database Setup

1. **Install PostgreSQL and create database:**
```sql
CREATE DATABASE lunch_delivery;
CREATE USER lunch_user WITH PASSWORD '1234';
GRANT ALL PRIVILEGES ON DATABASE lunch_delivery TO lunch_user;
```

2. **Connect and run migrations:**
```sql
\c lunch_delivery

-- Create tables
CREATE TABLE menu_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    address TEXT,
    contact VARCHAR(255),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    company_id INTEGER REFERENCES companies(id),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    wa_contact VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE daily_menus (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    menu_item_ids INTEGER[] NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_sessions (
    id SERIAL PRIMARY KEY,
    company_id INTEGER REFERENCES companies(id),
    date DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'OPEN',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMP,
    UNIQUE(company_id, date)
);

CREATE TABLE individual_orders (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES order_sessions(id),
    employee_id INTEGER REFERENCES employees(id),
    menu_item_ids INTEGER[] NOT NULL,
    total_price INTEGER NOT NULL,
    paid BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(session_id, employee_id)
);

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO lunch_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO lunch_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO lunch_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO lunch_user;

-- Add nutritionist feature tables
CREATE TABLE nutritionist_selections (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    menu_item_ids BIGINT[] NOT NULL,
    selected_indices INTEGER[] NOT NULL,
    reasoning TEXT,
    nutritional_summary JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_nutritionist_selections_date ON nutritionist_selections(date);

-- Add reset flag to daily_menus table
ALTER TABLE daily_menus ADD COLUMN nutritionist_reset BOOLEAN DEFAULT FALSE;

-- Add tracking table for users who used nutritionist selection
CREATE TABLE nutritionist_user_selections (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER REFERENCES employees(id),
    date DATE NOT NULL,
    order_id INTEGER REFERENCES individual_orders(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employee_id, date)
);

CREATE INDEX idx_nutritionist_user_selections_date ON nutritionist_user_selections(date);
CREATE INDEX idx_nutritionist_user_selections_employee ON nutritionist_user_selections(employee_id);
```

### Application Setup

1. **Clone and install dependencies:**
```bash
git clone <repository-url>
cd lunch-delivery
go mod tidy
```

2. **Configure database connection in `cmd/server/main.go`:**
```go
db, err := database.NewConnection("localhost", "5432", "lunch_user", "1234", "lunch_delivery")
```

3. **Run the application:**
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### Sample Data

Insert test data to explore the application:

```sql
-- Sample menu items
INSERT INTO menu_items (name, price) VALUES
('Cah jagung muda', 1500),
('Cah labu', 1500),
('Cah toge', 1200),
('Cah kembang kol', 1500),
('Cah oyong telur', 1800),
('Ceker cabe ijo', 2000),
('Ayam goreng kandar merah', 2500),
('Udang crispy cabe garam', 3000);

-- Sample companies
INSERT INTO companies (name, address, contact) VALUES
('Tech Corp', 'Jakarta Selatan', 'tech@corp.com'),
('Marketing Inc', 'Jakarta Pusat', 'hello@marketing.com'),
('Finance Ltd', 'Jakarta Barat', 'contact@finance.com');

-- Sample employees (password: 'password')
INSERT INTO employees (company_id, name, email, wa_contact, password_hash) VALUES
(1, 'Jemmy', 'jemmy@techcorp.com', '+628123456789', '$2a$10$dummy.hash.for.testing'),
(1, 'Hafidh', 'hafidh@techcorp.com', '+628234567890', '$2a$10$dummy.hash.for.testing'),
(1, 'Jeri', 'jeri@techcorp.com', '+628345678901', '$2a$10$dummy.hash.for.testing');
```

## User Workflows

### Admin Workflow (Daily Operations)

**Morning Setup (9:00 AM):**
1. **Access Admin Panel**: `http://localhost:8080/admin/`
2. **Set Daily Menu**: Navigate to Daily Menu, select available items from master menu
3. **Create Order Sessions**: Go to Order Sessions, create sessions for each company requiring lunch delivery
4. **Monitor Dashboard**: Track session status and order counts

**Order Management (9:00-11:30 AM):**
1. **Monitor Sessions**: View real-time order status on dashboard
2. **Handle Changes**: Reopen/close sessions as needed for special circumstances
3. **View Orders**: Click "View Orders" to see individual employee orders

**Delivery Preparation (11:30 AM):**
1. **Close Sessions**: Close all order sessions to finalize orders
2. **Generate Delivery List**: View session orders for packaging and delivery manifest
3. **Update Status**: Change session status to DELIVERED after delivery

**Payment Tracking:**
1. **Mark Payments**: Update individual order payment status as received
2. **Monitor Collection**: Track payment completion per company
3. **Complete Sessions**: Mark sessions as COMPLETED when all payments received

### Employee Workflow

**Registration:**
1. **Sign Up**: Visit `http://localhost:8080/signup`
2. **Select Company**: Choose from registered companies
3. **Provide Details**: Enter name, email, WhatsApp number, and password
4. **Account Activation**: Immediate access after successful registration

**Daily Ordering:**
1. **Login**: Access `http://localhost:8080/login`
2. **View Dashboard**: See today's order session and recent order history
3. **Place Order**: Click "Place Order Now" if session is open
4. **Select Items**: Choose from today's available menu with real-time price calculation
5. **🤖 AI Nutritionist**: Click "AI Nutritionist" button for intelligent meal recommendations based on nutritional balance
6. **Submit/Update**: Confirm order (can modify until session closes)
7. **Track Status**: Monitor order and payment status on dashboard

## API Endpoints

### Public Routes
- `GET /` - Redirect to login
- `GET /login` - Login form
- `POST /login` - Process login
- `GET /signup` - Registration form
- `POST /signup` - Process registration

### Protected Customer Routes (Authentication Required)
- `GET /logout` - Logout user
- `GET /my-orders` - Customer dashboard
- `GET /order/:company/:date` - Order form for specific company/date
- `POST /order` - Submit/update order
- `POST /order/:company/:date/nutritionist-select` - Get AI nutritionist recommendations

### Admin Routes
- `GET /admin/` - Admin dashboard
- `GET /admin/menu` - Menu items management
- `POST /admin/menu` - Create menu item
- `PUT /admin/menu/:id` - Update menu item
- `DELETE /admin/menu/:id` - Delete menu item
- `GET /admin/companies` - Company management
- `POST /admin/companies` - Create company
- `PUT /admin/companies/:id` - Update company
- `DELETE /admin/companies/:id` - Delete company
- `GET /admin/companies/:id/employees` - Employee management for company
- `POST /admin/employees` - Create employee
- `PUT /admin/employees/:id` - Update employee
- `DELETE /admin/employees/:id` - Delete employee
- `GET /admin/daily-menu` - Daily menu management
- `POST /admin/daily-menu` - Set daily menu
- `GET /admin/sessions` - Order sessions management
- `POST /admin/sessions` - Create order session
- `POST /admin/sessions/:id/close` - Close order session
- `POST /admin/sessions/:id/reopen` - Reopen order session
- `GET /admin/sessions/:id/orders` - View orders in session
- `POST /admin/orders/:id/paid` - Mark order as paid
- `POST /admin/orders/:id/unpaid` - Mark order as unpaid

## Features Details

### 🤖 AI Nutritionist Feature

The AI Nutritionist is an intelligent meal recommendation system that leverages Large Language Models (LLMs) to provide personalized, nutritionally balanced meal suggestions from the daily menu.

#### How It Works

**Smart Recommendations:**
- Analyzes the entire daily menu using advanced AI
- Selects 2-4 items that provide optimal nutritional balance
- Prioritizes protein sources, vegetables, whole grains, and balanced portions
- Avoids excessive fried foods, sugar, and unbalanced combinations

**Nutritional Analysis:**
- **Protein Assessment**: Evaluates protein content (high/moderate/low)
- **Vegetable Content**: Assesses vegetable intake (high/moderate/low/none)
- **Carbohydrate Balance**: Analyzes carbohydrate levels (high/moderate/low)  
- **Overall Rating**: Provides overall nutritional rating (excellent/good/balanced/adequate)

**Intelligent Caching System:**
- Caches AI recommendations by date to reduce API costs
- Automatically invalidates cache when admin updates menu
- Tracks user adoption for analytics and notifications
- Supports admin-triggered cache resets via `nutritionist_reset` flag

**User Experience:**
- One-click "AI Nutritionist" button on order form
- Real-time AI processing with loading indicators
- Clear reasoning explanation for recommendations
- Visual nutritional summary with color-coded ratings
- Seamless integration with existing order flow

#### Technical Implementation

**Architecture:**
- `internal/services/nutritionist.go`: Core AI service with LLM integration
- Smart caching with PostgreSQL JSONB storage
- Fallback parsing for robust AI response handling
- User tracking for notification systems

**Database Schema:**
- `nutritionist_selections`: Stores AI recommendations with reasoning
- `nutritionist_user_selections`: Tracks user adoption patterns
- `daily_menus.nutritionist_reset`: Admin control for cache invalidation

**AI Integration:**
- LLM client with structured JSON response parsing
- Robust error handling and fallback mechanisms
- Configurable prompts for nutritional optimization
- Index-based menu item selection for accuracy

#### Admin Features

**Cache Management:**
- Reset nutritionist recommendations when menu changes
- Automatic user notification system for menu updates
- Analytics on user adoption of AI recommendations

**Menu Integration:**
- Seamless integration with existing daily menu system
- Automatic cache invalidation on menu modifications
- Support for menu item additions/removals

### Menu Management
- **Master Menu**: Comprehensive list of all possible menu items with fixed prices
- **Daily Selection**: Admin selects subset of master menu available for each day
- **Price Management**: Prices stored in cents for accuracy, displayed in Rupiah
- **Active Status**: Soft delete functionality for menu items

### Company & Employee Management
- **Company Profiles**: Name, address, contact information
- **Employee Self-Service**: Registration with company association
- **WhatsApp Integration**: Contact numbers for communication
- **Account Management**: Edit/delete capabilities with proper referential integrity

### Order Session Management
- **Date-Based Sessions**: One session per company per day
- **Status Tracking**: Complete lifecycle from OPEN to COMPLETED
- **Flexible Control**: Admins can close/reopen sessions for emergency situations
- **Duplicate Prevention**: System prevents duplicate sessions for same company/date

### Order Processing
- **Real-Time Calculation**: Dynamic price calculation as items are selected
- **Order Updates**: Employees can modify orders until session closes
- **Visual Feedback**: Selected items highlighted with summary display
- **Validation**: Server-side validation for all order data

### Payment Tracking
- **Individual Tracking**: Each employee order tracked separately
- **Status Toggle**: Admin can mark paid/unpaid (handles mistakes)
- **Session Summary**: Total revenue and payment completion rates
- **Company Billing**: Aggregate view for company payment processing

## Technical Features

### Security
- **Password Hashing**: bcrypt with appropriate cost factor
- **Session Management**: HTTP-only cookies for authentication
- **Input Validation**: Both client and server-side validation
- **SQL Injection Prevention**: Parameterized queries throughout

### Database Design
- **Referential Integrity**: Proper foreign key relationships
- **Array Storage**: PostgreSQL arrays for menu item selections
- **Soft Deletes**: Active flags instead of hard deletion
- **Optimized Queries**: Proper indexing and query optimization

### User Experience
- **Responsive Design**: Works on desktop and mobile devices
- **Intuitive Navigation**: Clear menu structure and breadcrumbs
- **Real-Time Updates**: Dynamic price calculation and status updates
- **Error Handling**: Comprehensive error messages and graceful degradation

### Performance
- **Connection Pooling**: Efficient database connection management
- **Template Caching**: Compiled templates for faster rendering
- **Minimal Dependencies**: Lightweight framework choices
- **Static Asset Management**: Proper CSS/JS organization

## Deployment Considerations

### Production Setup
- Change default passwords and secrets
- Configure proper TLS certificates
- Set up database backups
- Implement logging and monitoring
- Configure reverse proxy (nginx/Apache)
- Set up proper firewall rules

### Scaling Options
- Database read replicas for high query loads
- Load balancer for multiple application instances
- Redis for session storage in multi-instance setup
- CDN for static assets

### Monitoring
- Database query performance monitoring
- Application response time tracking
- Error rate monitoring
- User activity analytics

## Contributing

This application was built with a focus on simplicity, reliability, and maintainability. The codebase follows Go best practices and maintains clear separation of concerns.

### Development Guidelines
- Follow Go formatting standards (`go fmt`)
- Write tests for new functionality
- Update documentation for API changes
- Use meaningful commit messages
- Maintain backwards compatibility when possible

## License

MIT License

Copyright (c) 2025 Lunch Delivery System

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Support

For technical issues or feature requests, please create an issue in this repo.
