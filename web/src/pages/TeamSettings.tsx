import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getMembers, addMember, removeMember, updateMemberRole, getShareSettings, updateShareSettings, regenerateShareToken } from '../api/organizations';
import { Trash2, UserPlus, AlertCircle, Users, Link as LinkIcon, RefreshCw, Copy, Check } from 'lucide-react';
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import { ChevronDownIcon } from '@heroicons/react/20/solid';
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
      toast.success("Member removed");
    },
    onError: (err: unknown) => {
       if (axios.isAxiosError(err) && err.response) {
        toast.error(String(err.response.data));
      } else {
        toast.error('Failed to remove member');
      }
    },
  });

  const updateRoleMutation = useMutation({
    mutationFn: ({ userId, role }: { userId: string, role: string }) => updateMemberRole(orgId!, userId, role),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['members', orgId] });
      toast.success("Role updated");
    },
    onError: (err: unknown) => {
       if (axios.isAxiosError(err) && err.response) {
        toast.error(String(err.response.data));
      } else {
        toast.error('Failed to update role');
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

  const handleRoleChange = (memberId: string, role: string) => {
      updateRoleMutation.mutate({ userId: memberId, role });
  };

  // Share Link Logic
  const { data: shareSettings } = useQuery({
    queryKey: ['shareSettings', orgId],
    queryFn: () => getShareSettings(orgId!),
    enabled: !!orgId,
  });

  const updateShareMutation = useMutation({
    mutationFn: (enabled: boolean) => updateShareSettings(orgId!, enabled),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['shareSettings', orgId] });
      toast.success("Share settings updated");
    },
    onError: () => toast.error("Failed to update share settings"),
  });

  const regenerateTokenMutation = useMutation({
    mutationFn: () => regenerateShareToken(orgId!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['shareSettings', orgId] });
      toast.success("Link regenerated");
    },
    onError: () => toast.error("Failed to regenerate link"),
  });

  const [copied, setCopied] = useState(false);
  const copyLink = () => {
    if (shareSettings?.share_link_token) {
      const url = `${window.location.origin}/submit-ticket?token=${shareSettings.share_link_token}`;
      navigator.clipboard.writeText(url);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
      toast.success("Link copied to clipboard");
    }
  };

  if (isLoading) return <div className="p-6">Loading members...</div>;

  // We can also check if currentOrg.id === orgId if we want to be strict on UI side

  return (
    <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div className="px-4 py-6 sm:px-0 space-y-8">

        {/* Share Link Section */}
        <div className="bg-white shadow sm:rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-medium text-gray-900 flex items-center">
                    <LinkIcon className="h-5 w-5 mr-2 text-gray-500" />
                    Public Share Link
                </h2>
                <div className="flex items-center">
                     <span className="mr-3 text-sm text-gray-700">
                        {shareSettings?.share_link_enabled ? 'Enabled' : 'Disabled'}
                     </span>
                     <button
                        onClick={() => updateShareMutation.mutate(!shareSettings?.share_link_enabled)}
                        type="button"
                        className={clsx(
                            shareSettings?.share_link_enabled ? 'bg-indigo-600' : 'bg-gray-200',
                            'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2'
                        )}
                        role="switch"
                        aria-checked={shareSettings?.share_link_enabled}
                    >
                        <span
                            aria-hidden="true"
                            className={clsx(
                                shareSettings?.share_link_enabled ? 'translate-x-5' : 'translate-x-0',
                                'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out'
                            )}
                        />
                    </button>
                </div>
            </div>

            <p className="text-sm text-gray-500 mb-4">
                Allow anyone with the link to submit tickets to this organization.
            </p>

            {shareSettings?.share_link_enabled && (
                <div className="mt-2 flex rounded-md shadow-sm">
                    <div className="relative flex-grow focus-within:z-10">
                        <input
                            type="text"
                            readOnly
                            className="block w-full rounded-none rounded-l-md border-gray-300 focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border bg-gray-50 text-gray-500"
                            value={shareSettings?.share_link_token ? `${window.location.origin}/submit-ticket?token=${shareSettings.share_link_token}` : ''}
                        />
                    </div>
                    <button
                        type="button"
                        onClick={copyLink}
                        className="relative -ml-px inline-flex items-center border border-gray-300 bg-gray-50 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                    >
                        {copied ? <Check className="h-5 w-5 text-green-600" /> : <Copy className="h-5 w-5 text-gray-400" />}
                        <span className="sr-only">Copy link</span>
                    </button>
                    <button
                        type="button"
                        onClick={() => {
                            if (confirm('Regenerating the link will invalidate the old one. Continue?')) {
                                regenerateTokenMutation.mutate();
                            }
                        }}
                        className="relative -ml-px inline-flex items-center rounded-r-md border border-gray-300 bg-gray-50 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                        title="Regenerate Link"
                    >
                        <RefreshCw className={clsx("h-5 w-5 text-gray-400", regenerateTokenMutation.isPending && "animate-spin")} />
                    </button>
                </div>
            )}
        </div>

        <h1 className="text-2xl font-bold text-gray-900">Team Members</h1>

        {/* Add Member Form */}
        <div className="bg-white shadow sm:rounded-lg p-6">
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
          <ul className="divide-y divide-gray-200">
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
                  <Menu as="div" className="relative inline-block text-left">
                    <div>
                      <MenuButton className={clsx(
                        "inline-flex w-full justify-center gap-x-1.5 rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50",
                        member.role === 'owner' ? "text-green-800 bg-green-50 ring-green-300" : ""
                      )}>
                        {member.role}
                        <ChevronDownIcon className="-mr-1 h-5 w-5 text-gray-400" aria-hidden="true" />
                      </MenuButton>
                    </div>

                    <MenuItems
                      transition
                      anchor="bottom end"
                      className="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none data-[closed]:scale-95 data-[closed]:transform data-[closed]:opacity-0 data-[enter]:duration-100 data-[leave]:duration-75 data-[enter]:ease-out data-[leave]:ease-in"
                    >
                      <div className="py-1">
                        <MenuItem>
                            <button
                                onClick={() => handleRoleChange(member.id, 'member')}
                                className="block w-full px-4 py-2 text-left text-sm text-gray-700 data-[focus]:bg-gray-100 data-[focus]:text-gray-900"
                              >
                                Member
                              </button>
                        </MenuItem>
                        <MenuItem>
                            <button
                                onClick={() => handleRoleChange(member.id, 'admin')}
                                className="block w-full px-4 py-2 text-left text-sm text-gray-700 data-[focus]:bg-gray-100 data-[focus]:text-gray-900"
                              >
                                Admin
                              </button>
                        </MenuItem>
                        <MenuItem>
                            <button
                                onClick={() => handleRoleChange(member.id, 'owner')}
                                className="block w-full px-4 py-2 text-left text-sm text-gray-700 data-[focus]:bg-gray-100 data-[focus]:text-gray-900"
                              >
                                Owner
                              </button>
                        </MenuItem>
                      </div>
                    </MenuItems>
                  </Menu>

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
