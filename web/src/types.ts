export interface User {
  id: string;
  email: string;
  name: string;
  role: string;
  avatar_url: string;
  created_at: string;
  updated_at: string;
}

export interface Organization {
  id: string;
  name: string;
  slug: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface File {
  id: string;
  ticket_id: string;
  filename: string;
  content_type: string;
  size: number;
  created_at: string;
}

export interface Ticket {
  id: string;
  organization_id: string;
  title: string;
  description: string;
  location: string;
  status_id: string;
  priority_id: string;
  reporter_id: string;
  assignee_user_id: string | null;
  assignee_name?: string;
  reporter_name?: string;
  sensitive: boolean;
  created_at: string;
  updated_at: string;
  completed_at: string | null;
}

export interface TicketDetail extends Ticket {
  reporter_name: string;
  assignee_name?: string;
  files?: File[];
}

export interface CreateTicketRequest {
  title: string;
  description: string;
  priority_id: string;
  organization_id: string;
  location?: string;
  sensitive?: boolean;
}

export interface Member {
  id: string;
  email: string;
  name: string;
  avatar_url: string;
  role: string;
}

export const TICKET_STATUSES = [
  { id: 'new', label: 'New', isFinished: false },
  { id: 'in_progress', label: 'In Progress', isFinished: false },
  { id: 'on_hold', label: 'On Hold', isFinished: false },
  { id: 'done', label: 'Done', isFinished: true },
  { id: 'canceled', label: 'Canceled', isFinished: true },
] as const;

export const TICKET_PRIORITIES = [
  { id: 'critical', label: 'Critical', level: 4 },
  { id: 'high', label: 'High', level: 3 },
  { id: 'medium', label: 'Medium', level: 2 },
  { id: 'low', label: 'Low', level: 1 },
] as const;

export interface ScheduledTask {
  id: string;
  organization_id: string;
  title: string;
  description: string;
  frequency: string;
  start_date: string;
  next_run_at: string;
  created_by: string;
  assignee_user_id: string | null;
  assignee_user_name?: string;
  created_by_name?: string;
  priority_id: string;
  location: string;
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

export const FREQUENCIES = [
  { id: 'daily', label: 'Daily' },
  { id: 'weekly', label: 'Weekly' },
  { id: 'monthly', label: 'Monthly' },
  { id: 'yearly', label: 'Yearly' },
] as const;
