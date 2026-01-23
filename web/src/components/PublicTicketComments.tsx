import { useQuery } from '@tanstack/react-query';
import { getPublicTicketComments } from '../api/public';
import { formatDistanceToNow } from 'date-fns';
import { useParams } from 'react-router-dom';

export default function PublicTicketComments() {
  const { token, ticketId } = useParams<{ token: string; ticketId: string }>();

  const { data: comments, isLoading, isError } = useQuery({
    queryKey: ['publicComments', token, ticketId],
    queryFn: () => getPublicTicketComments(token!, ticketId!),
    enabled: !!token && !!ticketId,
  });

  if (isLoading) {
    return <div className="p-4 text-center text-gray-500">Loading comments...</div>;
  }

  if (isError) {
    return <div className="p-4 text-center text-red-500">Failed to load comments</div>;
  }

  return (
    <div className="bg-white shadow sm:rounded-lg mt-8">
      <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
        <h3 className="text-lg leading-6 font-medium text-gray-900">Discussion</h3>
      </div>

      <div className="px-4 py-5 sm:px-6 space-y-6">
        {!comments || comments.length === 0 ? (
          <p className="text-gray-500 text-sm text-center py-4">No comments yet.</p>
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
                  <div className="text-sm flex items-center">
                    <span className="font-medium text-gray-900 mr-2">{comment.user.name}</span>
                    <span className="text-gray-500 mr-2">{formatDistanceToNow(new Date(comment.created_at), { addSuffix: true })}</span>
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
    </div>
  );
}
