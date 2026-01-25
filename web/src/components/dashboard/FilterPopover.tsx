import { Popover, Transition } from '@headlessui/react';
import { Fragment } from 'react';
import { Filter, X } from 'lucide-react';
import { TICKET_STATUSES, TICKET_PRIORITIES } from '../../types';
import clsx from 'clsx';

interface FilterPopoverProps {
  status: string[] | undefined;
  setStatus: (status: string[] | undefined) => void;
  priority: string[] | undefined;
  setPriority: (priority: string[] | undefined) => void;
}

export default function FilterPopover({
  status,
  setStatus,
  priority,
  setPriority,
}: FilterPopoverProps) {
  // Helper: Is this status checked?
  const isStatusChecked = (id: string) => {
    if (status === undefined) {
      // Default: Active statuses are checked
      const s = TICKET_STATUSES.find((ts) => ts.id === id);
      return s ? !s.isFinished : false;
    }
    return status.includes(id);
  };

  // Helper: Is this priority checked?
  const isPriorityChecked = (id: string) => {
    if (!priority || priority.length === 0) return true; // All are checked by default
    return priority.includes(id);
  };

  const toggleStatus = (id: string) => {
    let newStatus: string[];

    // If currently undefined (Default), start with Active set
    if (status === undefined) {
      newStatus = TICKET_STATUSES.filter((t) => !t.isFinished).map((t) => t.id);
    } else {
      newStatus = [...status];
    }

    if (newStatus.includes(id)) {
      newStatus = newStatus.filter((s) => s !== id);
    } else {
      newStatus.push(id);
    }

    // If empty, user probably wants to see "All" (cleared filter), so we select ALL explicitly
    // because sending nothing triggers "Active" default in backend.
    if (newStatus.length === 0) {
       newStatus = TICKET_STATUSES.map(t => t.id);
    }

    setStatus(newStatus);
  };

  const togglePriority = (id: string) => {
    let newPriority: string[];
    // If undefined/empty (All), start with All
    if (!priority || priority.length === 0) {
       newPriority = TICKET_PRIORITIES.map(p => p.id);
    } else {
       newPriority = [...priority];
    }

    if (newPriority.includes(id)) {
      newPriority = newPriority.filter((p) => p !== id);
    } else {
      newPriority.push(id);
    }

    // If empty, set to undefined (All)
    if (newPriority.length === 0) {
        setPriority(undefined);
    } else {
        setPriority(newPriority);
    }
  };

  const handleMinPriorityChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const value = e.target.value;
    if (!value || value === 'any') {
      setPriority(undefined); // Any
      return;
    }

    const level = parseInt(value, 10);
    const newPrio = TICKET_PRIORITIES.filter((p) => p.level >= level).map((p) => p.id);
    setPriority(newPrio);
  };

  const resetToDefault = () => {
    setStatus(undefined);
    setPriority(undefined);
  };

  const isDefault = status === undefined && (priority === undefined || priority.length === 0);

  return (
    <Popover className="relative">
      {({ open }) => (
        <>
          <Popover.Button
            className={clsx(
              "inline-flex items-center justify-center rounded-md px-4 py-2 text-sm font-medium shadow-sm border focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2",
              open || !isDefault
                ? "bg-indigo-50 border-indigo-200 text-indigo-700"
                : "bg-white border-gray-300 text-gray-700 hover:bg-gray-50"
            )}
          >
            <Filter className="h-4 w-4 mr-2" />
            Filter
            {!isDefault && (
               <span className="ml-2 bg-indigo-200 text-indigo-800 py-0.5 px-2 rounded-full text-xs">
                 •
               </span>
            )}
          </Popover.Button>

          <Transition
            as={Fragment}
            enter="transition ease-out duration-200"
            enterFrom="opacity-0 translate-y-1"
            enterTo="opacity-100 translate-y-0"
            leave="transition ease-in duration-150"
            leaveFrom="opacity-100 translate-y-0"
            leaveTo="opacity-0 translate-y-1"
          >
            <Popover.Panel className="absolute right-0 z-10 mt-2 w-72 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
              <div className="p-4 space-y-6">

                {/* Header with Reset */}
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-medium text-gray-900">Filters</h3>
                  {!isDefault && (
                    <button
                        onClick={resetToDefault}
                        className="text-xs text-indigo-600 hover:text-indigo-800 flex items-center"
                    >
                        Reset to default
                    </button>
                  )}
                </div>

                {/* Status Section */}
                <div>
                    <h4 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Status</h4>
                    <div className="space-y-2">
                        {TICKET_STATUSES.map((s) => (
                            <label key={s.id} className="flex items-center space-x-2 cursor-pointer">
                                <input
                                    type="checkbox"
                                    className="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                                    checked={isStatusChecked(s.id)}
                                    onChange={() => toggleStatus(s.id)}
                                />
                                <span className={clsx("text-sm", isStatusChecked(s.id) ? "text-gray-900" : "text-gray-500")}>
                                    {s.label}
                                </span>
                            </label>
                        ))}
                    </div>
                </div>

                {/* Priority Section */}
                <div>
                    <div className="flex items-center justify-between mb-2">
                         <h4 className="text-xs font-semibold text-gray-500 uppercase tracking-wider">Priority</h4>
                    </div>

                    {/* Min Priority Selector */}
                     <select
                        className="block w-full rounded-md border-gray-300 py-1.5 pl-3 pr-8 text-xs focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 mb-3"
                        onChange={handleMinPriorityChange}
                        value=""
                     >
                        <option value="" disabled>Select minimum priority...</option>
                        <option value="any">Any Priority</option>
                        {TICKET_PRIORITIES.slice().reverse().map(p => (
                            <option key={p.id} value={p.level}>At least {p.label}</option>
                        ))}
                     </select>

                    <div className="space-y-2">
                        {TICKET_PRIORITIES.map((p) => (
                            <label key={p.id} className="flex items-center space-x-2 cursor-pointer">
                                <input
                                    type="checkbox"
                                    className="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                                    checked={isPriorityChecked(p.id)}
                                    onChange={() => togglePriority(p.id)}
                                />
                                <span className={clsx("text-sm", isPriorityChecked(p.id) ? "text-gray-900" : "text-gray-500")}>
                                    {p.label}
                                </span>
                            </label>
                        ))}
                    </div>
                </div>

              </div>
            </Popover.Panel>
          </Transition>
        </>
      )}
    </Popover>
  );
}
