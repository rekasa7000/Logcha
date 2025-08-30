-- Logcha Database Schema Migration for Supabase PostgreSQL
-- Based on DATABASE_SCHEME.md

-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create custom types (using text with constraints instead of ENUM for better Supabase compatibility)
CREATE TYPE user_role AS ENUM ('company_admin', 'trainee');
CREATE TYPE trainee_type AS ENUM ('paid_intern', 'unpaid_intern', 'ojt');
CREATE TYPE trainee_status AS ENUM ('active', 'completed', 'terminated');
CREATE TYPE time_record_status AS ENUM ('present', 'half_day_am', 'half_day_pm', 'absent');

-- 1. Users Table (extends Supabase auth.users)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT auth.uid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    role user_role NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_users_auth_user FOREIGN KEY (id) REFERENCES auth.users(id) ON DELETE CASCADE
);

-- Create indexes for users table
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- 2. Companies Table
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    address TEXT,
    contact_person VARCHAR(100),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for companies table
CREATE INDEX idx_companies_name ON companies(name);

-- 3. Trainees Table
CREATE TABLE trainees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    company_id UUID NOT NULL,
    employee_id VARCHAR(50),
    trainee_type trainee_type NOT NULL,
    
    -- Intern-specific fields
    hourly_rate DECIMAL(10,2) NULL,
    max_weekly_hours INTEGER NOT NULL,
    
    -- OJT-specific fields
    total_required_hours INTEGER NULL,
    
    -- Common fields
    start_date DATE NOT NULL,
    end_date DATE,
    status trainee_status DEFAULT 'active',
    school_name VARCHAR(255),
    course VARCHAR(255),
    year_level VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_trainees_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_trainees_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT unique_employee_company UNIQUE (employee_id, company_id)
);

-- Create indexes for trainees table
CREATE INDEX idx_trainees_company_status ON trainees(company_id, status);
CREATE INDEX idx_trainees_type ON trainees(trainee_type);
CREATE INDEX idx_trainees_user_id ON trainees(user_id);

-- 4. Time Records Table
CREATE TABLE time_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trainee_id UUID NOT NULL,
    date DATE NOT NULL,
    
    -- AM Session
    am_time_in TIME,
    am_time_out TIME,
    am_hours DECIMAL(4,2) GENERATED ALWAYS AS (
        CASE
            WHEN am_time_in IS NOT NULL AND am_time_out IS NOT NULL
            THEN EXTRACT(EPOCH FROM (am_time_out - am_time_in)) / 3600.0
            ELSE 0
        END
    ) STORED,
    
    -- PM Session
    pm_time_in TIME,
    pm_time_out TIME,
    pm_hours DECIMAL(4,2) GENERATED ALWAYS AS (
        CASE
            WHEN pm_time_in IS NOT NULL AND pm_time_out IS NOT NULL
            THEN EXTRACT(EPOCH FROM (pm_time_out - pm_time_in)) / 3600.0
            ELSE 0
        END
    ) STORED,
    
    -- Daily total (calculate directly from time fields to avoid referencing other generated columns)
    total_hours DECIMAL(4,2) GENERATED ALWAYS AS (
        CASE
            WHEN am_time_in IS NOT NULL AND am_time_out IS NOT NULL
            THEN EXTRACT(EPOCH FROM (am_time_out - am_time_in)) / 3600.0
            ELSE 0
        END +
        CASE
            WHEN pm_time_in IS NOT NULL AND pm_time_out IS NOT NULL
            THEN EXTRACT(EPOCH FROM (pm_time_out - pm_time_in)) / 3600.0
            ELSE 0
        END
    ) STORED,
    
    notes TEXT,
    status time_record_status DEFAULT 'present',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_time_records_trainee FOREIGN KEY (trainee_id) REFERENCES trainees(id) ON DELETE CASCADE,
    CONSTRAINT unique_trainee_date UNIQUE (trainee_id, date),
    CONSTRAINT chk_am_times CHECK (am_time_out IS NULL OR am_time_in IS NULL OR am_time_out > am_time_in),
    CONSTRAINT chk_pm_times CHECK (pm_time_out IS NULL OR pm_time_in IS NULL OR pm_time_out > pm_time_in)
);

-- Create indexes for time_records table
CREATE INDEX idx_time_records_trainee_date ON time_records(trainee_id, date);
CREATE INDEX idx_time_records_date ON time_records(date);
CREATE INDEX idx_time_records_status ON time_records(status);

