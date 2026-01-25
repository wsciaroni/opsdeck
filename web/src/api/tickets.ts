import { client } from './client';
import type { Ticket, CreateTicketRequest, TicketDetail } from '../types';

export async function getTickets(orgID: string, filters?: { status?: string | string[]; priority?: string | string[]; search?: string; sort_by?: string; sort_order?: string }): Promise<Ticket[]> {
  const params = new URLSearchParams();
  params.append('organization_id', orgID);

  if (filters?.status) {
    if (Array.isArray(filters.status)) {
      filters.status.forEach(s => params.append('status', s));
    } else {
      params.append('status', filters.status);
    }
  }

  if (filters?.priority) {
    if (Array.isArray(filters.priority)) {
      filters.priority.forEach(p => params.append('priority', p));
    } else {
      params.append('priority', filters.priority);
    }
  }

  if (filters?.search) params.append('search', filters.search);
  if (filters?.sort_by) params.append('sort_by', filters.sort_by);
  if (filters?.sort_order) params.append('sort_order', filters.sort_order);

  const response = await client.get(`/tickets?${params.toString()}`);
  return response.data;
}

export async function getTicket(id: string): Promise<TicketDetail> {
  const response = await client.get(`/tickets/${id}`);
  return response.data;
}

export async function createTicket(data: CreateTicketRequest | FormData): Promise<Ticket> {
  const response = await client.post('/tickets', data);
  return response.data;
}

export async function updateTicket(id: string, data: { status_id?: string; priority_id?: string; sensitive?: boolean; assignee_id?: string | null }): Promise<Ticket> {
  const response = await client.patch(`/tickets/${id}`, data);
  return response.data;
}

export async function createPublicTicket(data: { token: string; title: string; description: string; name: string; email: string; priority_id: string } | FormData): Promise<Ticket> {
  const response = await client.post('/public/tickets', data);
  return response.data;
}
