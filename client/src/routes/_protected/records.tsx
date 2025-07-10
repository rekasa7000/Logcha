import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { useAuth } from "@/lib/auth-context";
import { useTimeRecords } from "@/lib/time-queries";
import { Calendar, Clock, FileText } from "lucide-react";

export const Route = createFileRoute("/_protected/records")({
  component: RouteComponent,
});

function RouteComponent() {
  const { user } = useAuth();
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  
  const { data: recordsData, isLoading } = useTimeRecords(
    user?.id || 0,
    startDate || undefined,
    endDate || undefined
  );

  const records = recordsData?.data || [];

  const formatTime = (timeString: string | null) => {
    if (!timeString) return '--:--';
    return new Date(timeString).toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      weekday: 'short',
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'present':
        return 'bg-green-100 text-green-800';
      case 'half_day_am':
      case 'half_day_pm':
        return 'bg-yellow-100 text-yellow-800';
      case 'absent':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const totalHours = records.reduce((sum, record) => sum + record.total_hours, 0);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading records...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Time Records</h1>
          <p className="text-muted-foreground">
            View your historical time tracking data
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <FileText className="h-5 w-5 text-muted-foreground" />
          <span className="text-sm text-muted-foreground">
            {records.length} record{records.length !== 1 ? 's' : ''}
          </span>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Calendar className="h-5 w-5" />
            <span>Filter Records</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="startDate">Start Date</Label>
              <Input
                id="startDate"
                type="date"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="endDate">End Date</Label>
              <Input
                id="endDate"
                type="date"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {records.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Clock className="h-5 w-5" />
              <span>Summary</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-3">
              <div className="text-center">
                <div className="text-2xl font-bold">{totalHours.toFixed(1)}h</div>
                <p className="text-sm text-muted-foreground">Total Hours</p>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold">{records.length}</div>
                <p className="text-sm text-muted-foreground">Days Logged</p>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold">
                  {records.length > 0 ? (totalHours / records.length).toFixed(1) : 0}h
                </div>
                <p className="text-sm text-muted-foreground">Average per Day</p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      <div className="space-y-4">
        {records.length === 0 ? (
          <Card>
            <CardContent className="p-8 text-center">
              <FileText className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-semibold mb-2">No records found</h3>
              <p className="text-muted-foreground">
                {startDate || endDate
                  ? "No time records found for the selected date range."
                  : "You haven't logged any time yet. Start by clocking in from the dashboard."}
              </p>
            </CardContent>
          </Card>
        ) : (
          records.map((record) => (
            <Card key={record.id}>
              <CardContent className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center space-x-4">
                    <div>
                      <h3 className="font-semibold">{formatDate(record.date)}</h3>
                      <p className="text-sm text-muted-foreground">
                        Total: {record.total_hours.toFixed(1)} hours
                      </p>
                    </div>
                  </div>
                  <Badge className={getStatusColor(record.status)}>
                    {record.status}
                  </Badge>
                </div>
                
                <div className="grid gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <h4 className="font-medium text-sm">Morning Session</h4>
                    <div className="flex items-center space-x-2 text-sm font-mono">
                      <span>{formatTime(record.am_time_in)}</span>
                      <span>-</span>
                      <span>{formatTime(record.am_time_out)}</span>
                      <span className="text-muted-foreground">
                        ({record.am_hours.toFixed(1)}h)
                      </span>
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <h4 className="font-medium text-sm">Afternoon Session</h4>
                    <div className="flex items-center space-x-2 text-sm font-mono">
                      <span>{formatTime(record.pm_time_in)}</span>
                      <span>-</span>
                      <span>{formatTime(record.pm_time_out)}</span>
                      <span className="text-muted-foreground">
                        ({record.pm_hours.toFixed(1)}h)
                      </span>
                    </div>
                  </div>
                </div>
                
                {record.notes && (
                  <div className="mt-4 p-3 bg-muted rounded-lg">
                    <p className="text-sm">{record.notes}</p>
                  </div>
                )}
              </CardContent>
            </Card>
          ))
        )}
      </div>
    </div>
  );
}