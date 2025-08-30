# Database Schema Documentation

## Internship & OJT Time Tracking System

This document outlines the complete database schema for the web application that tracks and monitors staff, interns, and OJT (On-the-Job Training) time records.

---

## System Overview

The system handles three types of trainees:

- **Paid Interns**: Company-defined maximum weekly hours, paid hourly rate
- **Unpaid Interns**: Company-defined maximum weekly hours, no payment
- **OJT Students**: Total required hours (e.g., 500h), no payment, academic requirement

### Key Features

- AM/PM session time tracking with lunch break handling
- Company-defined maximum weekly hours for each intern
- Progress tracking for OJT total hour requirements
- Monthly DTR (Daily Time Record) report generation
- Multi-company support

---

## Core Database Tables

### 1. Users Table

```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('company_admin', 'trainee') NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_email (email),
    INDEX idx_role (role)
);
```

### 2. Companies Table

```sql
CREATE TABLE companies (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    contact_person VARCHAR(100),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_name (name)
);
```

### 3. Trainees Table

```sql
CREATE TABLE trainees (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    company_id BIGINT NOT NULL,
    employee_id VARCHAR(50), -- Company-specific ID
    trainee_type ENUM('paid_intern', 'unpaid_intern', 'ojt') NOT NULL,

    -- Intern-specific fields
    hourly_rate DECIMAL(10,2) NULL, -- NULL for unpaid interns and OJT
    max_weekly_hours INT NOT NULL, -- Company-defined maximum hours per week for interns

    -- OJT-specific fields
    total_required_hours INT NULL, -- Required total hours for OJT (e.g., 500)

    -- Common fields
    start_date DATE NOT NULL,
    end_date DATE,
    status ENUM('active', 'completed', 'terminated') DEFAULT 'active',
    school_name VARCHAR(255),
    course VARCHAR(255),
    year_level VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,

    INDEX idx_company_status (company_id, status),
    INDEX idx_trainee_type (trainee_type),
    INDEX idx_user_id (user_id),
    UNIQUE KEY unique_employee_company (employee_id, company_id)
);
```

### 4. Time Records Table

```sql
CREATE TABLE time_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    trainee_id BIGINT NOT NULL,
    date DATE NOT NULL,

    -- AM Session (typically 8AM-12PM)
    am_time_in TIME NULL,
    am_time_out TIME NULL,
    am_hours DECIMAL(4,2) GENERATED ALWAYS AS (
        CASE
            WHEN am_time_in IS NOT NULL AND am_time_out IS NOT NULL
            THEN TIME_TO_SEC(TIMEDIFF(am_time_out, am_time_in)) / 3600.0
            ELSE 0
        END
    ) STORED,

    -- PM Session (typically 1PM-5PM)
    pm_time_in TIME NULL,
    pm_time_out TIME NULL,
    pm_hours DECIMAL(4,2) GENERATED ALWAYS AS (
        CASE
            WHEN pm_time_in IS NOT NULL AND pm_time_out IS NOT NULL
            THEN TIME_TO_SEC(TIMEDIFF(pm_time_out, pm_time_in)) / 3600.0
            ELSE 0
        END
    ) STORED,

    -- Daily total
    total_hours DECIMAL(4,2) GENERATED ALWAYS AS (
        COALESCE(am_hours, 0) + COALESCE(pm_hours, 0)
    ) STORED,

    notes TEXT,
    status ENUM('present', 'half_day_am', 'half_day_pm', 'absent') DEFAULT 'present',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (trainee_id) REFERENCES trainees(id) ON DELETE CASCADE,

    UNIQUE KEY unique_trainee_date (trainee_id, date),
    INDEX idx_trainee_date (trainee_id, date),
    INDEX idx_date (date),
    INDEX idx_status (status)
);
```

### 5. Weekly Summaries Table

```sql
CREATE TABLE weekly_summaries (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    trainee_id BIGINT NOT NULL,
    week_start_date DATE NOT NULL, -- Monday of the week
    week_end_date DATE NOT NULL,   -- Sunday of the week
    total_hours_worked DECIMAL(6,2) DEFAULT 0,
    billable_hours DECIMAL(6,2) DEFAULT 0, -- Capped at max_weekly_hours for paid interns
    gross_pay DECIMAL(10,2) DEFAULT 0, -- billable_hours * hourly_rate
    days_present INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (trainee_id) REFERENCES trainees(id) ON DELETE CASCADE,

    UNIQUE KEY unique_trainee_week (trainee_id, week_start_date),
    INDEX idx_week_dates (week_start_date, week_end_date)
);
```

### 6. Monthly Reports Table

