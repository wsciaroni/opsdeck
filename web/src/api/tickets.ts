import { client } from './client';
import type { Ticket, CreateTicketRequest } from '../types';

export async function getTickets(orgID: string): Promise<Ticket[]> {
  const response = await client.get(`/tickets?organization_id=${orgID}`);
  return response.data;
}

export async function createTicket(data: CreateTicketRequest): Promise<Ticket> {
  const response = await client.post('/tickets', data);
  return response.data;
}
