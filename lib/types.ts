export type UserRole = 'company_admin' | 'trainee'

export type TraineeType = 'paid_intern' | 'unpaid_intern' | 'ojt'

export type TraineeStatus = 'active' | 'completed' | 'terminated'

export type TimeRecordStatus = 'present' | 'half_day_am' | 'half_day_pm' | 'absent'

export interface User {
  id: string
  email: string
  role: UserRole
  first_name: string
  last_name: string
  phone?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Company {
  id: string
  name: string
  address?: string
  contact_person?: string
  contact_email?: string
  contact_phone?: string
  created_at: string
  updated_at: string
}

export interface Trainee {
  id: string
  user_id: string
  company_id: string
  employee_id?: string
  trainee_type: TraineeType
  hourly_rate?: number
  max_weekly_hours: number
  total_required_hours?: number
  start_date: string
  end_date?: string
  status: TraineeStatus
  school_name?: string
  course?: string
  year_level?: string
  created_at: string
  updated_at: string
  // Relations
  user?: User
  company?: Company
}

export interface TimeRecord {
  id: string
  trainee_id: string
  date: string
  am_time_in?: string
  am_time_out?: string
  am_hours?: number
  pm_time_in?: string
  pm_time_out?: string
  pm_hours?: number
  total_hours?: number
  notes?: string
  status: TimeRecordStatus
  created_at: string
  updated_at: string
  // Relations
  trainee?: Trainee
}

export interface WeeklySummary {
  id: string
  trainee_id: string
  week_start_date: string
  week_end_date: string
  total_hours_worked: number
  billable_hours: number
  gross_pay: number
  days_present: number
  created_at: string
  updated_at: string
  // Relations
  trainee?: Trainee
}

export interface MonthlyReport {
  id: string
  trainee_id: string
  month: number
  year: number
  total_hours_worked: number
  total_billable_hours: number
  total_gross_pay: number
  days_present: number
  days_absent: number
  generated_at: string
  // Relations
  trainee?: Trainee
}

// DTOs for forms
export interface CreateTraineeDto {
  user_id: string
  company_id: string
  employee_id?: string
  trainee_type: TraineeType
  hourly_rate?: number
  max_weekly_hours: number
  total_required_hours?: number
  start_date: string
  end_date?: string
  school_name?: string
  course?: string
  year_level?: string
}

export interface CreateTimeRecordDto {
  trainee_id: string
  date: string
  am_time_in?: string
  am_time_out?: string
  pm_time_in?: string
  pm_time_out?: string
  notes?: string
  status: TimeRecordStatus
}

export interface OJTProgress {
  trainee_id: string
  first_name: string
  last_name: string
  total_required_hours: number
  hours_rendered: number
  remaining_hours: number
  completion_percentage: number
}