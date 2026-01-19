import axios from 'axios';

export async function logout(): Promise<void> {
  await axios.post('/auth/logout');
}
