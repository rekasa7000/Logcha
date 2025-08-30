import { z } from 'zod'

// Auth schemas
export const loginSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
})

export const registerSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
  first_name: z.string().min(1, 'First name is required'),
  last_name: z.string().min(1, 'Last name is required'),
  phone: z.string().optional(),
  role: z.enum(['company_admin', 'trainee']),
})

// Company schemas
export const companySchema = z.object({
  name: z.string().min(1, 'Company name is required'),
  address: z.string().optional(),
  contact_person: z.string().optional(),
  contact_email: z.string().email('Invalid email address').optional().or(z.literal('')),
  contact_phone: z.string().optional(),
})

// Trainee schemas
export const traineeSchema = z.object({
  user_id: z.string().min(1, 'User is required'),
  company_id: z.string().min(1, 'Company is required'),
  employee_id: z.string().optional(),
  trainee_type: z.enum(['paid_intern', 'unpaid_intern', 'ojt']),
  hourly_rate: z.number().min(0).optional(),
  max_weekly_hours: z.number().min(1, 'Maximum weekly hours must be at least 1'),
  total_required_hours: z.number().min(1).optional(),
  start_date: z.string().min(1, 'Start date is required'),
  end_date: z.string().optional(),
  school_name: z.string().optional(),
  course: z.string().optional(),
  year_level: z.string().optional(),
}).superRefine((data, ctx) => {
  // Paid interns must have hourly_rate
  if (data.trainee_type === 'paid_intern' && !data.hourly_rate) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'Hourly rate is required for paid interns',
      path: ['hourly_rate'],
    })
  }
  
  // OJT students must have total_required_hours
  if (data.trainee_type === 'ojt' && !data.total_required_hours) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'Total required hours is required for OJT students',
      path: ['total_required_hours'],
    })
  }

  // End date must be after start date if provided
  if (data.end_date && data.start_date && new Date(data.end_date) <= new Date(data.start_date)) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'End date must be after start date',
      path: ['end_date'],
    })
  }
})

// Time record schemas
export const timeRecordSchema = z.object({
  trainee_id: z.string().min(1, 'Trainee is required'),
  date: z.string().min(1, 'Date is required'),
  am_time_in: z.string().optional(),
  am_time_out: z.string().optional(),
  pm_time_in: z.string().optional(),
  pm_time_out: z.string().optional(),
  notes: z.string().optional(),
  status: z.enum(['present', 'half_day_am', 'half_day_pm', 'absent']),
}).superRefine((data, ctx) => {
  // Validate AM session
  if (data.am_time_in && data.am_time_out) {
    const amIn = new Date(`2000-01-01T${data.am_time_in}:00`)
    const amOut = new Date(`2000-01-01T${data.am_time_out}:00`)
    
    if (amOut <= amIn) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'AM time out must be after AM time in',
        path: ['am_time_out'],
      })
    }
  }

  // Validate PM session
  if (data.pm_time_in && data.pm_time_out) {
    const pmIn = new Date(`2000-01-01T${data.pm_time_in}:00`)
    const pmOut = new Date(`2000-01-01T${data.pm_time_out}:00`)
    
    if (pmOut <= pmIn) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'PM time out must be after PM time in',
        path: ['pm_time_out'],
      })
    }
  }

  // Validate PM is after AM
  if (data.am_time_out && data.pm_time_in) {
    const amOut = new Date(`2000-01-01T${data.am_time_out}:00`)
    const pmIn = new Date(`2000-01-01T${data.pm_time_in}:00`)
    
    if (pmIn <= amOut) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'PM time in should be after AM time out',
        path: ['pm_time_in'],
      })
    }
  }

  // At least one session must be present for 'present' status
  if (data.status === 'present' && 
      (!data.am_time_in || !data.am_time_out) && 
      (!data.pm_time_in || !data.pm_time_out)) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'At least one complete session (AM or PM) is required for present status',
      path: ['status'],
    })
  }

  // Validate future dates
  const recordDate = new Date(data.date)
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  
  if (recordDate > today) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'Cannot enter future dates',
      path: ['date'],
    })
  }
})