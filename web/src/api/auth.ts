import { client } from './client';

export async function logout(): Promise<void> {
  // Override baseURL since auth routes are at root, not under /api
  await client.post('/auth/logout', null, { baseURL: '/' });
}
