import { client } from './client';
import type { Organization } from '../types';

export const createOrganization = async (name: string): Promise<Organization> => {
  const response = await client.post('/organizations', { name });
  return response.data;
};
