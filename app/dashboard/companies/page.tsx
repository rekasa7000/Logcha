'use client'

import { useState, useEffect } from 'react'
import { supabase } from '@/lib/supabase'
import { Company } from '@/lib/types'
import { ProtectedRoute } from '@/components/protected-route'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Plus, Building2, Edit, Trash2 } from 'lucide-react'
import { CompanyDialog } from '@/components/company-dialog'
import { toast } from 'sonner'

export default function CompaniesPage() {
  const [companies, setCompanies] = useState<Company[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [selectedCompany, setSelectedCompany] = useState<Company | null>(null)

  useEffect(() => {
    fetchCompanies()
  }, [])

  const fetchCompanies = async () => {
    try {
      const { data, error } = await supabase
        .from('companies')
        .select('*')
        .order('name')

      if (error) throw error
      setCompanies(data || [])
    } catch (error) {
      toast.error('Failed to fetch companies')
      console.error('Error:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleDeleteCompany = async (company: Company) => {
    if (!confirm(`Are you sure you want to delete "${company.name}"?`)) {
      return
    }

    try {
      const { error } = await supabase
        .from('companies')
        .delete()
        .eq('id', company.id)

      if (error) throw error

      toast.success('Company deleted successfully')
      fetchCompanies()
    } catch (error) {
      toast.error('Failed to delete company')
      console.error('Error:', error)
    }
  }

  const handleEditCompany = (company: Company) => {
    setSelectedCompany(company)
    setDialogOpen(true)
  }

  const handleAddCompany = () => {
    setSelectedCompany(null)
    setDialogOpen(true)
  }

  const handleDialogClose = (success?: boolean) => {
    setDialogOpen(false)
    setSelectedCompany(null)
    if (success) {
      fetchCompanies()
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-gray-900"></div>
      </div>
    )
  }

  return (
    <ProtectedRoute allowedRoles={['company_admin']}>
      <div className="container mx-auto py-6">
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="flex items-center gap-2">
                  <Building2 className="h-5 w-5" />
                  Company Management
                </CardTitle>
                <CardDescription>
                  Manage companies and their information
                </CardDescription>
              </div>
              <Button onClick={handleAddCompany}>
                <Plus className="h-4 w-4 mr-2" />
                Add Company
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            {companies.length === 0 ? (
              <div className="text-center py-8">
                <Building2 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-gray-900 mb-2">No companies yet</h3>
                <p className="text-gray-600 mb-4">Get started by adding your first company.</p>
                <Button onClick={handleAddCompany}>
                  <Plus className="h-4 w-4 mr-2" />
                  Add Company
                </Button>
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Contact Person</TableHead>
                    <TableHead>Contact Email</TableHead>
                    <TableHead>Contact Phone</TableHead>
                    <TableHead>Address</TableHead>
                    <TableHead className="w-24">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {companies.map((company) => (
                    <TableRow key={company.id}>
                      <TableCell className="font-medium">{company.name}</TableCell>
                      <TableCell>{company.contact_person || '-'}</TableCell>
                      <TableCell>{company.contact_email || '-'}</TableCell>
                      <TableCell>{company.contact_phone || '-'}</TableCell>
                      <TableCell className="max-w-xs truncate" title={company.address || ''}>
                        {company.address || '-'}
                      </TableCell>
                      <TableCell>
                        <div className="flex space-x-2">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleEditCompany(company)}
                          >
                            <Edit className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleDeleteCompany(company)}
                            className="text-red-600 hover:text-red-700"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </CardContent>
        </Card>

        <CompanyDialog
          open={dialogOpen}
          company={selectedCompany}
          onClose={handleDialogClose}
        />
      </div>
    </ProtectedRoute>
  )
}