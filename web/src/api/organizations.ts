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

export const updateMemberRole = async (orgID: string, userID: string, role: string): Promise<void> => {
  await client.put(`/organizations/${orgID}/members/${userID}/role`, { role });
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

export const getPublicViewSettings = async (orgID: string): Promise<{ public_view_enabled: boolean; public_view_token?: string }> => {
  const response = await client.get(`/organizations/${orgID}/public-view`);
  return response.data;
};

export const updatePublicViewSettings = async (orgID: string, enabled: boolean): Promise<{ public_view_enabled: boolean; public_view_token?: string }> => {
  const response = await client.put(`/organizations/${orgID}/public-view`, { enabled });
  return response.data;
};

export const regeneratePublicViewToken = async (orgID: string): Promise<{ public_view_enabled: boolean; public_view_token?: string }> => {
  const response = await client.post(`/organizations/${orgID}/public-view/regenerate`);
  return response.data;
};
