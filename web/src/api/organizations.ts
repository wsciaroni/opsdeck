import { client } from './client';
import type { Organization, Member } from '../types';

export const createOrganization = async (name: string): Promise<Organization> => {
  const response = await client.post('/organizations', { name });
  return response.data;
};

export const getMembers = async (orgID: string): Promise<Member[]> => {
  const response = await client.get(`/organizations/${orgID}/members`);
  return response.data;
};

export const addMember = async (orgID: string, email: string): Promise<void> => {
  await client.post(`/organizations/${orgID}/members`, { email });
};

export const removeMember = async (orgID: string, userID: string): Promise<void> => {
  await client.delete(`/organizations/${orgID}/members/${userID}`);
};
