import { type Ticket, TICKET_PRIORITIES } from './types';

export function formatBytes(bytes: number, decimals = 2): string {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

export function formatDate(date: string | Date): string {
  if (!date) return '';
  const d = new Date(date);
  if (isNaN(d.getTime())) return String(date);
  return d.toLocaleDateString();
}

export function getTicketAssignee(ticket: Ticket): string {
  return ticket.assignee_name || ticket.assignee_user_id || 'Unassigned';
}

export function getPriorityLabel(priorityId: string): string {
  const priority = TICKET_PRIORITIES.find((p) => p.id === priorityId);
  return priority ? priority.label : priorityId;
}
