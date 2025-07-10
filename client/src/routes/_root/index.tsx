import { createFileRoute, redirect } from "@tanstack/react-router";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { useAuth } from "@/lib/auth-context";
import { useTodayRecord, useTimeIn, useTimeOut } from "@/lib/time-queries";
import { Clock, LogIn, LogOut, Calendar, User, FileText } from "lucide-react";
import { ThemeToggleButton } from "@/components/shared/theme-button";

export const Route = createFileRoute("/")({
  component: App,
  beforeLoad: () => {
    const token = localStorage.getItem('token');
    if (!token) {
      throw redirect({
        to: '/signin',
      });
    }
  },
});

function App() {
  const { user, logout, isLoading } = useAuth();
  const [selectedSession, setSelectedSession] = useState<'am' | 'pm'>('am');
  
  const { data: todayData, isLoading: todayLoading } = useTodayRecord(user?.id || 0);
  const timeInMutation = useTimeIn();
  const timeOutMutation = useTimeOut();

  const todayRecord = todayData?.data;

  const handleTimeIn = async () => {
    if (!user) return;
    try {
      await timeInMutation.mutateAsync({
        trainee_id: user.id,
        session: selectedSession,
      });
    } catch (error) {
      console.error('Failed to clock in:', error);
    }
  };

  const handleTimeOut = async () => {
    if (!user) return;
    try {
      await timeOutMutation.mutateAsync({
        trainee_id: user.id,
        session: selectedSession,
      });
    } catch (error) {
      console.error('Failed to clock out:', error);
    }
  };

  const formatTime = (timeString: string | null) => {
    if (!timeString) return '--:--';
    return new Date(timeString).toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
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

  const canClockIn = (session: 'am' | 'pm') => {
    if (!todayRecord) return true;
    if (session === 'am') return !todayRecord.am_time_in;
    return !todayRecord.pm_time_in;
  };

  const canClockOut = (session: 'am' | 'pm') => {
    if (!todayRecord) return false;
    if (session === 'am') return todayRecord.am_time_in && !todayRecord.am_time_out;
    return todayRecord.pm_time_in && !todayRecord.pm_time_out;
  };

  if (isLoading || todayLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b">
        <div className="container mx-auto px-4 py-3 flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <img
              src="/logo-transparent.png"
              alt="Logcha"
              className="h-8 w-auto dark:invert"
            />
            <div className="text-sm text-muted-foreground">
              Welcome, {user.first_name} {user.last_name}
            </div>
          </div>
          <div className="flex items-center space-x-2">
            <ThemeToggleButton />
            <Button
              variant="outline"
              size="sm"
              onClick={logout}
              className="flex items-center space-x-2"
            >
              <LogOut className="h-4 w-4" />
              <span>Logout</span>
            </Button>
          </div>
        </div>
      </header>
      
      <main className="container mx-auto px-4 py-6">
        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
              <p className="text-muted-foreground">
                Track your time and view your records
              </p>
            </div>
            <div className="flex items-center space-x-2">
              <Calendar className="h-5 w-5 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">
                {new Date().toLocaleDateString('en-US', {
                  weekday: 'long',
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                })}
              </span>
            </div>
          </div>

          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Today's Status</CardTitle>
                <Badge className={getStatusColor(todayRecord?.status || 'absent')}>
                  {todayRecord?.status || 'absent'}
                </Badge>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {todayRecord?.total_hours.toFixed(1) || '0.0'}h
                </div>
                <p className="text-xs text-muted-foreground">
                  Total hours logged today
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Morning Session</CardTitle>
                <Clock className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-lg font-mono">
                  {formatTime(todayRecord?.am_time_in)} - {formatTime(todayRecord?.am_time_out)}
                </div>
                <p className="text-xs text-muted-foreground">
                  {todayRecord?.am_hours.toFixed(1) || '0.0'} hours
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Afternoon Session</CardTitle>
                <Clock className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-lg font-mono">
                  {formatTime(todayRecord?.pm_time_in)} - {formatTime(todayRecord?.pm_time_out)}
                </div>
                <p className="text-xs text-muted-foreground">
                  {todayRecord?.pm_hours.toFixed(1) || '0.0'} hours
                </p>
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <User className="h-5 w-5" />
                <span>Time Clock</span>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center space-x-4">
                <div className="flex space-x-2">
                  <Button
                    variant={selectedSession === 'am' ? 'default' : 'outline'}
                    size="sm"
                    onClick={() => setSelectedSession('am')}
                  >
                    Morning
                  </Button>
                  <Button
                    variant={selectedSession === 'pm' ? 'default' : 'outline'}
                    size="sm"
                    onClick={() => setSelectedSession('pm')}
                  >
                    Afternoon
                  </Button>
                </div>
              </div>

              <div className="flex space-x-4">
                <Button
                  onClick={handleTimeIn}
                  disabled={!canClockIn(selectedSession) || timeInMutation.isPending}
                  className="flex items-center space-x-2"
                >
                  <LogIn className="h-4 w-4" />
                  <span>
                    {timeInMutation.isPending ? 'Clocking in...' : 'Clock In'}
                  </span>
                </Button>
                <Button
                  onClick={handleTimeOut}
                  disabled={!canClockOut(selectedSession) || timeOutMutation.isPending}
                  variant="outline"
                  className="flex items-center space-x-2"
                >
                  <LogOut className="h-4 w-4" />
                  <span>
                    {timeOutMutation.isPending ? 'Clocking out...' : 'Clock Out'}
                  </span>
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  );
}
