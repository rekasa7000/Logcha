#!/usr/bin/env bun

import { createClient } from '@supabase/supabase-js'
import { config } from 'dotenv'
import { resolve } from 'path'

// Load environment variables
config({ path: resolve(process.cwd(), '.env.local') })

const supabaseUrl = process.env.NEXT_PUBLIC_SUPABASE_URL
const supabaseServiceKey = process.env.SUPABASE_SERVICE_ROLE_KEY

if (!supabaseUrl || !supabaseServiceKey) {
  console.error('❌ Missing required environment variables:')
  console.error('   - NEXT_PUBLIC_SUPABASE_URL')
  console.error('   - SUPABASE_SERVICE_ROLE_KEY')
  console.error('\nPlease check your .env.local file')
  process.exit(1)
}

// Create Supabase client with service role key for admin operations
const supabase = createClient(supabaseUrl, supabaseServiceKey, {
  auth: {
    autoRefreshToken: false,
    persistSession: false
  }
})

interface AdminUserConfig {
  email: string
  password: string
  firstName: string
  lastName: string
  phone?: string
}

const DEFAULT_ADMIN: AdminUserConfig = {
  email: 'admin@logcha.com',
  password: 'admin123456',
  firstName: 'System',
  lastName: 'Administrator',
  phone: '+1234567890'
}

async function seedAdminUser(config: AdminUserConfig = DEFAULT_ADMIN) {
  console.log('🚀 Starting admin user seeding process...')
  
  try {
    // Check if admin user already exists
    console.log(`📋 Checking if user ${config.email} already exists...`)
    
    const { data: existingUsers, error: checkError } = await supabase
      .from('users')
      .select('email')
      .eq('email', config.email)
      .limit(1)

    if (checkError) {
      throw new Error(`Failed to check existing users: ${checkError.message}`)
    }

    if (existingUsers && existingUsers.length > 0) {
      console.log('⚠️  Admin user already exists with email:', config.email)
      console.log('   Skipping user creation...')
      return
    }

    // Create auth user
    console.log('👤 Creating authentication user...')
    
    const { data: authData, error: authError } = await supabase.auth.admin.createUser({
      email: config.email,
      password: config.password,
      email_confirm: true, // Auto-confirm email
      user_metadata: {
        first_name: config.firstName,
        last_name: config.lastName,
      }
    })

    if (authError) {
      throw new Error(`Failed to create auth user: ${authError.message}`)
    }

    if (!authData.user) {
      throw new Error('No user data returned from auth creation')
    }

    console.log('✅ Auth user created with ID:', authData.user.id)

    // Create user profile
    console.log('📝 Creating user profile...')
    
    const { error: profileError } = await supabase
      .from('users')
      .insert([
        {
          id: authData.user.id,
          email: config.email,
          first_name: config.firstName,
          last_name: config.lastName,
          phone: config.phone,
          role: 'company_admin',
          is_active: true
        }
      ])

    if (profileError) {
      // Clean up auth user if profile creation fails
      await supabase.auth.admin.deleteUser(authData.user.id)
      throw new Error(`Failed to create user profile: ${profileError.message}`)
    }

    console.log('✅ User profile created successfully')

    // Create a default company for the admin
    console.log('🏢 Creating default company...')
    
    const { data: companyData, error: companyError } = await supabase
      .from('companies')
      .insert([
        {
          name: 'Default Company',
          address: '123 Business Street, City, State 12345',
          contact_person: `${config.firstName} ${config.lastName}`,
          contact_email: config.email,
          contact_phone: config.phone
        }
      ])
      .select()
      .single()

    if (companyError) {
      console.log('⚠️  Warning: Could not create default company:', companyError.message)
    } else {
      console.log('✅ Default company created successfully')
    }

    console.log('\n🎉 Admin user seeding completed successfully!')
    console.log('\n📋 Admin User Details:')
    console.log(`   Email: ${config.email}`)
    console.log(`   Password: ${config.password}`)
    console.log(`   Name: ${config.firstName} ${config.lastName}`)
    console.log(`   Role: Company Administrator`)
    
    if (config.phone) {
      console.log(`   Phone: ${config.phone}`)
    }
    
    console.log('\n🔗 You can now log in at: http://localhost:3000/auth/login')
    
  } catch (error) {
    console.error('\n❌ Error seeding admin user:')
    console.error('  ', error instanceof Error ? error.message : String(error))
    process.exit(1)
  }
}

// Handle command line arguments
function parseArguments() {
  const args = process.argv.slice(2)
  const config = { ...DEFAULT_ADMIN }
  
  for (let i = 0; i < args.length; i += 2) {
    const key = args[i]
    const value = args[i + 1]
    
    switch (key) {
      case '--email':
        config.email = value
        break
      case '--password':
        config.password = value
        break
      case '--first-name':
        config.firstName = value
        break
      case '--last-name':
        config.lastName = value
        break
      case '--phone':
        config.phone = value
        break
      case '--help':
        showHelp()
        process.exit(0)
      default:
        if (key && !value) {
          console.error(`❌ Missing value for argument: ${key}`)
          showHelp()
          process.exit(1)
        }
    }
  }
  
  return config
}

function showHelp() {
  console.log('\n🔧 Logcha Admin User Seeder')
  console.log('\nUsage: bun run seed:admin [options]')
  console.log('\nOptions:')
  console.log('  --email <email>         Admin email (default: admin@logcha.com)')
  console.log('  --password <password>   Admin password (default: admin123456)')
  console.log('  --first-name <name>     Admin first name (default: System)')
  console.log('  --last-name <name>      Admin last name (default: Administrator)')
  console.log('  --phone <phone>         Admin phone number (optional)')
  console.log('  --help                  Show this help message')
  console.log('\nExample:')
  console.log('  bun run seed:admin --email admin@mycompany.com --password secretpass123')
  console.log('')
}

// Main execution
if (import.meta.main) {
  const config = parseArguments()
  await seedAdminUser(config)
}