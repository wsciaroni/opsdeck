import { useState, Fragment } from 'react';
import { useAuth } from '../context/AuthContext';
import { createOrganization } from '../api/organizations';
import { Menu, Transition, Dialog } from '@headlessui/react';
import { User, LogOut, Check, Plus, Building } from 'lucide-react';
import { Link } from 'react-router-dom';
import clsx from 'clsx';

export default function ProfileDropdown() {
  const { user, currentOrg, organizations, switchOrganization, refreshOrganizations, logout } = useAuth();
  const [isOpen, setIsOpen] = useState(false);
  const [newOrgName, setNewOrgName] = useState('');
  const [isCreating, setIsCreating] = useState(false);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newOrgName.trim()) return;

    setIsCreating(true);
    try {
      const newOrg = await createOrganization(newOrgName);
      await refreshOrganizations();
      switchOrganization(newOrg.id);
      setIsOpen(false);
      setNewOrgName('');
    } catch (error) {
      console.error('Failed to create organization', error);
      // Ideally show an error message to the user
    } finally {
      setIsCreating(false);
    }
  };

  if (!user) return null;

  return (
    <>
      <Menu as="div" className="relative ml-3">
        <div>
          <Menu.Button className="flex rounded-full bg-white text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
            <span className="sr-only">Open user menu</span>
            <div className="h-8 w-8 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-600">
               <User className="h-5 w-5" />
            </div>
          </Menu.Button>
        </div>

        <Menu.Items
          anchor="bottom end"
          className="z-50 w-56 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none transition duration-200 ease-out data-[closed]:scale-95 data-[closed]:opacity-0"
        >
          <div className="px-4 py-2 border-b border-gray-100">
            <p className="text-sm text-gray-900 truncate font-medium">{user.email}</p>
          </div>

          <div className="py-1">
            <Menu.Item>
              {({ active }) => (
                <Link
                  to="/profile"
                  className={clsx(
                    active ? 'bg-gray-100' : '',
                    'flex px-4 py-2 text-sm text-gray-700 items-center'
                  )}
                >
                  <User className="mr-3 h-4 w-4 text-gray-400" aria-hidden="true" />
                  Profile Settings
                </Link>
              )}
            </Menu.Item>
            {currentOrg && (
              <Menu.Item>
                {({ active }) => (
                  <Link
                    to={`/organizations/${currentOrg.id}/settings/team`}
                    className={clsx(
                      active ? 'bg-gray-100' : '',
                      'flex px-4 py-2 text-sm text-gray-700 items-center'
                    )}
                  >
                    <Building className="mr-3 h-4 w-4 text-gray-400" aria-hidden="true" />
                    Organization Settings
                  </Link>
                )}
              </Menu.Item>
            )}
          </div>

          <div className="border-t border-gray-100 py-1">
             <div className="px-4 py-2">
                <p className="text-xs font-semibold text-gray-500 uppercase tracking-wider">
                  Switch Organization
                </p>
             </div>
             {organizations.map((org) => (
              <Menu.Item key={org.id}>
                {({ active }) => (
                  <button
                    onClick={() => switchOrganization(org.id)}
                    className={clsx(
                      active ? 'bg-gray-100' : '',
                      'group flex w-full items-center px-4 py-2 text-sm text-gray-700'
                    )}
                  >
                    <span className="flex-1 text-left truncate">{org.name}</span>
                    {currentOrg?.id === org.id && (
                      <Check className="ml-2 h-4 w-4 text-indigo-600" />
                    )}
                  </button>
                )}
              </Menu.Item>
            ))}
            <Menu.Item>
              {({ active }) => (
                <button
                  onClick={() => setIsOpen(true)}
                  className={clsx(
                    active ? 'bg-gray-100' : '',
                    'group flex w-full items-center px-4 py-2 text-sm text-indigo-600 font-medium'
                  )}
                >
                  <Plus className="mr-3 h-4 w-4" aria-hidden="true" />
                  Create New Organization
                </button>
              )}
            </Menu.Item>
          </div>

          <div className="border-t border-gray-100 py-1">
            <Menu.Item>
              {({ active }) => (
                <button
                  onClick={logout}
                  className={clsx(
                    active ? 'bg-gray-100' : '',
                    'flex w-full px-4 py-2 text-sm text-gray-700 items-center'
                  )}
                >
                  <LogOut className="mr-3 h-4 w-4 text-gray-400" aria-hidden="true" />
                  Logout
                </button>
              )}
            </Menu.Item>
          </div>
        </Menu.Items>
      </Menu>

      <Transition appear show={isOpen} as={Fragment}>
        <Dialog as="div" className="relative z-10" onClose={() => setIsOpen(false)}>
          {/* Dialog content preserved */}
          <Transition.Child
            as={Fragment}
            enter="ease-out duration-300"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            leave="ease-in duration-200"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <div className="fixed inset-0 bg-black bg-opacity-25" />
          </Transition.Child>

          <div className="fixed inset-0 overflow-y-auto">
            <div className="flex min-h-full items-center justify-center p-4 text-center">
              <Transition.Child
                as={Fragment}
                enter="ease-out duration-300"
                enterFrom="opacity-0 scale-95"
                enterTo="opacity-100 scale-100"
                leave="ease-in duration-200"
                leaveFrom="opacity-100 scale-100"
                leaveTo="opacity-0 scale-95"
              >
                <Dialog.Panel className="w-full max-w-md transform overflow-hidden rounded-2xl bg-white p-6 text-left align-middle shadow-xl transition-all">
                  <Dialog.Title
                    as="h3"
                    className="text-lg font-medium leading-6 text-gray-900"
                  >
                    Create Organization
                  </Dialog.Title>
                  <form onSubmit={handleCreate} className="mt-4">
                    <div>
                      <label htmlFor="orgName" className="block text-sm font-medium text-gray-700">
                        Organization Name
                      </label>
                      <input
                        type="text"
                        name="orgName"
                        id="orgName"
                        className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border"
                        placeholder="My Awesome Project"
                        value={newOrgName}
                        onChange={(e) => setNewOrgName(e.target.value)}
                        required
                      />
                    </div>

                    <div className="mt-6 flex justify-end space-x-3">
                      <button
                        type="button"
                        className="inline-flex justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                        onClick={() => setIsOpen(false)}
                      >
                        Cancel
                      </button>
                      <button
                        type="submit"
                        disabled={isCreating}
                        className="inline-flex justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50"
                      >
                        {isCreating ? 'Creating...' : 'Create'}
                      </button>
                    </div>
                  </form>
                </Dialog.Panel>
              </Transition.Child>
            </div>
          </div>
        </Dialog>
      </Transition>
    </>
  );
}