```sql
CREATE TABLE monthly_reports (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    trainee_id BIGINT NOT NULL,
    month INT NOT NULL, -- 1-12
    year INT NOT NULL,
    total_hours_worked DECIMAL(8,2) DEFAULT 0,
    total_billable_hours DECIMAL(8,2) DEFAULT 0,
    total_gross_pay DECIMAL(12,2) DEFAULT 0,
    days_present INT DEFAULT 0,
    days_absent INT DEFAULT 0,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (trainee_id) REFERENCES trainees(id) ON DELETE CASCADE,

    UNIQUE KEY unique_trainee_month_year (trainee_id, month, year),
    INDEX idx_month_year (month, year)
);
```

---

## Business Logic & Calculations

### Time Entry Validation Rules

1. `am_time_out` must be greater than `am_time_in` (if both present)
2. `pm_time_out` must be greater than `pm_time_in` (if both present)
3. `pm_time_in` should be after `am_time_out` (if both present)
4. At least one session (AM or PM) must be present for status = 'present'
5. Cannot enter future dates
6. Cannot modify time records older than 7 days (business rule)

### Weekly Summary Calculation

```sql
-- Example: Calculate weekly summary for a trainee
INSERT INTO weekly_summaries (trainee_id, week_start_date, week_end_date, total_hours_worked, billable_hours, gross_pay, days_present)
SELECT
    tr.trainee_id,
    DATE_SUB(tr.date, INTERVAL WEEKDAY(tr.date) DAY) as week_start_date,
    DATE_ADD(DATE_SUB(tr.date, INTERVAL WEEKDAY(tr.date) DAY), INTERVAL 6 DAY) as week_end_date,
    SUM(tr.total_hours) as total_hours_worked,
    CASE
        WHEN t.trainee_type IN ('paid_intern', 'unpaid_intern') THEN LEAST(SUM(tr.total_hours), t.max_weekly_hours)
        ELSE SUM(tr.total_hours) -- OJT has no weekly limit
    END as billable_hours,
    CASE
        WHEN t.trainee_type = 'paid_intern' THEN LEAST(SUM(tr.total_hours), t.max_weekly_hours) * t.hourly_rate
        ELSE 0
    END as gross_pay,
    COUNT(CASE WHEN tr.status IN ('present', 'half_day_am', 'half_day_pm') THEN 1 END) as days_present
FROM time_records tr
JOIN trainees t ON tr.trainee_id = t.id
WHERE tr.trainee_id = ?
  AND tr.date BETWEEN ? AND ?
GROUP BY tr.trainee_id, week_start_date, week_end_date;
```

### OJT Progress Tracking

```sql
-- Check OJT student progress
SELECT
    u.first_name,
    u.last_name,
    t.total_required_hours,
    COALESCE(SUM(tr.total_hours), 0) as hours_rendered,
    (t.total_required_hours - COALESCE(SUM(tr.total_hours), 0)) as remaining_hours,
    ROUND((COALESCE(SUM(tr.total_hours), 0) / t.total_required_hours) * 100, 2) as completion_percentage
FROM trainees t
JOIN users u ON t.user_id = u.id
LEFT JOIN time_records tr ON t.id = tr.trainee_id
WHERE t.trainee_type = 'ojt' AND t.status = 'active'
GROUP BY t.id;
```

---

## Sample Data Examples

### Example 1: Paid Intern (Juan) - 16 hours max per week

```sql
-- Juan's profile (company sets 16 hours max per week)
INSERT INTO trainees (user_id, company_id, trainee_type, hourly_rate, max_weekly_hours)
VALUES (1, 1, 'paid_intern', 66.66, 16);

-- Juan's time record: 8AM-12PM, 1PM-5PM (8 hours total)
INSERT INTO time_records (trainee_id, date, am_time_in, am_time_out, pm_time_in, pm_time_out)
VALUES (1, '2025-01-15', '08:00:00', '12:00:00', '13:00:00', '17:00:00');
```

### Example 2: Paid Intern (Juancho) - Works 20 hours but only paid for 16

```sql
-- Juancho's profile (company sets 16 hours max per week)
INSERT INTO trainees (user_id, company_id, trainee_type, hourly_rate, max_weekly_hours)
VALUES (3, 1, 'paid_intern', 75.00, 16);

-- Juancho works 4 hours every day (Monday-Friday = 20 hours total)
-- But only gets paid for 16 hours (his max_weekly_hours)
INSERT INTO time_records (trainee_id, date, am_time_in, am_time_out, pm_time_in, pm_time_out) VALUES
(3, '2025-01-13', '08:00:00', '12:00:00', NULL, NULL), -- Monday: 4 hours
(3, '2025-01-14', '08:00:00', '12:00:00', NULL, NULL), -- Tuesday: 4 hours
(3, '2025-01-15', '08:00:00', '12:00:00', NULL, NULL), -- Wednesday: 4 hours
(3, '2025-01-16', '08:00:00', '12:00:00', NULL, NULL), -- Thursday: 4 hours
(3, '2025-01-17', '08:00:00', '12:00:00', NULL, NULL); -- Friday: 4 hours
-- Total: 20 hours worked, but billable_hours = 16 (capped at max_weekly_hours)
```

