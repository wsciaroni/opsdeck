import { client } from './client';

export interface Comment {
  id: string;
  body: string;
  created_at: string;
  user: {
    id: string;
    name: string;
    avatar_url: string;
  };
}

export interface CreateCommentRequest {
  body: string;
}

export async function getComments(ticketID: string): Promise<Comment[]> {
  const response = await client.get(`/tickets/${ticketID}/comments`);
  return response.data;
}

export async function createComment(ticketID: string, data: CreateCommentRequest): Promise<Comment> {
  const response = await client.post(`/tickets/${ticketID}/comments`, data);
  return response.data;
}