-- 5. Weekly Summaries Table
CREATE TABLE weekly_summaries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trainee_id UUID NOT NULL,
    week_start_date DATE NOT NULL,
    week_end_date DATE NOT NULL,
    total_hours_worked DECIMAL(6,2) DEFAULT 0,
    billable_hours DECIMAL(6,2) DEFAULT 0,
    gross_pay DECIMAL(10,2) DEFAULT 0,
    days_present INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_weekly_summaries_trainee FOREIGN KEY (trainee_id) REFERENCES trainees(id) ON DELETE CASCADE,
    CONSTRAINT unique_trainee_week UNIQUE (trainee_id, week_start_date)
);

-- Create indexes for weekly_summaries table
CREATE INDEX idx_weekly_summaries_week_dates ON weekly_summaries(week_start_date, week_end_date);

-- 6. Monthly Reports Table
CREATE TABLE monthly_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trainee_id UUID NOT NULL,
    month INTEGER NOT NULL CHECK (month >= 1 AND month <= 12),
    year INTEGER NOT NULL,
    total_hours_worked DECIMAL(8,2) DEFAULT 0,
    total_billable_hours DECIMAL(8,2) DEFAULT 0,
    total_gross_pay DECIMAL(12,2) DEFAULT 0,
    days_present INTEGER DEFAULT 0,
    days_absent INTEGER DEFAULT 0,
    generated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_monthly_reports_trainee FOREIGN KEY (trainee_id) REFERENCES trainees(id) ON DELETE CASCADE,
    CONSTRAINT unique_trainee_month_year UNIQUE (trainee_id, month, year)
);

-- Create indexes for monthly_reports table
CREATE INDEX idx_monthly_reports_month_year ON monthly_reports(month, year);

-- Create updated_at triggers
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON companies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_trainees_updated_at BEFORE UPDATE ON trainees FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_time_records_updated_at BEFORE UPDATE ON time_records FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_weekly_summaries_updated_at BEFORE UPDATE ON weekly_summaries FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Row Level Security (RLS) Policies
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE companies ENABLE ROW LEVEL SECURITY;
ALTER TABLE trainees ENABLE ROW LEVEL SECURITY;
ALTER TABLE time_records ENABLE ROW LEVEL SECURITY;
ALTER TABLE weekly_summaries ENABLE ROW LEVEL SECURITY;
ALTER TABLE monthly_reports ENABLE ROW LEVEL SECURITY;

-- Users can read/update their own profile
CREATE POLICY "Users can view their own profile" ON users
    FOR SELECT USING (auth.uid() = id);

CREATE POLICY "Users can update their own profile" ON users
    FOR UPDATE USING (auth.uid() = id);

-- Company admins can manage companies they belong to
CREATE POLICY "Company admins can view companies" ON companies
    FOR SELECT USING (
        EXISTS (
            SELECT 1 FROM users u 
            JOIN trainees t ON u.id = t.user_id 
            WHERE u.id = auth.uid() 
            AND u.role = 'company_admin' 
            AND t.company_id = companies.id
        )
    );

-- Company admins can manage trainees in their companies
CREATE POLICY "Company admins can view trainees in their companies" ON trainees
    FOR SELECT USING (
        EXISTS (
            SELECT 1 FROM users u 
            JOIN trainees t ON u.id = t.user_id 
            WHERE u.id = auth.uid() 
            AND u.role = 'company_admin' 
            AND t.company_id = trainees.company_id
        )
        OR user_id = auth.uid()
    );

-- Trainees can view their own records
CREATE POLICY "Trainees can view their own time records" ON time_records
    FOR SELECT USING (
        EXISTS (
            SELECT 1 FROM trainees t 
            WHERE t.id = time_records.trainee_id 
            AND t.user_id = auth.uid()
        )
        OR EXISTS (
            SELECT 1 FROM users u 
            JOIN trainees t ON u.id = t.user_id 
            WHERE u.id = auth.uid() 
            AND u.role = 'company_admin' 
            AND t.company_id IN (
                SELECT company_id FROM trainees WHERE id = time_records.trainee_id
            )
        )
    );

-- Trainees can insert their own time records
CREATE POLICY "Trainees can insert their own time records" ON time_records
    FOR INSERT WITH CHECK (
        EXISTS (
            SELECT 1 FROM trainees t 
            WHERE t.id = trainee_id 
            AND t.user_id = auth.uid()
        )
    );