### Example 3: Paid Intern (Pedro) - Works minimal hours

```sql
-- Pedro's profile (company sets 10 hours max per week)
INSERT INTO trainees (user_id, company_id, trainee_type, hourly_rate, max_weekly_hours)
VALUES (4, 1, 'paid_intern', 80.00, 10);

-- Pedro works only 1 day this week (3 hours total)
-- Gets paid for 3 hours (under his max_weekly_hours of 10)
INSERT INTO time_records (trainee_id, date, am_time_in, am_time_out, pm_time_in, pm_time_out)
VALUES (4, '2025-01-15', '09:00:00', '12:00:00', NULL, NULL); -- 3 hours only
```

### Example 4: OJT Student (Maria) - No maximum weekly hours

```sql
-- Maria's profile (500 required hours total, no weekly limit)
INSERT INTO trainees (user_id, company_id, trainee_type, total_required_hours, max_weekly_hours)
VALUES (5, 1, 'ojt', 500, NULL); -- max_weekly_hours is NULL for OJT

-- Maria's time record: 8AM-12PM, 1PM-6PM (9 hours total, all counted towards 500 hours)
INSERT INTO time_records (trainee_id, date, am_time_in, am_time_out, pm_time_in, pm_time_out)
VALUES (5, '2025-01-15', '08:00:00', '12:00:00', '13:00:00', '18:00:00');
```

---

## Report Queries

### Daily Time Record (DTR) Query

```sql
SELECT
    DATE_FORMAT(tr.date, '%Y-%m-%d') as date,
    tr.am_time_in,
    tr.am_time_out,
    tr.pm_time_in,
    tr.pm_time_out,
    tr.total_hours,
    tr.status,
    tr.notes
FROM time_records tr
WHERE tr.trainee_id = ?
  AND tr.date BETWEEN ? AND ?
ORDER BY tr.date;
```

### Monthly Summary Report

```sql
SELECT
    u.first_name,
    u.last_name,
    t.trainee_type,
    mr.month,
    mr.year,
    mr.total_hours_worked,
    mr.total_billable_hours,
    mr.total_gross_pay,
    mr.days_present,
    mr.days_absent
FROM monthly_reports mr
JOIN trainees t ON mr.trainee_id = t.id
JOIN users u ON t.user_id = u.id
WHERE t.company_id = ?
  AND mr.month = ?
  AND mr.year = ?
ORDER BY u.last_name, u.first_name;
```

---

## Database Indexes & Performance

### Recommended Indexes

- `time_records(trainee_id, date)` - For individual DTR queries
- `time_records(date)` - For company-wide daily reports
- `trainees(company_id, status)` - For active trainee lists
- `weekly_summaries(trainee_id, week_start_date)` - For payroll processing
- `monthly_reports(month, year)` - For company reports

### Data Retention Policy

- Keep `time_records` for 5 years (legal requirement)
- Archive `weekly_summaries` and `monthly_reports` after 3 years
- Soft delete `trainees` (set status to 'terminated' instead of hard delete)

---

## Future Enhancements

### Phase 2 Features (Post-MVP)

- Email notifications table
- Document management (contracts, NDAs)
- Requirement checklists for OJT
- Leave/absence request system
- Overtime tracking
- Holiday calendar management

### Additional Tables for Phase 2

```sql
-- Email communications
CREATE TABLE email_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    sender_id BIGINT,
    recipient_id BIGINT,
    subject VARCHAR(255),
    message TEXT,
    sent_at TIMESTAMP
);

-- Document storage
CREATE TABLE documents (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    trainee_id BIGINT,
    document_type ENUM('contract', 'nda', 'requirements', 'certificate'),
    file_name VARCHAR(255),
    file_path VARCHAR(500),
    uploaded_at TIMESTAMP
);
```

---

## Migration Notes

### Initial Data Setup

1. Create company admin accounts first
2. Set up company profiles
3. Import existing trainee data (if any)
4. Configure default hourly rates and weekly hour limits

### Data Migration from Existing System

If migrating from spreadsheets or another system:

1. Export existing time records to CSV format
2. Use bulk insert scripts for historical data
3. Generate weekly/monthly summaries for historical periods
4. Validate data integrity with sample calculations
