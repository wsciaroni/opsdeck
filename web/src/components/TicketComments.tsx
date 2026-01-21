import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getComments, createComment } from '../api/comments';
import { formatDistanceToNow } from 'date-fns';
import clsx from 'clsx';
import toast from 'react-hot-toast';

interface TicketCommentsProps {
  ticketId: string;
}

export default function TicketComments({ ticketId }: TicketCommentsProps) {
  const [body, setBody] = useState('');
  const queryClient = useQueryClient();

  const { data: comments, isLoading, isError } = useQuery({
    queryKey: ['comments', ticketId],
    queryFn: () => getComments(ticketId),
    enabled: !!ticketId,
  });

  const mutation = useMutation({
    mutationFn: (newBody: string) => createComment(ticketId, { body: newBody }),
    onSuccess: () => {
      setBody('');
      queryClient.invalidateQueries({ queryKey: ['comments', ticketId] });
      toast.success("Comment posted!");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!body.trim()) return;
    mutation.mutate(body);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      if (!body.trim()) return;
      mutation.mutate(body);
    }
  };

  if (isLoading) {
    return <div className="p-4 text-center text-gray-500">Loading comments...</div>;
  }

  if (isError) {
    return <div className="p-4 text-center text-red-500">Failed to load comments</div>;
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="bg-white shadow sm:rounded-lg">
        <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
          <h3 className="text-lg leading-6 font-medium text-gray-900">Discussion</h3>
        </div>

        {/* Comments List */}
        <div className="px-4 py-5 sm:px-6 space-y-6 max-h-[500px] overflow-y-auto">
          {!comments || comments.length === 0 ? (
            <p className="text-gray-500 text-sm text-center py-4">No comments yet. Start the conversation!</p>
          ) : (
            comments.map((comment) => (
              <div key={comment.id} className="flex space-x-3">
                <div className="flex-shrink-0">
                  <img
                    className="h-10 w-10 rounded-full bg-gray-300"
                    src={comment.user.avatar_url || `https://ui-avatars.com/api/?name=${encodeURIComponent(comment.user.name)}`}
                    alt={comment.user.name}
                  />
                </div>
                <div className="min-w-0 flex-1">
                  <div>
                    <div className="text-sm">
                      <span className="font-medium text-gray-900 mr-2">{comment.user.name}</span>
                      <span className="text-gray-500">{formatDistanceToNow(new Date(comment.created_at), { addSuffix: true })}</span>
                    </div>
                    <div className="mt-1 text-sm text-gray-700">
                      <p className="whitespace-pre-wrap">{comment.body}</p>
                    </div>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>

        {/* Comment Input */}
        <div className="bg-gray-50 px-4 py-4 sm:px-6">
          <form onSubmit={handleSubmit}>
            <div>
              <label htmlFor="comment" className="sr-only">
                Add your comment
              </label>
              <textarea
                id="comment"
                name="comment"
                rows={3}
                className="shadow-sm block w-full focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm border border-gray-300 rounded-md p-2"
                placeholder="Add a comment... (Press Enter to send)"
                value={body}
                onChange={(e) => setBody(e.target.value)}
                onKeyDown={handleKeyDown}
                disabled={mutation.isPending}
              />
            </div>
            <div className="mt-3 flex items-center justify-end">
              <button
                type="submit"
                disabled={mutation.isPending || !body.trim()}
                className={clsx(
                  "inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500",
                  mutation.isPending || !body.trim() ? "bg-indigo-400 cursor-not-allowed" : "bg-indigo-600 hover:bg-indigo-700"
                )}
              >
                {mutation.isPending ? 'Posting...' : 'Post Comment'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