-- Trainees can update their own time records (within 7 days)
CREATE POLICY "Trainees can update their own recent time records" ON time_records
    FOR UPDATE USING (
        EXISTS (
            SELECT 1 FROM trainees t 
            WHERE t.id = time_records.trainee_id 
            AND t.user_id = auth.uid()
        )
        AND date >= CURRENT_DATE - INTERVAL '7 days'
    );

-- Similar policies for weekly_summaries and monthly_reports
CREATE POLICY "View own weekly summaries" ON weekly_summaries
    FOR SELECT USING (
        EXISTS (
            SELECT 1 FROM trainees t 
            WHERE t.id = weekly_summaries.trainee_id 
            AND t.user_id = auth.uid()
        )
        OR EXISTS (
            SELECT 1 FROM users u 
            JOIN trainees t ON u.id = t.user_id 
            WHERE u.id = auth.uid() 
            AND u.role = 'company_admin' 
            AND t.company_id IN (
                SELECT company_id FROM trainees WHERE id = weekly_summaries.trainee_id
            )
        )
    );

CREATE POLICY "View own monthly reports" ON monthly_reports
    FOR SELECT USING (
        EXISTS (
            SELECT 1 FROM trainees t 
            WHERE t.id = monthly_reports.trainee_id 
            AND t.user_id = auth.uid()
        )
        OR EXISTS (
            SELECT 1 FROM users u 
            JOIN trainees t ON u.id = t.user_id 
            WHERE u.id = auth.uid() 
            AND u.role = 'company_admin' 
            AND t.company_id IN (
                SELECT company_id FROM trainees WHERE id = monthly_reports.trainee_id
            )
        )
    );

-- Functions for calculating weekly summaries
CREATE OR REPLACE FUNCTION calculate_weekly_summary(p_trainee_id UUID, p_week_start DATE)
RETURNS TABLE (
    total_hours_worked DECIMAL(6,2),
    billable_hours DECIMAL(6,2),
    gross_pay DECIMAL(10,2),
    days_present INTEGER
) AS $$
DECLARE
    trainee_record RECORD;
BEGIN
    -- Get trainee info
    SELECT * INTO trainee_record 
    FROM trainees 
    WHERE id = p_trainee_id;
    
    -- Calculate totals
    SELECT 
        COALESCE(SUM(tr.total_hours), 0),
        CASE
            WHEN trainee_record.trainee_type IN ('paid_intern', 'unpaid_intern') THEN 
                LEAST(COALESCE(SUM(tr.total_hours), 0), trainee_record.max_weekly_hours)
            ELSE 
                COALESCE(SUM(tr.total_hours), 0) -- OJT has no weekly limit
        END,
        CASE
            WHEN trainee_record.trainee_type = 'paid_intern' THEN 
                LEAST(COALESCE(SUM(tr.total_hours), 0), trainee_record.max_weekly_hours) * trainee_record.hourly_rate
            ELSE 
                0
        END,
        COUNT(CASE WHEN tr.status IN ('present', 'half_day_am', 'half_day_pm') THEN 1 END)::INTEGER
    INTO total_hours_worked, billable_hours, gross_pay, days_present
    FROM time_records tr
    WHERE tr.trainee_id = p_trainee_id
      AND tr.date BETWEEN p_week_start AND (p_week_start + INTERVAL '6 days')::DATE;
      
    RETURN NEXT;
END;
$$ LANGUAGE plpgsql;

-- Function to get OJT progress
CREATE OR REPLACE FUNCTION get_ojt_progress(p_trainee_id UUID DEFAULT NULL)
RETURNS TABLE (
    trainee_id UUID,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    total_required_hours INTEGER,
    hours_rendered DECIMAL,
    remaining_hours DECIMAL,
    completion_percentage DECIMAL
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        t.id,
        u.first_name,
        u.last_name,
        t.total_required_hours,
        COALESCE(SUM(tr.total_hours), 0),
        (t.total_required_hours - COALESCE(SUM(tr.total_hours), 0)),
        ROUND((COALESCE(SUM(tr.total_hours), 0) / t.total_required_hours) * 100, 2)
    FROM trainees t
    JOIN users u ON t.user_id = u.id
    LEFT JOIN time_records tr ON t.id = tr.trainee_id
    WHERE t.trainee_type = 'ojt' 
      AND t.status = 'active'
      AND (p_trainee_id IS NULL OR t.id = p_trainee_id)
    GROUP BY t.id, u.first_name, u.last_name, t.total_required_hours;
END;
$$ LANGUAGE plpgsql;