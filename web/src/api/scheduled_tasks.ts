import { client as api } from './client';
import type { ScheduledTask } from '../types';

export const listScheduledTasks = async (organizationId: string) => {
  const response = await api.get<ScheduledTask[]>('/scheduled-tasks', {
    params: { organization_id: organizationId },
  });
  return response.data;
};

export const createScheduledTask = async (data: Partial<ScheduledTask>) => {
  // Optimization: send organization_id in query params to enable early authorization on backend
  const queryParams = data.organization_id ? { organization_id: data.organization_id } : undefined;
  const response = await api.post<ScheduledTask>('/scheduled-tasks', data, {
    params: queryParams,
  });
  return response.data;
};

export const updateScheduledTask = async (id: string, data: Partial<ScheduledTask>) => {
  const response = await api.patch<ScheduledTask>(`/scheduled-tasks/${id}`, data);
  return response.data;
};

export const deleteScheduledTask = async (id: string) => {
  await api.delete(`/scheduled-tasks/${id}`);
};
