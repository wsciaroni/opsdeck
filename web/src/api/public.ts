import { client } from './client';
import type { Ticket } from '../types';
import type { Comment } from './comments';

export interface PublicOrganization {
  id: string;
  name: string;
  slug: string;
}

export type PublicTicket = Ticket;

export async function getPublicOrganization(token: string): Promise<PublicOrganization> {
  const response = await client.get(`/public/view/${token}/organization`);
  return response.data;
}

export async function getPublicTickets(token: string): Promise<PublicTicket[]> {
  const response = await client.get(`/public/view/${token}/tickets`);
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
