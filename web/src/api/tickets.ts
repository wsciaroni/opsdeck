import { client } from './client';
import type { Ticket, CreateTicketRequest, TicketDetail } from '../types';

export async function getTickets(orgID: string): Promise<Ticket[]> {
  const response = await client.get(`/tickets?organization_id=${orgID}`);
  return response.data;
}

export async function getTicket(id: string): Promise<TicketDetail> {
  const response = await client.get(`/tickets/${id}`);
  return response.data;
}

export async function createTicket(data: CreateTicketRequest): Promise<Ticket> {
  const response = await client.post('/tickets', data);
  return response.data;
}

export async function updateTicket(id: string, data: { status_id?: string; priority_id?: string }): Promise<Ticket> {
  const response = await client.patch(`/tickets/${id}`, data);
  return response.data;
}
