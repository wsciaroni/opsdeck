import { useState, memo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getComments, createComment, type Comment } from '../api/comments';
import { formatDistanceToNow } from 'date-fns';
import clsx from 'clsx';
import toast from 'react-hot-toast';
import { Lock, Loader2 } from 'lucide-react';

interface TicketCommentsProps {
  ticketId: string;
}

// Memoized component to prevent re-rendering all comments when typing in the input (parent state change)
const CommentItem = memo(function CommentItem({ comment }: { readonly comment: Comment }) {
  return (
    <div className="flex space-x-3">
      <div className="flex-shrink-0">
        <img
          className="h-10 w-10 rounded-full bg-gray-300"
          src={comment.user.avatar_url || `https://ui-avatars.com/api/?name=${encodeURIComponent(comment.user.name)}`}
          alt={comment.user.name}
        />
      </div>
      <div className="min-w-0 flex-1">
        <div>
          <div className="text-sm flex items-center">
            <span className="font-medium text-gray-900 mr-2">{comment.user.name}</span>
            <span className="text-gray-500 mr-2">{formatDistanceToNow(new Date(comment.created_at), { addSuffix: true })}</span>
            {comment.sensitive && (
              <span className="inline-flex items-center rounded-full bg-red-100 px-2 py-0.5 text-xs font-medium text-red-800">
                <Lock className="w-3 h-3 mr-1" />
                Sensitive
              </span>
            )}
          </div>
          <div className="mt-1 text-sm text-gray-700">
            <p className="whitespace-pre-wrap">{comment.body}</p>
          </div>
        </div>
      </div>
    </div>
  );
});

export default function TicketComments({ ticketId }: TicketCommentsProps) {
  const [body, setBody] = useState('');
  const [sensitive, setSensitive] = useState(false);
  const queryClient = useQueryClient();

  const { data: comments, isLoading, isError } = useQuery({
    queryKey: ['comments', ticketId],
    queryFn: () => getComments(ticketId),
    enabled: !!ticketId,
  });

  const mutation = useMutation({
    mutationFn: () => createComment(ticketId, { body, sensitive }),
    onSuccess: () => {
      setBody('');
      setSensitive(false);
      queryClient.invalidateQueries({ queryKey: ['comments', ticketId] });
      toast.success("Comment posted!");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!body.trim()) return;
    mutation.mutate();
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      if (!body.trim()) return;
      mutation.mutate();
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
              <CommentItem key={comment.id} comment={comment} />
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
            <div className="mt-3 flex items-center justify-between">
              <div className="flex items-center">
                <input
                  id="comment-sensitive"
                  name="comment-sensitive"
                  type="checkbox"
                  className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  checked={sensitive}
                  onChange={(e) => setSensitive(e.target.checked)}
                />
                <label htmlFor="comment-sensitive" className="ml-2 block text-sm text-gray-900">
                  Mark as sensitive
                </label>
              </div>
              <button
                type="submit"
                disabled={mutation.isPending || !body.trim()}
                className={clsx(
                  "inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500",
                  mutation.isPending || !body.trim() ? "bg-indigo-400 cursor-not-allowed" : "bg-indigo-600 hover:bg-indigo-700"
                )}
              >
                {mutation.isPending ? (
                  <>
                    <Loader2 className="animate-spin h-4 w-4 mr-2" />
                    Posting...
                  </>
                ) : (
                  'Post Comment'
                )}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
