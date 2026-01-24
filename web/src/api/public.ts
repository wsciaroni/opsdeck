import { client } from './client';
import type { Ticket } from '../types';
import type { Comment } from './comments';

export interface PublicOrganization {
  id: string;
  name: string;
  slug: string;
}

export interface PublicTicket extends Ticket {
  reporter_name?: string;
  assignee_name?: string;
}

export async function getPublicOrganization(token: string): Promise<PublicOrganization> {
  const response = await client.get(`/public/view/${token}/organization`);
  return response.data;
}

export async function getPublicTickets(token: string, search?: string): Promise<PublicTicket[]> {
  const params = new URLSearchParams();
  if (search) {
    params.append('search', search);
  }
  const queryString = params.toString() ? `?${params.toString()}` : '';
  const response = await client.get(`/public/view/${token}/tickets${queryString}`);
  return response.data;
}

export async function getPublicTicket(token: string, ticketID: string): Promise<PublicTicket> {
  const response = await client.get(`/public/view/${token}/tickets/${ticketID}`);
  return response.data;
}

export async function getPublicTicketComments(token: string, ticketID: string): Promise<Comment[]> {
  const response = await client.get(`/public/view/${token}/tickets/${ticketID}/comments`);
  return response.data;
}
