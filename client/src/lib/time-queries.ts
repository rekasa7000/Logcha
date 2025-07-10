import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { apiClient, type TimeRecord, type TimeInRequest, type TimeOutRequest } from './api';

export const timeQueryKeys = {
  all: ['time'] as const,
  records: (traineeId: number) => [...timeQueryKeys.all, 'records', traineeId] as const,
  today: (traineeId: number) => [...timeQueryKeys.all, 'today', traineeId] as const,
  recordsWithDates: (traineeId: number, startDate?: string, endDate?: string) => 
    [...timeQueryKeys.records(traineeId), { startDate, endDate }] as const,
};

export function useTimeRecords(traineeId: number, startDate?: string, endDate?: string) {
  return useQuery({
    queryKey: timeQueryKeys.recordsWithDates(traineeId, startDate, endDate),
    queryFn: () => apiClient.getTimeRecords(traineeId, startDate, endDate),
    enabled: !!traineeId,
  });
}

export function useTodayRecord(traineeId: number) {
  return useQuery({
    queryKey: timeQueryKeys.today(traineeId),
    queryFn: () => apiClient.getTodayRecord(traineeId),
    enabled: !!traineeId,
    refetchInterval: 60000, // Refetch every minute
  });
}

export function useTimeIn() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: TimeInRequest) => apiClient.timeIn(data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: timeQueryKeys.today(variables.trainee_id),
      });
      queryClient.invalidateQueries({
        queryKey: timeQueryKeys.records(variables.trainee_id),
      });
    },
  });
}

export function useTimeOut() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: TimeOutRequest) => apiClient.timeOut(data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: timeQueryKeys.today(variables.trainee_id),
      });
      queryClient.invalidateQueries({
        queryKey: timeQueryKeys.records(variables.trainee_id),
      });
    },
  });
}