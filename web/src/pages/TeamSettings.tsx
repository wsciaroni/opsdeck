import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getMembers, addMember, removeMember } from '../api/organizations';
import { Trash2, UserPlus, AlertCircle, Users } from 'lucide-react';
import clsx from 'clsx';
import axios from 'axios';
import toast from 'react-hot-toast';
import EmptyState from '../components/EmptyState';

export default function TeamSettings() {
  const { orgId } = useParams<{ orgId: string }>();
  const queryClient = useQueryClient();
  const [newMemberEmail, setNewMemberEmail] = useState('');
  const [error, setError] = useState('');

  // Determine if we can show content (must match current org)
  // Although the route should probably be protected or we rely on API 403.

  const { data: members, isLoading } = useQuery({
    queryKey: ['members', orgId],
    queryFn: () => getMembers(orgId!),
    enabled: !!orgId,
  });

  const addMutation = useMutation({
    mutationFn: (email: string) => addMember(orgId!, email),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['members', orgId] });
      setNewMemberEmail('');
      setError('');
      toast.success("Member added!");
    },
    onError: (err: unknown) => {
      if (axios.isAxiosError(err) && err.response) {
        setError(String(err.response.data));
      } else {
        setError('Failed to add member');
      }
    },
  });

  const removeMutation = useMutation({
    mutationFn: (userId: string) => removeMember(orgId!, userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['members', orgId] });
    },
    onError: (err: unknown) => {
       if (axios.isAxiosError(err) && err.response) {
        alert(String(err.response.data));
      } else {
        alert('Failed to remove member');
      }
    },
  });

  const handleAddMember = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMemberEmail.trim()) return;
    addMutation.mutate(newMemberEmail);
  };

  const handleRemoveMember = (memberId: string) => {
    if (confirm('Are you sure you want to remove this member?')) {
      removeMutation.mutate(memberId);
    }
  };

  if (isLoading) return <div className="p-6">Loading members...</div>;

  // We can also check if currentOrg.id === orgId if we want to be strict on UI side

  return (
    <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div className="px-4 py-6 sm:px-0">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">Team Members</h1>

        {/* Add Member Form */}
        <div className="bg-white shadow sm:rounded-lg mb-8 p-6">
          <h2 className="text-lg font-medium text-gray-900 mb-4">Add Team Member</h2>
          <form onSubmit={handleAddMember} className="flex gap-4 items-start">
            <div className="flex-1">
              <label htmlFor="email" className="sr-only">Email address</label>
              <input
                type="email"
                name="email"
                id="email"
                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border"
                placeholder="Enter email address"
                value={newMemberEmail}
                onChange={(e) => setNewMemberEmail(e.target.value)}
                required
              />
              {error && (
                <p className="mt-2 text-sm text-red-600 flex items-center">
                  <AlertCircle className="h-4 w-4 mr-1" />
                  {error}
                </p>
              )}
            </div>
            <button
              type="submit"
              disabled={addMutation.isPending}
              className="inline-flex items-center justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50"
            >
              <UserPlus className="h-4 w-4 mr-2" />
              {addMutation.isPending ? 'Adding...' : 'Add Member'}
            </button>
          </form>
        </div>

        {/* Members List */}
        <div className="bg-white shadow overflow-hidden sm:rounded-lg">
          <ul role="list" className="divide-y divide-gray-200">
            {members?.map((member) => (
              <li key={member.id} className="px-4 py-4 sm:px-6 flex items-center justify-between hover:bg-gray-50">
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    {member.avatar_url ? (
                      <img className="h-10 w-10 rounded-full" src={member.avatar_url} alt="" />
                    ) : (
                      <span className="inline-block h-10 w-10 rounded-full overflow-hidden bg-gray-100">
                        <svg className="h-full w-full text-gray-300" fill="currentColor" viewBox="0 0 24 24">
                          <path d="M24 20.993V24H0v-2.996A14.977 14.977 0 0112.004 15c4.904 0 9.26 2.354 11.996 5.993zM16.002 8.999a4 4 0 11-8 0 4 4 0 018 0z" />
                        </svg>
                      </span>
                    )}
                  </div>
                  <div className="ml-4">
                    <div className="text-sm font-medium text-gray-900">{member.name || 'Unknown Name'}</div>
                    <div className="text-sm text-gray-500">{member.email}</div>
                  </div>
                </div>
                <div className="flex items-center space-x-4">
                  <span className={clsx(
                    "px-2 inline-flex text-xs leading-5 font-semibold rounded-full",
                    member.role === 'owner' ? "bg-green-100 text-green-800" : "bg-gray-100 text-gray-800"
                  )}>
                    {member.role}
                  </span>

                  {/* Show remove button if user is owner or removing themselves (and not the only owner logic which backend handles or we ignore for MVP) */}
                  {/* Actually, user can remove themselves. Owner can remove anyone. */}
                  {/* We need to know current user's role in this org. */}
                  {/* `currentOrg` from context might have role if we update it, or we check members list for current user. */}

                  <button
                    onClick={() => handleRemoveMember(member.id)}
                    className="text-gray-400 hover:text-red-600 transition-colors"
                    title="Remove member"
                  >
                    <Trash2 className="h-5 w-5" />
                  </button>
                </div>
              </li>
            ))}
          </ul>
          {members?.length === 0 && (
            <EmptyState
                title="No team members"
                description="Invite your team members to collaborate on tickets."
                icon={Users}
            />
          )}
        </div>
      </div>
    </div>
  );
}
