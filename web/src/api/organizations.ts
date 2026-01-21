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

export const getShareSettings = async (orgID: string): Promise<{ share_link_enabled: boolean; share_link_token?: string }> => {
  const response = await client.get(`/organizations/${orgID}/share`);
  return response.data;
};

export const updateShareSettings = async (orgID: string, enabled: boolean): Promise<{ share_link_enabled: boolean; share_link_token?: string }> => {
  const response = await client.put(`/organizations/${orgID}/share`, { enabled });
  return response.data;
};

export const regenerateShareToken = async (orgID: string): Promise<{ share_link_enabled: boolean; share_link_token?: string }> => {
  const response = await client.post(`/organizations/${orgID}/share/regenerate`);
  return response.data;
};
