import { createFileRoute } from '@tanstack/react-router'
import { useForm } from '@tanstack/react-form'
import { z } from 'zod'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

const signupSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Please enter a valid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ["confirmPassword"],
})

export const Route = createFileRoute('/signup')({
  component: SignUp,
})

function SignUp() {
  const form = useForm({
    defaultValues: {
      name: '',
      email: '',
      password: '',
      confirmPassword: '',
    },
    validators: {
      onChange: signupSchema,
    },
    onSubmit: async ({ value }) => {
      console.log('Sign up data:', value)
      // API call would go here
    },
  })

  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl text-center">Create an account</CardTitle>
          <CardDescription className="text-center">
            Enter your information to create a new account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form
            onSubmit={(e) => {
              e.preventDefault()
              e.stopPropagation()
              form.handleSubmit()
            }}
            className="space-y-4"
          >
            <div className="space-y-2">
              <form.Field
                name="name"
                children={(field) => (
                  <>
                    <Label htmlFor={field.name}>Name</Label>
                    <Input
                      id={field.name}
                      type="text"
                      placeholder="Enter your name"
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                    />
                    {field.state.meta.isTouched && field.state.meta.errors.length ? (
                      <div className="text-sm text-destructive">
                        {field.state.meta.errors.join(', ')}
                      </div>
                    ) : null}
                  </>
                )}
              />
            </div>
            <div className="space-y-2">
              <form.Field
                name="email"
                children={(field) => (
                  <>
                    <Label htmlFor={field.name}>Email</Label>
                    <Input
                      id={field.name}
                      type="email"
                      placeholder="Enter your email"
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                    />
                    {field.state.meta.isTouched && field.state.meta.errors.length ? (
                      <div className="text-sm text-destructive">
                        {field.state.meta.errors.join(', ')}
                      </div>
                    ) : null}
                  </>
                )}
              />
            </div>
            <div className="space-y-2">
              <form.Field
                name="password"
                children={(field) => (
                  <>
                    <Label htmlFor={field.name}>Password</Label>
                    <Input
                      id={field.name}
                      type="password"
                      placeholder="Enter your password"
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                    />
                    {field.state.meta.isTouched && field.state.meta.errors.length ? (
                      <div className="text-sm text-destructive">
                        {field.state.meta.errors.join(', ')}
                      </div>
                    ) : null}
                  </>
                )}
              />
            </div>
            <div className="space-y-2">
              <form.Field
                name="confirmPassword"
                children={(field) => (
                  <>
                    <Label htmlFor={field.name}>Confirm Password</Label>
                    <Input
                      id={field.name}
                      type="password"
                      placeholder="Confirm your password"
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                    />
                    {field.state.meta.isTouched && field.state.meta.errors.length ? (
                      <div className="text-sm text-destructive">
                        {field.state.meta.errors.join(', ')}
                      </div>
                    ) : null}
                  </>
                )}
              />
            </div>
            <form.Subscribe
              selector={(state) => [state.canSubmit, state.isSubmitting]}
              children={([canSubmit, isSubmitting]) => (
                <Button 
                  type="submit" 
                  className="w-full" 
                  disabled={!canSubmit}
                >
                  {isSubmitting ? 'Creating account...' : 'Create account'}
                </Button>
              )}
            />
          </form>
        </CardContent>
      </Card>
    </div>
  )
}